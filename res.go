package yin

import (
	"encoding/json"
	"net/http"
)

type response struct {
	Status     func(statusCode int) *response
	SetCookie  func(cookie *http.Cookie) *response
	SetHeader  func(key string, value string) *response
	JSON       func(interface{})
	String     func(s string)
	SendStatus func(statusCode int)
	File       func(r *http.Request, filepath string)
}

func Res(w http.ResponseWriter) *response {
	r := &response{}

	r.Status = func(statusCode int) *response {
		w.WriteHeader(statusCode)
		return r
	}

	r.SetCookie = func(cookie *http.Cookie) *response {
		http.SetCookie(w, cookie)
		return r
	}

	r.SetHeader = func(key string, value string) *response {
		w.Header().Set(key, value)
		return r
	}

	r.JSON = func(u interface{}) {
		w.Header().Set(Headers.ContentType, "application/json")
		json.NewEncoder(w).Encode(u)
	}

	r.String = func(s string) {
		w.Write([]byte(s))
	}

	r.SendStatus = func(statusCode int) {
		w.WriteHeader(statusCode)
		w.Write([]byte(""))
	}

	r.File = func(r *http.Request, filepath string) {
		http.ServeFile(w, r, filepath)
	}

	return r
}