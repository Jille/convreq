// Library convreq is a library to make HTTP requests more convenient.
package convreq

import (
	"context"

	"github.com/Jille/convreq/internal"
)

// HttpResponse is what is to be returned from request handlers.
// Respond gets executed to write the response to the client.
type HttpResponse = internal.HttpResponse

// ErrorHandler is a callback type that you can register with ContextWithErrorHandler or WithErrorHandler to have your own callback called to render errors.
type ErrorHandler = internal.ErrorHandler

// ContextWithErrorHandler returns a new context within which all errors are rendered with ErrorHandler.
func ContextWithErrorHandler(ctx context.Context, f ErrorHandler) context.Context {
	return context.WithValue(ctx, internal.ErrorHandlerContextKey, f)
}
