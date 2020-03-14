// Library respond creates convreq.HttpResponse objects. These can be returned from HTTP handlers to respond something to the client.
package respond

import (
	"io"
	"net/http"
	"strconv"

	internal "github.com/Jille/convreq/internal"
)

type httpError struct {
	code int
	msg  string
}

// Respond implements convreq.HttpResponse.
func (e httpError) Respond(w http.ResponseWriter, r *http.Request) error {
	http.Error(w, e.msg, e.code)
	return nil
}

// BadRequest creates a HTTP 400 Bad Request response.
func BadRequest(err error) internal.HttpResponse {
	return httpError{400, err.Error()}
}

// InternalServerError creates a HTTP 500 Internal Server Error response.
func InternalServerError(err error) internal.HttpResponse {
	return httpError{500, err.Error()}
}

type handlerResponse struct {
	h http.Handler
}

// Respond implements convreq.HttpResponse.
func (h handlerResponse) Respond(w http.ResponseWriter, r *http.Request) error {
	h.h.ServeHTTP(w, r)
	return nil
}

// RepondWithHandler creates a response that delegates to a regular http.Handler.
func RepondWithHandler(h http.Handler) internal.HttpResponse {
	return handlerResponse{h}
}

type withHeaders struct {
	parent internal.HttpResponse
	header http.Header
}

// Respond implements convreq.HttpResponse.
func (h withHeaders) Respond(w http.ResponseWriter, r *http.Request) error {
	for k, v := range h.header {
		w.Header()[k] = v
	}
	return nil
}

// WithHeader wraps a response and adds an additional header to be set.
func WithHeader(hr internal.HttpResponse, header, value string) internal.HttpResponse {
	ret := withHeaders{
		parent: hr,
	}
	ret.header.Set(header, value)
	return ret
}

// WithHeaders wraps a response and adds additional headers to be set.
func WithHeaders(hr internal.HttpResponse, headers http.Header) internal.HttpResponse {
	return withHeaders{
		parent: hr,
		header: headers,
	}
}

type redirect struct {
	code int
	url  string
}

// Respond implements convreq.HttpResponse.
func (re redirect) Respond(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Location", re.url)
	w.WriteHeader(re.code)
	return nil
}

// Redirect creates a response that redirects to another URL.
func Redirect(code int, url string) internal.HttpResponse {
	return redirect{code, url}
}

type respondString struct {
	data string
}

// Respond implements convreq.HttpResponse.
func (s respondString) Respond(w http.ResponseWriter, r *http.Request) error {
	if w.Header().Get("Content-Length") == "" {
		w.Header().Set("Content-Length", strconv.Itoa(len(s.data)))
	}
	_, err := io.WriteString(w, s.data)
	return err
}

// String creates a response that sends a string back to the client.
func String(data string) internal.HttpResponse {
	return respondString{data}
}
