package web

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"reflect"
	"search/domains/users"
	"search/util"
)

func (s *HTTPServer) FindUser(w http.ResponseWriter, r *http.Request) {
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
	orgs, err := s.UsersFinder.Find(r.Context(), users.FindAllInput{Field: field, Value: value})
	if err != nil {
		http.Error(w, "invalid field or value.  Please refer to /user-details for a list of valid fields", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(orgs)
}

func (s *HTTPServer) ListUserFields(w http.ResponseWriter, r *http.Request) {
	val := reflect.ValueOf(users.User{})
	fields := util.GetJsonFields(val)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(fields)
}
