package util

import (
	"reflect"
	s "strings"
	"time"
)

//Join joins together a slice of strings
func Join(items []string) string {
	return s.Join(items, ",")
}

//Split returns a slice of strings
func Split(items string) []string {
	return s.Split(items, ",")
}

//ToTime transfer the string to a time object
func ToTime(str string) time.Time {
	layout := "2016-05-21T11:10:28 -10:00"
	t, err := time.Parse(layout, str)
	if err != nil {
		return time.Now()
	}
	return t
}

//GetJsonFields are the fields that a user can query for in a specific domain
func GetJsonFields(val reflect.Value) []string {
	fields := make([]string, 0)
	for i := 0; i < val.Type().NumField(); i++ {
		field := val.Type().Field(i).Tag.Get("json")
		//todo fix this
		if field == "_id" {
			field = "id"
		}
		fields = append(fields, field)
	}
	return fields
}
