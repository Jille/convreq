package respond

import (
	"net/http"

	"github.com/Jille/convreq/internal"
)

// Func exposes the raw http.ResponseWriter and http.Request in a function that should return a convreq.HttpResponse.
type Func func(w http.ResponseWriter, r *http.Request) error

// Respond implements convreq.HttpResponse.
func (fn Func) Respond(w http.ResponseWriter, r *http.Request) error {
	return fn(w, r)
}

// Skip does nothing. Useful for when you have already responded yourself but still have to return a convreq.HttpResponse.
func Skip() internal.HttpResponse {
	return Func(func(w http.ResponseWriter, r *http.Request) error {
		return nil
	})
}
