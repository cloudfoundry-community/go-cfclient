package client

import (
	"github.com/cloudfoundry-community/go-cfclient/v3/resource"
	"github.com/cloudfoundry-community/go-cfclient/v3/test"
	"net/http"
	"testing"
)

func TestIsolationSegments(t *testing.T) {
	g := test.NewObjectJSONGenerator(1)
	iso := g.IsolationSegment()
	iso2 := g.IsolationSegment()
	iso3 := g.IsolationSegment()
	iso4 := g.IsolationSegment()
	isoRelations := g.IsolationSegmentRelationships()

	tests := []RouteTest{
		{
			Description: "Create isolation segment",
			Route: MockRoute{
				Method:   "POST",
				Endpoint: "/v3/isolation_segments",
				Output:   []string{iso},
				Status:   http.StatusCreated,
				PostForm: `{ "name": "my-iso" }`,
			},
			Expected: iso,
			Action: func(c *Client, t *testing.T) (any, error) {
				r := resource.NewIsolationSegmentCreate("my-iso")
				return c.IsolationSegments.Create(r)
			},
		},
		{
			Description: "Delete iso",
			Route: MockRoute{
				Method:   "DELETE",
				Endpoint: "/v3/isolation_segments/a45d5da8-67dc-4523-b34b-ffa68b8d8821",
				Status:   http.StatusAccepted,
			},
			Action: func(c *Client, t *testing.T) (any, error) {
				return nil, c.IsolationSegments.Delete("a45d5da8-67dc-4523-b34b-ffa68b8d8821")
			},
		},
		{
			Description: "Entitle isolation segment for org",
			Route: MockRoute{
				Method:   "POST",
				Endpoint: "/v3/isolation_segments/a45d5da8-67dc-4523-b34b-ffa68b8d8821/relationships/organizations",
				Output:   []string{isoRelations},
				Status:   http.StatusCreated,
				PostForm: `{ "data": [{ "guid":"5700e458-283d-4528-806f-c3509e038f05" }]}`,
			},
			Expected: isoRelations,
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.IsolationSegments.EntitleOrg("a45d5da8-67dc-4523-b34b-ffa68b8d8821", "5700e458-283d-4528-806f-c3509e038f05")
			},
		},
		{
			Description: "Get iso",
			Route: MockRoute{
				Method:   "GET",
				Endpoint: "/v3/isolation_segments/a45d5da8-67dc-4523-b34b-ffa68b8d8821",
				Output:   []string{iso},
				Status:   http.StatusOK},
			Expected: iso,
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.IsolationSegments.Get("a45d5da8-67dc-4523-b34b-ffa68b8d8821")
			},
		},
		{
			Description: "List all isolation segments",
			Route: MockRoute{
				Method:   "GET",
				Endpoint: "/v3/isolation_segments",
				Output:   g.Paged([]string{iso, iso2}, []string{iso3, iso4}),
				Status:   http.StatusOK},
			Expected: g.Array(iso, iso2, iso3, iso4),
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.IsolationSegments.ListAll(nil)
			},
		},
		{
			Description: "List all isolation segment related orgs",
			Route: MockRoute{
				Method:   "GET",
				Endpoint: "/v3/isolation_segments/a45d5da8-67dc-4523-b34b-ffa68b8d8821/relationships/organizations",
				Output: []string{`{
				  "data": [
					{
					  "guid": "68d54d31-9b3a-463b-ba94-e8e4c32edbac"
					},
					{
					  "guid": "b19f6525-cbd3-4155-b156-dc0c2a431b4c"
					}
				  ],
				  "links": {
					"self": {
					  "href": "https://api.example.org/v3/isolation_segments/bdeg4371-cbd3-4155-b156-dc0c2a431b4c/relationships/organizations"
					},
					"related": {
					  "href": "https://api.example.org/v3/isolation_segments/bdeg4371-cbd3-4155-b156-dc0c2a431b4c/organizations"
					}
				  }
				}`},
				Status: http.StatusOK,
			},
			Expected: `["68d54d31-9b3a-463b-ba94-e8e4c32edbac", "b19f6525-cbd3-4155-b156-dc0c2a431b4c"]`,
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.IsolationSegments.ListOrgRelationships("a45d5da8-67dc-4523-b34b-ffa68b8d8821")
			},
		},
		{
			Description: "List all isolation segment related spaces",
			Route: MockRoute{
				Method:   "GET",
				Endpoint: "/v3/isolation_segments/a45d5da8-67dc-4523-b34b-ffa68b8d8821/relationships/spaces",
				Output: []string{`{
				  "data": [
					{
					  "guid": "885735b5-aea4-4cf5-8e44-961af0e41920"
					},
					{
					  "guid": "d4c91047-7b29-4fda-b7f9-04033e5c9c9f"
					}
				  ],
				  "links": {
					"self": {
					  "href": "https://api.example.org/v3/isolation_segments/bdeg4371-cbd3-4155-b156-dc0c2a431b4c/relationships/organizations"
					},
					"related": {
					  "href": "https://api.example.org/v3/isolation_segments/bdeg4371-cbd3-4155-b156-dc0c2a431b4c/organizations"
					}
				  }
				}`},
				Status: http.StatusOK,
			},
			Expected: `["885735b5-aea4-4cf5-8e44-961af0e41920", "d4c91047-7b29-4fda-b7f9-04033e5c9c9f"]`,
			Action: func(c *Client, t *testing.T) (any, error) {
				return c.IsolationSegments.ListSpaceRelationships("a45d5da8-67dc-4523-b34b-ffa68b8d8821")
			},
		},
		{
			Description: "Revoke isolation segment for org",
			Route: MockRoute{
				Method:   "DELETE",
				Endpoint: "/v3/isolation_segments/a45d5da8-67dc-4523-b34b-ffa68b8d8821/relationships/organizations/5700e458-283d-4528-806f-c3509e038f05",
				Status:   http.StatusNoContent,
			},
			Action: func(c *Client, t *testing.T) (any, error) {
				err := c.IsolationSegments.RevokeOrg("a45d5da8-67dc-4523-b34b-ffa68b8d8821", "5700e458-283d-4528-806f-c3509e038f05")
				return nil, err
			},
		},
		{
			Description: "Update iso",
			Route: MockRoute{
				Method:   "PATCH",
				Endpoint: "/v3/isolation_segments/a45d5da8-67dc-4523-b34b-ffa68b8d8821",
				Output:   []string{iso},
				Status:   http.StatusOK,
				PostForm: `{ "name": "new-name" }`,
			},
			Expected: iso,
			Action: func(c *Client, t *testing.T) (any, error) {
				name := "new-name"
				r := &resource.IsolationSegmentUpdate{
					Name: &name,
				}
				return c.IsolationSegments.Update("a45d5da8-67dc-4523-b34b-ffa68b8d8821", r)
			},
		},
	}
	executeTests(tests, t)
}
