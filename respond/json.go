package respond

import (
	"encoding/json"

	"github.com/Jille/convreq/internal"
)

// RespondJSON creates a JSON response
func RespondJSON(data interface{}) internal.HttpResponse {
	response, err := json.Marshal(data)
	if err != nil {
		return InternalServerError(err.Error())
	}
	return WithHeader(Bytes(response), "Content-Type", "application/json")
}
