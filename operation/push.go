package operation

import (
	"context"
	"fmt"
	"github.com/cloudfoundry-community/go-cfclient/v3/client"
	"github.com/cloudfoundry-community/go-cfclient/v3/resource"
	"gopkg.in/yaml.v3"
	"io"
)

// AppPushOperation can be used to push buildpack apps
type AppPushOperation struct {
	orgName   string
	spaceName string
	client    *client.Client
}

// NewAppPushOperation creates a new AppPushOperation
func NewAppPushOperation(client *client.Client, orgName, spaceName string) *AppPushOperation {
	return &AppPushOperation{
		orgName:   orgName,
		spaceName: spaceName,
		client:    client,
	}
}

// Push creates or updates an application using the specified manifest and zipped source files
func (p *AppPushOperation) Push(ctx context.Context, appManifest *AppManifest, zipFile io.Reader) (*resource.App, error) {
	org, err := p.findOrg(ctx)
	if err != nil {
		return nil, err
	}
	space, err := p.findSpace(ctx, org.GUID)
	if err != nil {
		return nil, err
	}
	return p.pushApp(ctx, space, appManifest, zipFile)
}

// pushApp pushes an application
//
// After an application is created and packages are uploaded, a droplet must be created via a build in order for
// an application to be deployed or tasks to be run. The current droplet must be assigned to an application before
// it may be started. When tasks are created, they either use a specific droplet guid, or use the current droplet
// assigned to an application.
func (p *AppPushOperation) pushApp(ctx context.Context, space *resource.Space, manifest *AppManifest, zipFile io.Reader) (*resource.App, error) {
	err := p.applySpaceManifest(ctx, space, manifest)
	if err != nil {
		return nil, err
	}

	app, err := p.findApp(ctx, manifest.Name, space)
	if err != nil {
		return nil, err
	}

	pkg, err := p.uploadPackage(ctx, app, zipFile)
	if err != nil {
		return nil, err
	}

	droplet, err := p.buildDroplet(ctx, pkg, manifest)
	if err != nil {
		return nil, err
	}

	_, err = p.client.Droplets.SetCurrentAssociationForApp(ctx, app.GUID, droplet.GUID)
	if err != nil {
		return nil, err
	}

	return p.client.Applications.Start(ctx, app.GUID)
}

func (p *AppPushOperation) applySpaceManifest(ctx context.Context, space *resource.Space, manifest *AppManifest) error {
	// wrap it in a manifest that has an applications array as required by the API
	multiAppsManifest := &Manifest{
		Applications: []*AppManifest{manifest},
	}
	manifestBytes, err := yaml.Marshal(&multiAppsManifest)
	if err != nil {
		return fmt.Errorf("error marshalling application manifest: %w", err)
	}

	jobGUID, err := p.client.Manifests.ApplyManifest(ctx, space.GUID, string(manifestBytes))
	if err != nil {
		return fmt.Errorf("error applying application manifest to space %s: %w", space.Name, err)
	}
	err = p.client.Jobs.PollComplete(ctx, jobGUID, nil)
	if err != nil {
		return fmt.Errorf("error waiting for application manifest to finish applying to space %s: %w", space.Name, err)
	}
	return nil
}

func (p *AppPushOperation) findApp(ctx context.Context, appName string, space *resource.Space) (*resource.App, error) {
	appOpts := client.NewAppListOptions()
	appOpts.Names.Values = []string{appName}
	appOpts.SpaceGUIDs.Values = []string{space.GUID}
	apps, err := p.client.Applications.ListAll(ctx, appOpts)
	if err != nil {
		return nil, err
	}
	if len(apps) != 1 {
		return nil, fmt.Errorf("expected to find one application named %s in space %s, but found %d",
			appName, space.Name, len(apps))
	}
	return apps[0], nil
}

func (p *AppPushOperation) uploadPackage(ctx context.Context, app *resource.App, zipFile io.Reader) (*resource.Package, error) {
	newPkg := resource.NewPackageCreate(app.GUID)
	pkg, err := p.client.Packages.Create(ctx, newPkg)
	if err != nil {
		return nil, fmt.Errorf("error creating package for app %s: %w", app.Name, err)
	}

	err = p.client.Packages.UploadBits(ctx, pkg.GUID, zipFile)
	if err != nil {
		return nil, fmt.Errorf("error uploading package bits for app %s: %w", app.Name, err)
	}
	err = p.client.Packages.PollReady(ctx, pkg.GUID, nil)
	if err != nil {
		return nil, fmt.Errorf("error while waiting for package to process for app %s: %w", app.Name, err)
	}
	return pkg, nil
}

func (p *AppPushOperation) buildDroplet(ctx context.Context, pkg *resource.Package, manifest *AppManifest) (*resource.Droplet, error) {
	newBuild := resource.NewBuildCreate(pkg.GUID)
	newBuild.Lifecycle = &resource.Lifecycle{
		Type: "buildpack",
		BuildpackData: resource.BuildpackLifecycle{
			Buildpacks: manifest.Buildpacks,
			Stack:      manifest.Stack,
		},
	}
	build, err := p.client.Builds.Create(ctx, newBuild)
	if err != nil {
		return nil, fmt.Errorf("error creating build from package for app %s: %w", manifest.Name, err)
	}
	err = p.client.Builds.PollStaged(ctx, build.GUID, nil)
	if err != nil {
		return nil, fmt.Errorf("error while waiting for app %s package to build: %w", manifest.Name, err)
	}

	opts := client.NewDropletPackageListOptions()
	opts.States.Values = []string{string(resource.DropletStateStaged)}
	droplets, err := p.client.Droplets.ListForPackageAll(ctx, pkg.GUID, opts)
	if err != nil {
		return nil, fmt.Errorf("error finding droplet for app %s: %w", manifest.Name, err)
	}
	if len(droplets) != 1 {
		return nil, fmt.Errorf("expected one droplet, but found %d", len(droplets))
	}
	return droplets[0], nil
}

func (p *AppPushOperation) findOrg(ctx context.Context) (*resource.Organization, error) {
	opts := client.NewOrgListOptions()
	opts.Names.Values = []string{p.orgName}
	orgs, err := p.client.Organizations.ListAll(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("could not find org %s: %w", p.orgName, err)
	}
	if len(orgs) != 1 {
		return nil, fmt.Errorf("expected to find one org named %s, but found %d", p.orgName, len(orgs))
	}
	return orgs[0], nil
}

func (p *AppPushOperation) findSpace(ctx context.Context, orgGUID string) (*resource.Space, error) {
	opts := client.NewSpaceListOptions()
	opts.Names.Values = []string{p.spaceName}
	opts.OrganizationGUIDs.Values = []string{orgGUID}
	spaces, err := p.client.Spaces.ListAll(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("could not find space %s: %w", p.spaceName, err)
	}
	if len(spaces) != 1 {
		return nil, fmt.Errorf("expected to find one space named %s, but found %d", p.spaceName, len(spaces))
	}
	return spaces[0], nil
}
