// Copyright 2020 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
