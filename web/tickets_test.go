package web

import (
	"context"
	"encoding/json"
	"github.com/tj/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"search/domains/organizations"
	"search/domains/tickets"
	"testing"
)

type mockTicketFinder struct {
}

func (m mockTicketFinder) Find(ctx context.Context, input tickets.FindAllInput) ([]tickets.Ticket, error) {
	return []tickets.Ticket{
		{
			ID:             "436bf9b0-1147-4c0a-8439-6f79833bff5b",
			Url:            "http://initech.zendesk.com/api/v2/tickets/436bf9b0-1147-4c0a-8439-6f79833bff5b.json",
			ExternalID:     "9210cdc9-4bee-485f-a078-35396cd74063",
			CreatedAt:      "2016-04-28T11:19:34 -10:00",
			Type:           "incident",
			Subject:        "A Catastrophe in Korea (North)",
			Description:    "Nostrud ad sit velit cupidatat laboris ipsum nisi amet laboris ex exercitation amet et proident. Ipsum fugiat aute dolore tempor nostrud velit ipsum.",
			Priority:       "low",
			Status:         "pending",
			SubmitterID:    38,
			AssigneeID:      24,
			OrganizationID: 116,
			Tags:           []string{"Ohio", "Pennsylvania", "American Samoa", "Northern Mariana Islands"},
			HasIncidents:   false,
			DueAt:          "2016-07-31T02:37:50 -10:00",
			Via:            "web",
		},
	}, nil
}

func (m mockTicketFinder) ImportAll(ctx context.Context, input organizations.ImportAllInput) error {
	return nil
}

func TestHTTPServer_FindTicket(t *testing.T) {
	server := HTTPServer{
		TicketsFinder: &mockTicketFinder{},
	}
	handler := MakeHandler(server)
	req, err := http.NewRequest("GET", "/tickets?field=id&value=436bf9b0-1147-4c0a-8439-6f79833bff5b", nil)
	assert.Nil(t, err)

	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	resp := w.Result()
	body, err := ioutil.ReadAll(resp.Body)
	assert.Nil(t, err)

	var tic []tickets.Ticket
	err = json.Unmarshal(body, &tic)
	assert.Nil(t, err)

	assert.Len(t, tic, 1)
	assert.Equal(t, int32(116), tic[0].OrganizationID)
}

func TestHTTPServer_FindTickets_MissingField(t *testing.T) {
	server := HTTPServer{
		OrgFinder: &mockOrgFinder{},
	}
	handler := MakeHandler(server)
	req, err := http.NewRequest("GET", "/tickets?value=101", nil)
	assert.Nil(t, err)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	assert.Equal(t, 404, w.Code)
}

func TestHTTPServer_ListTicketFields(t *testing.T) {
	server := HTTPServer{
		OrgFinder: &mockOrgFinder{},
	}
	handler := MakeHandler(server)
	req, err := http.NewRequest("GET", "/ticket-details", nil)
	assert.Nil(t, err)

	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
}

func TestHTTPServer_BadTicketURL(t *testing.T) {
	server := HTTPServer{
		OrgFinder: &mockOrgFinder{},
	}
	handler := MakeHandler(server)
	req, err := http.NewRequest("GET", "/ticket", nil)
	assert.Nil(t, err)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	assert.Equal(t, 404, w.Code)
}
