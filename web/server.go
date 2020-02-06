package web

import (
	"github.com/gorilla/mux"
	"net/http"
	"search/domains/organizations"
	"search/domains/tickets"
	"search/domains/users"
)

type HTTPServer struct {
	OrgFinder     organizations.Finder
	TicketsFinder tickets.Finder
	UsersFinder   users.Finder
}

func MakeHandler(server HTTPServer) http.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/organizations", server.FindOrg).Queries("field", "{field}").Queries("value", "{value}").Methods("GET")
	r.HandleFunc("/users", server.FindUser).Queries("field", "{field}").Queries("value", "{value}").Methods("GET")
	r.HandleFunc("/tickets", server.FindTicket).Queries("field", "{field}").Queries("value", "{value}").Methods("GET")
	r.HandleFunc("/org-details", server.ListOrgFields).Methods("GET")
	r.HandleFunc("/ticket-details", server.ListTicketFields).Methods("GET")
	r.HandleFunc("/user-details", server.ListUserFields).Methods("GET")
	return r
}
