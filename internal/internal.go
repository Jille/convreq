package convreq

import (
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
)

var decoder = schema.NewDecoder()

// HttpResponse is what is to be returned from request handlers.
// Respond gets executed to write the response to the client.
type HttpResponse interface {
	Respond(w http.ResponseWriter, r *http.Request) error
}

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
