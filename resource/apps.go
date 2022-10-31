package resource

import (
	"encoding/json"
	"time"
)

type App struct {
	GUID          string                       `json:"guid"`
	CreatedAt     time.Time                    `json:"created_at"`
	UpdatedAt     time.Time                    `json:"updated_at"`
	Name          string                       `json:"name"`
	State         string                       `json:"state,omitempty"`
	Lifecycle     Lifecycle                    `json:"lifecycle,omitempty"`
	Relationships map[string]ToOneRelationship `json:"relationships,omitempty"`
	Links         map[string]Link              `json:"links,omitempty"`
	Metadata      Metadata                     `json:"metadata"`
}

type AppCreate struct {
	Name                 string            `json:"name"`
	Relationships        SpaceRelationship `json:"relationships"`
	EnvironmentVariables map[string]string `json:"environment_variables,omitempty"`
	Lifecycle            *Lifecycle        `json:"lifecycle,omitempty"`
	Metadata             *Metadata         `json:"metadata,omitempty"`
}

type AppUpdate struct {
	Name      string     `json:"name"`
	Lifecycle *Lifecycle `json:"lifecycle,omitempty"`
	Metadata  *Metadata  `json:"metadata,omitempty"`
}

type AppList struct {
	Pagination Pagination `json:"pagination"`
	Resources  []*App     `json:"resources"`
}

type AppSSHEnabled struct {
	Enabled bool   `json:"enabled"`
	Reason  string `json:"reason"`
}

type AppPermissions struct {
	ReadBasicData     bool `json:"read_basic_data"`
	ReadSensitiveData bool `json:"read_sensitive_data"`
}

type AppEnvironment struct {
	EnvVars       map[string]string          `json:"environment_variables,omitempty"`
	StagingEnv    map[string]string          `json:"staging_env_json,omitempty"`
	RunningEnv    map[string]string          `json:"running_env_json,omitempty"`
	SystemEnvVars map[string]json.RawMessage `json:"system_env_json,omitempty"`      // VCAP_SERVICES
	AppEnvVars    map[string]json.RawMessage `json:"application_env_json,omitempty"` // VCAP_APPLICATION
}

type EnvVar struct {
	Var map[string]*string `json:"var"`
}

type EnvVarResponse struct {
	EnvVar
	Links map[string]Link `json:"links"`
}

type Lifecycle struct {
	Type          string             `json:"type,omitempty"`
	BuildpackData BuildpackLifecycle `json:"data,omitempty"` // TODO: support other lifecycles
}

type BuildpackLifecycle struct {
	Buildpacks []string `json:"buildpacks,omitempty"`
	Stack      string   `json:"stack,omitempty"`
}

func NewAppCreate(name, spaceGUID string) *AppCreate {
	return &AppCreate{
		Name: name,
		Relationships: SpaceRelationship{
			Space: ToOneRelationship{
				Data: Relationship{
					GUID: spaceGUID,
				},
			},
		},
	}
}
