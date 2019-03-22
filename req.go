package yin

import (
	"encoding/json"
	"errors"
	"net/http"
)

type request struct {
	Body      func(body interface{}) error
	GetHeader func(key string) string
	GetQuery  func(key string) string
}

func Req(r *http.Request) *request {
	req := &request{}

	req.Body = func(body interface{}) error {
		if r.Body == nil {
			return errors.New("No request body found")
		}
		err := json.NewDecoder(r.Body).Decode(body)
		if err != nil {
			return err
		}
		return nil
	}

	req.GetHeader = func(key string) string {
		return r.Header.Get(key)
	}

	req.GetQuery = func(key string) string {
		return r.Header.Get(key)
	}

	// req.GetParam = func(key string) string {
	// 	return r.Header.Get(key)
	// }

	return req
}
