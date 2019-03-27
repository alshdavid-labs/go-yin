package yin

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Status     func(statusCode int) *Response
	SetCookie  func(cookie *http.Cookie) *Response
	SetHeader  func(key string, value string) *Response
	JSON       func(interface{})
	String     func(s string)
	SendStatus func(statusCode int)
	File       func(r *http.Request, filepath string)
	Redirect   func(statusCode int, url string)
}

func Res(w http.ResponseWriter) *Response {
	r := &Response{}

	r.Status = func(statusCode int) *Response {
		w.WriteHeader(statusCode)
		return r
	}

	r.SetCookie = func(cookie *http.Cookie) *Response {
		http.SetCookie(w, cookie)
		return r
	}

	r.SetHeader = func(key string, value string) *Response {
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

	r.Redirect = func(statusCode int, url string) {
		w.WriteHeader(statusCode)
		w.Header().Set(Headers.Location, url)
		w.Write([]byte(""))
	}

	return r
}
