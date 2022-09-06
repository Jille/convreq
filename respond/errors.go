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

package respond

import (
	"net/http"

	"github.com/Jille/convreq/internal"
)

type httpError struct {
	code int
	msg  string
}

// Respond implements convreq.HttpResponse.
func (e httpError) Respond(w http.ResponseWriter, r *http.Request) error {
	if f, ok := r.Context().Value(internal.ErrorHandlerContextKey).(internal.ErrorHandler); ok {
		return f(e.code, e.msg, r).Respond(w, r)
	}
	http.Error(w, e.msg, e.code)
	return nil
}

// Error creates a HTTP 500 Internal Server Error response.
func Error(err error) internal.HttpResponse {
	return httpError{500, err.Error()}
}

// Created creates a HTTP 201 Created response.
func Created(msg string) internal.HttpResponse {
	return httpError{201, msg}
}

// Accepted creates a HTTP 202 Accepted response.
func Accepted(msg string) internal.HttpResponse {
	return httpError{202, msg}
}

// NoContent creates a HTTP 204 No Content response.
// The param ignored exists to provide backwards compatibility with an old bug that accepted a message.
func NoContent(ignored ...string) internal.HttpResponse {
	return httpError{204, ""}
}

// ResetContent creates a HTTP 205 Reset Content response.
func ResetContent(msg string) internal.HttpResponse {
	return httpError{205, msg}
}

// NotModified creates a HTTP 304 Not Modified response.
func NotModified(msg string) internal.HttpResponse {
	return httpError{304, msg}
}

// BadRequest creates a HTTP 400 Bad Request response.
func BadRequest(msg string) internal.HttpResponse {
	return httpError{400, msg}
}

// Forbidden creates a HTTP 403 Forbidden response.
func Forbidden(msg string) internal.HttpResponse {
	return httpError{403, msg}
}

// NotFound creates a HTTP 404 Not Found response.
func NotFound(msg string) internal.HttpResponse {
	return httpError{404, msg}
}

// MethodNotAllowed creates a HTTP 405 Method Not Allowed response.
func MethodNotAllowed(msg string) internal.HttpResponse {
	return httpError{405, msg}
}

// NotAcceptable creates a HTTP 406 Not Acceptable response.
func NotAcceptable(msg string) internal.HttpResponse {
	return httpError{406, msg}
}

// RequestTimeout creates a HTTP 408 Request Timeout response.
func RequestTimeout(msg string) internal.HttpResponse {
	return httpError{408, msg}
}

// Conflict creates a HTTP 409 Conflict response.
func Conflict(msg string) internal.HttpResponse {
	return httpError{409, msg}
}

// Gone creates a HTTP 410 Gone response.
func Gone(msg string) internal.HttpResponse {
	return httpError{410, msg}
}

// LengthRequired creates a HTTP 411 Length Required response.
func LengthRequired(msg string) internal.HttpResponse {
	return httpError{411, msg}
}

// PreconditionFailed creates a HTTP 412 Precondition Failed response.
func PreconditionFailed(msg string) internal.HttpResponse {
	return httpError{412, msg}
}

// PayloadTooLarge creates a HTTP 413 Payload Too Large response.
func PayloadTooLarge(msg string) internal.HttpResponse {
	return httpError{413, msg}
}

// URITooLong creates a HTTP 414 URI Too Long response.
func URITooLong(msg string) internal.HttpResponse {
	return httpError{414, msg}
}

// UnsupportedMediaType creates a HTTP 415 Unsupported Media Type response.
func UnsupportedMediaType(msg string) internal.HttpResponse {
	return httpError{415, msg}
}

// RangeNotSatisfiable creates a HTTP 416 Range Not Satisfiable response.
func RangeNotSatisfiable(msg string) internal.HttpResponse {
	return httpError{416, msg}
}

// ExpectationFailed creates a HTTP 417 Expectation Failed response.
func ExpectationFailed(msg string) internal.HttpResponse {
	return httpError{417, msg}
}

// Imateapot creates a HTTP 418 I'm a teapot response.
func Imateapot(msg string) internal.HttpResponse {
	return httpError{418, msg}
}

// UnprocessableEntity creates a HTTP 422 Unprocessable Entity response.
func UnprocessableEntity(msg string) internal.HttpResponse {
	return httpError{422, msg}
}

// FailedDependency creates a HTTP 424 Failed Dependency response.
func FailedDependency(msg string) internal.HttpResponse {
	return httpError{424, msg}
}

// TooEarly creates a HTTP 425 Too Early response.
func TooEarly(msg string) internal.HttpResponse {
	return httpError{425, msg}
}

// UpgradeRequired creates a HTTP 426 Upgrade Required response.
func UpgradeRequired(msg string) internal.HttpResponse {
	return httpError{426, msg}
}

// PreconditionRequired creates a HTTP 428 Precondition Required response.
func PreconditionRequired(msg string) internal.HttpResponse {
	return httpError{428, msg}
}

// TooManyRequests creates a HTTP 429 Too Many Requests response.
func TooManyRequests(msg string) internal.HttpResponse {
	return httpError{429, msg}
}

// RequestHeaderFieldsTooLarge creates a HTTP 431 Request Header Fields Too Large response.
func RequestHeaderFieldsTooLarge(msg string) internal.HttpResponse {
	return httpError{431, msg}
}

// UnavailableForLegalReasons creates a HTTP 451 Unavailable For Legal Reasons response.
func UnavailableForLegalReasons(msg string) internal.HttpResponse {
	return httpError{451, msg}
}

// InternalServerError creates a HTTP 500 Internal Server Error response.
func InternalServerError(msg string) internal.HttpResponse {
	return httpError{500, msg}
}

// NotImplemented creates a HTTP 501 Not Implemented response.
func NotImplemented(msg string) internal.HttpResponse {
	return httpError{501, msg}
}

// BadGateway creates a HTTP 502 Bad Gateway response.
func BadGateway(msg string) internal.HttpResponse {
	return httpError{502, msg}
}

// ServiceUnavailable creates a HTTP 503 Service Unavailable response.
func ServiceUnavailable(msg string) internal.HttpResponse {
	return httpError{503, msg}
}

// GatewayTimeout creates a HTTP 504 Gateway Timeout response.
func GatewayTimeout(msg string) internal.HttpResponse {
	return httpError{504, msg}
}

// HTTPVersionNotSupported creates a HTTP 505 HTTP Version Not Supported response.
func HTTPVersionNotSupported(msg string) internal.HttpResponse {
	return httpError{505, msg}
}

// VariantAlsoNegotiates creates a HTTP 506 Variant Also Negotiates response.
func VariantAlsoNegotiates(msg string) internal.HttpResponse {
	return httpError{506, msg}
}

// InsufficientStorage creates a HTTP 507 Insufficient Storage response.
func InsufficientStorage(msg string) internal.HttpResponse {
	return httpError{507, msg}
}

// LoopDetected creates a HTTP 508 Loop Detected response.
func LoopDetected(msg string) internal.HttpResponse {
	return httpError{508, msg}
}

// NotExtended creates a HTTP 510 Not Extended response.
func NotExtended(msg string) internal.HttpResponse {
	return httpError{510, msg}
}

// NetworkAuthenticationRequired creates a HTTP 511 Network Authentication Required response.
func NetworkAuthenticationRequired(msg string) internal.HttpResponse {
	return httpError{511, msg}
}
