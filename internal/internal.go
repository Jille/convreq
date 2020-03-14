package internal

import (
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
)

var decoder = func() *schema.Decoder {
	d := schema.NewDecoder()
	d.IgnoreUnknownKeys(true)
	return d
}()

// HttpResponse is what is to be returned from request handlers.
// Respond gets executed to write the response to the client.
type HttpResponse interface {
	Respond(w http.ResponseWriter, r *http.Request) error
}

type ctxKey int

// ErrorHandlerContextKey is used to store an ErrorHandler in the context.
var ErrorHandlerContextKey ctxKey = 1

// ErrorHandler is a callback type that you can register with ContextWithErrorHandler or WithErrorHandler to have your own callback called to render errors.
type ErrorHandler func(code int, msg string, r *http.Request) HttpResponse

// DoRespond executes a HttpResponse and has it write to the ResponseWriter.
func DoRespond(w http.ResponseWriter, r *http.Request, hr HttpResponse) {
	if err := hr.Respond(w, r); err != nil {
		log.Printf("Failed to respond to request: %v", err)
	}
}

// DecodeGet parses the GET parameters of the request into `ret` using github.com/gorilla/schema.
func DecodeGet(r *http.Request, ret interface{}) error {
	vm, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		return fmt.Errorf("failed to parse query: %v", err)
	}
	for k, v := range mux.Vars(r) {
		vm.Set(k, v)
	}
	if err := decoder.Decode(ret, vm); err != nil {
		return fmt.Errorf("failed to parse url/query: %v", err)
	}
	return nil
}

// DecodePost parses the POST parameters of the request into `ret` using github.com/gorilla/schema.
func DecodePost(r *http.Request, ret interface{}) error {
	if err := r.ParseForm(); err != nil {
		return fmt.Errorf("failed to parse form input: %v", err)
	}
	if err := decoder.Decode(ret, r.PostForm); err != nil {
		return fmt.Errorf("failed to parse form input: %v", err)
	}
	return nil
}
