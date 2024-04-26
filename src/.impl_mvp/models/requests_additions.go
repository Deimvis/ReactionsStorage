package models

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
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

func (r *ReactionsGETRequest) Header() http.Header {
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

func (r *ReactionsPOSTRequest) Header() http.Header {
    return nil
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

func (r *ReactionsDELETERequest) Header() http.Header {
    return nil
}

func (r *ConfigurationPOSTRequest) Method() string {
	return "POST"
}

func (r *ConfigurationPOSTRequest) Path() string {
	return "/configuration"
}

func (r *ConfigurationPOSTRequest) QueryString() string {
	return ""
}

func (r *ConfigurationPOSTRequest) BodyRaw() []byte {
	return makeJsonBodyRaw(r.Body)
}

func (r *ConfigurationPOSTRequest) Header() http.Header {
    return r.Headers
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

func (r *AvailableReactionsGETRequest) Header() http.Header {
    return nil
}

func (r *NamespaceGETRequest) Method() string {
	return "GET"
}

func (r *NamespaceGETRequest) Path() string {
	return "/configuration/namespace"
}

func (r *NamespaceGETRequest) QueryString() string {
	return makeQueryString(r.Query)
}

func (r *NamespaceGETRequest) BodyRaw() []byte {
	return nil
}

func (r *NamespaceGETRequest) Header() http.Header {
    return nil
}

func makeQueryString(query interface{}) string {
	var sb strings.Builder
	v := reflect.ValueOf(query)
	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)
		if f.Kind() == reflect.Ptr && f.IsNil() {
			continue
		}
		queryKey := v.Type().Field(i).Tag.Get("query")
		queryValue := toString(f)
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

// converts to string
func toString(v reflect.Value) string {
	for v.Kind() == reflect.Ptr {
		v = getUnderlyingValue(v)
	}
	var s string
	switch v.Kind() {
	case reflect.String:
		s = v.String()
	case reflect.Bool:
		s = strconv.FormatBool(v.Bool())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		s = strconv.FormatInt(v.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		s = strconv.FormatUint(v.Uint(), 10)
	default:
		panic(fmt.Errorf("got unsuported kind: %s", v.Kind()))
	}
	return s
}

// returns dereferenced value (does nothing on non-pointer value)
func getUnderlyingValue(v reflect.Value) reflect.Value {
	if v.Kind() == reflect.Ptr {
		return reflect.Indirect(v)
	}
	return v
}
