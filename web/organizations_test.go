package web

import (
	"context"
	"encoding/json"
	"github.com/tj/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"search/domains/organizations"
	"search/domains/users"
	"testing"
)

type mockOrgFinder struct {
}

func (m mockOrgFinder) Find(ctx context.Context, input organizations.FindAllInput) ([]organizations.Organization, error) {
	return []organizations.Organization{
		{
			ID:   101,
			Name: "Sulfax",
			DomainNames: []string{
				"comvey.com",
				"velity.com",
				"enormo.com",
			},
			CreatedAt:     "2016-01-12T01:16:06 -11:00",
			Details:       "MegaCörp",
			SharedTickets: true,
			Tags: []string{
				"Travis",
				"Clarke",
				"Glenn",
				"Santos",
			},
		},
	}, nil
}

func (m mockOrgFinder) ImportAll(ctx context.Context, input organizations.ImportAllInput) error {
	return nil
}

func TestHTTPServer_FindOrg(t *testing.T) {
	server := HTTPServer{
		OrgFinder: &mockOrgFinder{},
	}
	handler := MakeHandler(server)
	req, err := http.NewRequest("GET", "/organizations?field=id&value=101", nil)
	assert.Nil(t, err)

	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	resp := w.Result()
	body, err := ioutil.ReadAll(resp.Body)
	assert.Nil(t, err)

	var u []users.User
	err = json.Unmarshal(body, &u)
	assert.Nil(t, err)

	assert.Equal(t, "Sulfax", u[0].Name)
	assert.Len(t, u, 1)
}

func TestHTTPServer_FindOrg_MissingField(t *testing.T) {
	server := HTTPServer{
		OrgFinder: &mockOrgFinder{},
	}
	handler := MakeHandler(server)
	req, err := http.NewRequest("GET", "/organizations?value=101", nil)
	assert.Nil(t, err)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	assert.Equal(t, 404, w.Code)
}

func TestHTTPServer_ListOrgFields(t *testing.T) {
	server := HTTPServer{
		OrgFinder: &mockOrgFinder{},
	}
	handler := MakeHandler(server)
	req, err := http.NewRequest("GET", "/org-details", nil)
	assert.Nil(t, err)

	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
}

func TestHTTPServer_BadURL(t *testing.T) {
	server := HTTPServer{
		OrgFinder: &mockOrgFinder{},
	}
	handler := MakeHandler(server)
	req, err := http.NewRequest("GET", "/org", nil)
	assert.Nil(t, err)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	assert.Equal(t, 404, w.Code)
}
