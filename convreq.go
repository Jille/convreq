// Library convreq is a library to make HTTP requests more convenient.
package convreq

import (
	"github.com/Jille/convreq/internal"
)

// HttpResponse is what is to be returned from request handlers.
// Respond gets executed to write the response to the client.
type HttpResponse = internal.HttpResponse
