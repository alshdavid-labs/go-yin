package yin

import (
	"encoding/json"
	"net/http"
)

type response struct {
	Status     func(statusCode int) *response
	JSON       func(interface{})
	String     func(s string)
	SendStatus func(statusCode int)
}

func Res(w http.ResponseWriter) *response {
	r := &response{}

	r.Status = func(statusCode int) *response {
		w.WriteHeader(statusCode)
		return r
	}

	r.JSON = func(u interface{}) {
		json.NewEncoder(w).Encode(u)
	}

	r.String = func(s string) {
		w.Write([]byte(s))
	}

	r.SendStatus = func(statusCode int) {
		w.WriteHeader(statusCode)
		w.Write([]byte(""))
	}

	return r
}
