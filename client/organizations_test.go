package client

import (
	"github.com/cloudfoundry-community/go-cfclient/resource"
	"github.com/cloudfoundry-community/go-cfclient/test"
	"net/http"
	"testing"
)

func TestOrgs(t *testing.T) {
	g := test.NewObjectJSONGenerator(15)
	org := g.Organization()
	org2 := g.Organization()
	org3 := g.Organization()
	org4 := g.Organization()

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
				Status:   http.StatusOK},
			Expected: org,
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.Organizations.Get("3691e277-eb88-4ddc-bec3-0111d9dd4ef5")
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
			Description: "List all orgs",
			Route: MockRoute{
				Method:   "GET",
				Endpoint: "/v3/organizations",
				Output:   g.Paged([]string{org, org2}, []string{org3, org4}),
				Status:   http.StatusOK},
			Expected: g.Array(org, org2, org3, org4),
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.Organizations.ListAll(nil)
			},
		},
	}
	executeTests(tests, t)
}
