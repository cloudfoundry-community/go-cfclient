package client

import (
	"github.com/cloudfoundry-community/go-cfclient/resource"
	"github.com/cloudfoundry-community/go-cfclient/test"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func TestOrgs(t *testing.T) {
	g := test.NewObjectJSONGenerator(15)
	org := g.Organization()
	org2 := g.Organization()
	org3 := g.Organization()
	org4 := g.Organization()
	domain := g.Domain()
	orgUsageSummary := g.OrganizationUsageSummary()
	user := g.User()
	user2 := g.User()

	tests := []RouteTest{
		{
			Description: "Create org",
			Route: MockRoute{
				Method:   "POST",
				Endpoint: "/v3/organizations",
				Output:   []string{org},
				Status:   http.StatusCreated,
				PostForm: `{ "name": "my-organization" }`,
			},
			Expected: org,
			Action: func(c *Client, t *testing.T) (any, error) {
				r := resource.NewOrganizationCreate("my-organization")
				return c.Organizations.Create(r)
			},
		},
		{
			Description: "Get org",
			Route: MockRoute{
				Method:   "GET",
				Endpoint: "/v3/organizations/3691e277-eb88-4ddc-bec3-0111d9dd4ef5",
				Output:   []string{org},
				Status:   http.StatusOK,
			},
			Expected: org,
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.Organizations.Get("3691e277-eb88-4ddc-bec3-0111d9dd4ef5")
			},
		},
		{
			Description: "Get org default iso segment",
			Route: MockRoute{
				Method:   "GET",
				Endpoint: "/v3/organizations/3691e277-eb88-4ddc-bec3-0111d9dd4ef5/relationships/default_isolation_segment",
				Output:   []string{`{ "data": { "guid": "443a1ea0-2403-4f0f-8c74-023a320bd1f2" }}`},
				Status:   http.StatusOK,
			},
			Action: func(c *Client, t *testing.T) (any, error) {
				iso, err := c.Organizations.GetDefaultIsoSegment("3691e277-eb88-4ddc-bec3-0111d9dd4ef5")
				require.NoError(t, err)
				require.Equal(t, "443a1ea0-2403-4f0f-8c74-023a320bd1f2", iso)
				return nil, nil
			},
		},
		{
			Description: "Get org default domain",
			Route: MockRoute{
				Method:   "GET",
				Endpoint: "/v3/organizations/3691e277-eb88-4ddc-bec3-0111d9dd4ef5/domains/default",
				Output:   []string{domain},
				Status:   http.StatusOK,
			},
			Expected: domain,
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.Organizations.GetDefaultDomain("3691e277-eb88-4ddc-bec3-0111d9dd4ef5")
			},
		},
		{
			Description: "Get org usage summary",
			Route: MockRoute{
				Method:   "GET",
				Endpoint: "/v3/organizations/3691e277-eb88-4ddc-bec3-0111d9dd4ef5/usage_summary",
				Output:   []string{orgUsageSummary},
				Status:   http.StatusOK,
			},
			Expected: orgUsageSummary,
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.Organizations.GetUsageSummary("3691e277-eb88-4ddc-bec3-0111d9dd4ef5")
			},
		},
		{
			Description: "Delete org",
			Route: MockRoute{
				Method:   "DELETE",
				Endpoint: "/v3/organizations/3691e277-eb88-4ddc-bec3-0111d9dd4ef5",
				Status:   http.StatusAccepted,
			},
			Action: func(c *Client, t *testing.T) (any, error) {
				return nil, c.Organizations.Delete("3691e277-eb88-4ddc-bec3-0111d9dd4ef5")
			},
		},
		{
			Description: "List all orgs",
			Route: MockRoute{
				Method:   "GET",
				Endpoint: "/v3/organizations",
				Output:   g.Paged([]string{org, org2}, []string{org3, org4}),
				Status:   http.StatusOK,
			},
			Expected: g.Array(org, org2, org3, org4),
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.Organizations.ListAll(nil)
			},
		},
		{
			Description: "List all orgs for iso segment",
			Route: MockRoute{
				Method:   "GET",
				Endpoint: "/v3/isolation_segments/571de34f-8067-44f0-8bec-4ac17bf8750f/organizations",
				Output:   g.Paged([]string{org, org2}, []string{org3, org4}),
				Status:   http.StatusOK,
			},
			Expected: g.Array(org, org2, org3, org4),
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.Organizations.ListForIsoSegmentAll("571de34f-8067-44f0-8bec-4ac17bf8750f", nil)
			},
		},
		{
			Description: "List all org users",
			Route: MockRoute{
				Method:   "GET",
				Endpoint: "/v3/organizations/3691e277-eb88-4ddc-bec3-0111d9dd4ef5/users",
				Output:   g.Paged([]string{user}, []string{user2}),
				Status:   http.StatusOK},
			Expected: g.Array(user, user2),
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.Organizations.ListUsersAll("3691e277-eb88-4ddc-bec3-0111d9dd4ef5", nil)
			},
		},
		{
			Description: "Update org",
			Route: MockRoute{
				Method:   "PATCH",
				Endpoint: "/v3/organizations/3691e277-eb88-4ddc-bec3-0111d9dd4ef5",
				Output:   []string{org},
				Status:   http.StatusOK,
				PostForm: `{ "name": "new_name" }`,
			},
			Expected: org,
			Action: func(c *Client, t *testing.T) (any, error) {
				r := &resource.OrganizationUpdate{
					Name: "new_name",
				}
				return c.Organizations.Update("3691e277-eb88-4ddc-bec3-0111d9dd4ef5", r)
			},
		},
		{
			Description: "Assign default org iso segment",
			Route: MockRoute{
				Method:   "PATCH",
				Endpoint: "/v3/organizations/3691e277-eb88-4ddc-bec3-0111d9dd4ef5/relationships/default_isolation_segment",
				Output:   []string{`{ "data": { "guid": "443a1ea0-2403-4f0f-8c74-023a320bd1f2" }}`},
				Status:   http.StatusOK,
				PostForm: `{ "data": { "guid": "443a1ea0-2403-4f0f-8c74-023a320bd1f2" }}`,
			},
			Action: func(c *Client, t *testing.T) (any, error) {
				err := c.Organizations.AssignDefaultIsoSegment("3691e277-eb88-4ddc-bec3-0111d9dd4ef5", "443a1ea0-2403-4f0f-8c74-023a320bd1f2")
				return nil, err
			},
		},
	}
	executeTests(tests, t)
}