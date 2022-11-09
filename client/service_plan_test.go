package client

import (
	"github.com/cloudfoundry-community/go-cfclient/v3/resource"
	"github.com/cloudfoundry-community/go-cfclient/v3/test"
	"net/http"
	"testing"
)

func TestServicePlans(t *testing.T) {
	g := test.NewObjectJSONGenerator(156)
	svcPlan := g.ServicePlan()
	svcPlan2 := g.ServicePlan()
	svcPlan3 := g.ServicePlan()
	svcPlan4 := g.ServicePlan()
	space := g.Space()
	space2 := g.Space()
	org := g.Organization()
	svcOffering := g.ServiceOffering()

	tests := []RouteTest{
		{
			Description: "Delete service plan",
			Route: MockRoute{
				Method:   "DELETE",
				Endpoint: "/v3/service_plans/79aae221-b2a6-4aaa-a134-76f605af46c9",
				Status:   http.StatusNoContent,
			},
			Action: func(c *Client, t *testing.T) (any, error) {
				err := c.ServicePlans.Delete("79aae221-b2a6-4aaa-a134-76f605af46c9")
				return nil, err
			},
		},
		{
			Description: "Get service plan",
			Route: MockRoute{
				Method:   "GET",
				Endpoint: "/v3/service_plans/79aae221-b2a6-4aaa-a134-76f605af46c9",
				Output:   []string{svcPlan},
				Status:   http.StatusOK},
			Expected: svcPlan,
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.ServicePlans.Get("79aae221-b2a6-4aaa-a134-76f605af46c9")
			},
		},
		{
			Description: "List all service plans",
			Route: MockRoute{
				Method:   "GET",
				Endpoint: "/v3/service_plans",
				Output:   g.Paged([]string{svcPlan}, []string{svcPlan2}),
				Status:   http.StatusOK},
			Expected: g.Array(svcPlan, svcPlan2),
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.ServicePlans.ListAll(nil)
			},
		},
		{
			Description: "List all service plans include service offerings",
			Route: MockRoute{
				Method:   "GET",
				Endpoint: "/v3/service_plans",
				Output: g.PagedWithInclude(
					test.PagedResult{
						Resources:        []string{svcPlan, svcPlan2},
						ServiceOfferings: []string{svcOffering},
					},
					test.PagedResult{
						Resources: []string{svcPlan3, svcPlan4},
					}),
				Status: http.StatusOK},
			Expected:  g.Array(svcPlan, svcPlan2, svcPlan3, svcPlan4),
			Expected2: g.Array(svcOffering),
			Action2: func(c *Client, t *testing.T) (any, any, error) {
				return c.ServicePlans.ListIncludeServiceOfferingAll(nil)
			},
		},
		{
			Description: "List all service plans include spaces and orgs",
			Route: MockRoute{
				Method:   "GET",
				Endpoint: "/v3/service_plans",
				Output: g.PagedWithInclude(
					test.PagedResult{
						Resources:     []string{svcPlan, svcPlan2},
						Spaces:        []string{space},
						Organizations: []string{org},
					},
					test.PagedResult{
						Resources: []string{svcPlan3, svcPlan4},
						Spaces:    []string{space2},
					}),
				Status: http.StatusOK},
			Expected:  g.Array(svcPlan, svcPlan2, svcPlan3, svcPlan4),
			Expected2: g.Array(space, space2),
			Expected3: g.Array(org),
			Action3: func(c *Client, t *testing.T) (any, any, any, error) {
				return c.ServicePlans.ListIncludeSpacesAndOrgsAll(nil)
			},
		},
		{
			Description: "Update service plan",
			Route: MockRoute{
				Method:   "PATCH",
				Endpoint: "/v3/service_plans/79aae221-b2a6-4aaa-a134-76f605af46c9",
				Output:   []string{svcPlan},
				Status:   http.StatusOK,
				PostForm: `{
					"metadata": {
					  "labels": {"key": "value"},
					  "annotations": {"note": "detailed information"}
					}
				  }`,
			},
			Expected: svcPlan,
			Action: func(c *Client, t *testing.T) (any, error) {
				r := &resource.ServicePlanUpdate{
					Metadata: resource.Metadata{
						Labels: map[string]string{
							"key": "value",
						},
						Annotations: map[string]string{
							"note": "detailed information",
						},
					},
				}
				return c.ServicePlans.Update("79aae221-b2a6-4aaa-a134-76f605af46c9", r)
			},
		},
	}
	executeTests(tests, t)
}
