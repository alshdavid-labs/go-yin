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
	File       func(filepath string)
	Redirect   func(statusCode int, url string)
}

func Res(w http.ResponseWriter, r *http.Request) *Response {
	res := &Response{}

	res.Status = func(statusCode int) *Response {
		w.WriteHeader(statusCode)
		return res
	}

	res.SetCookie = func(cookie *http.Cookie) *Response {
		http.SetCookie(w, cookie)
		return res
	}

	res.SetHeader = func(key string, value string) *Response {
		w.Header().Set(key, value)
		return res
	}

	res.JSON = func(u interface{}) {
		w.Header().Set(Headers.ContentType, "application/json")
		json.NewEncoder(w).Encode(u)
	}

	res.String = func(s string) {
		w.Write([]byte(s))
	}

	res.SendStatus = func(statusCode int) {
		w.WriteHeader(statusCode)
		w.Write([]byte(""))
	}

	res.File = func(filepath string) {
		http.ServeFile(w, r, filepath)
	}

	res.Redirect = func(statusCode int, url string) {
		http.Redirect(w, r, url, statusCode)
	}

	return res
}
