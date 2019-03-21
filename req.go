package yin

import (
	"encoding/json"
	"errors"
	"net/http"
)

type reqOpts struct {
	Body func(body interface{}) error
}

func Req(r *http.Request) *reqOpts {
	req := &reqOpts{}

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

	return req
}
