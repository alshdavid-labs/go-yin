package yin

import (
	"net/http"

	newrelic "github.com/newrelic/go-agent"
)

func NewRelic(app newrelic.Application) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			txn := app.StartTransaction(r.URL.Path, w, r)
			defer txn.End()

			r = newrelic.RequestWithTransactionContext(r, txn)

			next.ServeHTTP(txn, r)
		})
	}
}
