package operation

type Manifest struct {
	Applications []AppManifest `yaml:"applications"`
}

type AppManifest struct {
	Name       string   `yaml:"name"`
	Buildpacks []string `yaml:"buildpacks,omitempty"`
	Command    string   `yaml:"command,omitempty"`
	DiskQuota  string   `yaml:"disk_quota,omitempty"`
	Docker     struct {
		Image    string `yaml:"image,omitempty"`
		Username string `yaml:"username,omitempty"`
	} `yaml:"docker,omitempty"`
	Env                     map[string]string   `yaml:"env,omitempty"`
	HealthCheckType         string              `yaml:"health-check-type,omitempty"`
	HealthCheckHTTPEndpoint string              `yaml:"health-check-http-endpoint,omitempty"`
	Instances               int                 `yaml:"instances,omitempty"`
	LogRateLimit            string              `yaml:"log-rate-limit,omitempty"`
	Memory                  string              `yaml:"memory,omitempty"`
	NoRoute                 bool                `yaml:"no-route,omitempty"`
	Routes                  []AppManifestRoutes `yaml:"routes,omitempty"`
	Services                []string            `yaml:"services,omitempty"`
	Stack                   string              `yaml:"stack,omitempty"`
	Timeout                 int                 `yaml:"timeout,omitempty"`
}

type AppManifestRoutes struct {
	Route string `yaml:"route,omitempty"`
}
