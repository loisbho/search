package web

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"reflect"
	"search/domains/tickets"
	"search/util"
)

func (s *HTTPServer) FindTicket(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	field, hasField := params["field"]
	if !hasField {
		http.NotFound(w, r)
		return
	}
	value, hasValue := params["value"]
	if !hasValue {
		http.NotFound(w, r)
		return
	}
	orgs, err := s.TicketsFinder.Find(r.Context(), tickets.FindAllInput{Field: field, Value: value})
	if err != nil {
		http.Error(w, "invalid field or value.  Please refer to /ticket-details for a list of valid fields", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(orgs)
}

func (s *HTTPServer) ListTicketFields(w http.ResponseWriter, r *http.Request) {
	val := reflect.ValueOf(tickets.Ticket{})
	fields := util.GetJsonFields(val)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(fields)
}
