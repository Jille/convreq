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
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(rj.data)
}

// ServeJSON uses json.Marshal() to serve a JSON object.
func ServeJSON(data interface{}) internal.HttpResponse {
	return respondJSON{data}
}
