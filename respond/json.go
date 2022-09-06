package respond

import (
	"encoding/json"
	"net/http"

	"github.com/Jille/convreq/internal"
)

type respondJSON struct {
	data interface{}
}

// Respond implements convreq.HttpResponse.
func (rj respondJSON) Respond(w http.ResponseWriter, r *http.Request) error {
	if w.Header().Get("Content-Type") == "" {
		w.Header().Set("Content-Type", "application/json")
	}
	return json.NewEncoder(w).Encode(rj.data)
}

// JSON marshals the given data and sends it to the requester.
func JSON(data interface{}) internal.HttpResponse {
	return respondJSON{data}
}

// ServeJSON marshals the given data and sends it to the requester.
// Deprecated: Use JSON() instead.
func ServeJSON(data interface{}) internal.HttpResponse {
	return JSON(data)
}
