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

func (r *ReactionsGETRequest) BodyRaw() []byte {
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

func (r *ReactionsPOSTRequest) BodyRaw() []byte {
	return makeJsonBodyRaw(r.Body)
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

func (r *ReactionsDELETERequest) BodyRaw() []byte {
	return makeJsonBodyRaw(r.Body)
}

func (r *ConfiguratinPOSTRequest) Method() string {
	return "POST"
}

func (r *ConfiguratinPOSTRequest) Path() string {
	return "/configuration"
}

func (r *ConfiguratinPOSTRequest) QueryString() string {
	return ""
}

func (r *ConfiguratinPOSTRequest) BodyRaw() []byte {
	return makeJsonBodyRaw(r.Body)
}

func (r *AvailableReactionsGETRequest) Method() string {
	return "GET"
}

func (r *AvailableReactionsGETRequest) Path() string {
	return "/configuration/available_reactions"
}

func (r *AvailableReactionsGETRequest) QueryString() string {
	return makeQueryString(r.Query)
}

func (r *AvailableReactionsGETRequest) BodyRaw() []byte {
	return nil
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

func makeJsonBodyRaw(body interface{}) []byte {
	res, err := json.Marshal(body)
	if err != nil {
		panic(err)
	}
	return res
}
