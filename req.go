package yin

import (
	"encoding/json"
	"errors"
	"net/http"
)

type Request struct {
	Body     func(body interface{}) error
	Header   func(key string) string
	Query    func(key string) string
	Location func() *Location
}

func Req(r *http.Request) *Request {
	req := &Request{}

	req.Body = func(body interface{}) error {
		if r.Body == nil {
			return errors.New("No Request body found")
		}
		err := json.NewDecoder(r.Body).Decode(body)
		if err != nil {
			return err
		}
		return nil
	}

	req.Header = func(key string) string {
		return r.Header.Get(key)
	}

	req.Query = func(key string) string {
		return r.URL.Query().Get(key)
	}

	req.Location = func() *Location {
		return getLocation(r)
	}

	return req
}
