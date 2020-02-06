package web

import (
	"context"
	"encoding/json"
	"github.com/tj/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"search/domains/users"
	"testing"
)

type mockUserFinder struct {
}

func (m mockUserFinder) Find(ctx context.Context, input users.FindAllInput) ([]users.User, error) {
	return []users.User{
		{
			ID:             1,
			Url:            "http://initech.zendesk.com/api/v2/users/1.json",
			ExternalID:     "74341f74-9c79-49d5-9611-87ef9b6eb75f",
			Name:           "Francisca Rasmussen",
			Alias:          "Miss Coffey",
			CreatedAt:      "2016-04-15T05:19:46 -10:00",
			Active:         true,
			Verified:       true,
			Shared:         false,
			Locale:         "en-AU",
			Timezone:       "Sri Lanka",
			LastLoginAt:    "2013-08-04T01:03:27 -10:00",
			Email:          "coffeyrasmussen@flotonic.com",
			Phone:          "8335-422-718",
			Signature:      "Don't Worry Be Happy!",
			OrganizationID: 119,
			Tags:           []string{"Springville", "Sutton", "Hartsville/Hartley", "Diaperville"},
			Suspended:      true,
			Role:           "admin",
		},
	}, nil
}

func (m mockUserFinder) ImportAll(ctx context.Context, input users.ImportAllInput) error {
	return nil
}

func TestHTTPServer_FindUser(t *testing.T) {
	server := HTTPServer{
		UsersFinder: &mockUserFinder{},
	}
	handler := MakeHandler(server)
	req, err := http.NewRequest("GET", "/users?field=id&value=1", nil)
	assert.Nil(t, err)

	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	resp := w.Result()
	body, err := ioutil.ReadAll(resp.Body)
	assert.Nil(t, err)

	var users []users.User
	err = json.Unmarshal(body, &users)
	assert.Nil(t, err)

	assert.Len(t, users, 1)
	assert.Equal(t, "Francisca Rasmussen", users[0].Name)
}

func TestHTTPServer_FindUsers_MissingField(t *testing.T) {
	server := HTTPServer{
		OrgFinder: &mockOrgFinder{},
	}
	handler := MakeHandler(server)
	req, err := http.NewRequest("GET", "/users?value=101", nil)
	assert.Nil(t, err)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	assert.Equal(t, 404, w.Code)
}

func TestHTTPServer_ListUsersFields(t *testing.T) {
	server := HTTPServer{
		OrgFinder: &mockOrgFinder{},
	}
	handler := MakeHandler(server)
	req, err := http.NewRequest("GET", "/user-details", nil)
	assert.Nil(t, err)

	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
}

func TestHTTPServer_BadUsersURL(t *testing.T) {
	server := HTTPServer{
		OrgFinder: &mockOrgFinder{},
	}
	handler := MakeHandler(server)
	req, err := http.NewRequest("GET", "/users", nil)
	assert.Nil(t, err)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	assert.Equal(t, 404, w.Code)
}
