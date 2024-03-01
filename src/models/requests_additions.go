package models

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

func (r *ReactionsGETRequest) Method() string {
	return "GET"
}

func (r *ReactionsGETRequest) Path() string {
	return "/reactions"
}

func (r *ReactionsGETRequest) QueryString() string {
	return makeQueryString(r.Query)
}

func (r *ReactionsGETRequest) BodyJSON() []byte {
	return nil
}

func (r *ReactionsPOSTRequest) Method() string {
	return "POST"
}

func (r *ReactionsPOSTRequest) Path() string {
	return "/reactions"
}

func (r *ReactionsPOSTRequest) QueryString() string {
	return makeQueryString(r.Query)
}

func (r *ReactionsPOSTRequest) BodyJSON() []byte {
	return makeBodyJSON(r.Body)
}

func (r *ReactionsDELETERequest) Method() string {
	return "DELETE"
}

func (r *ReactionsDELETERequest) Path() string {
	return "/reactions"
}

func (r *ReactionsDELETERequest) QueryString() string {
	return ""
}

func (r *ReactionsDELETERequest) BodyJSON() []byte {
	return makeBodyJSON(r.Body)
}

func makeQueryString(query interface{}) string {
	var sb strings.Builder
	v := reflect.ValueOf(query)
	for i := 0; i < v.NumField(); i++ {
		if v.Field(i).IsNil() {
			continue
		}
		queryKey := v.Type().Field(i).Tag.Get("query")
		queryValue := v.Field(i).String()
		sb.WriteString(fmt.Sprintf("%s=%s", queryKey, queryValue))
		if i < v.NumField()-1 {
			sb.WriteString("&")
		}
	}
	return sb.String()
}

func makeBodyJSON(body interface{}) []byte {
	res, err := json.Marshal(body)
	if err != nil {
		panic(err)
	}
	return res
}
