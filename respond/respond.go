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

// Library respond creates convreq.HttpResponse objects. These can be returned from HTTP handlers to respond something to the client.
package respond

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Jille/convreq/internal"
)

type handlerResponse struct {
	h http.Handler
}

// Respond implements convreq.HttpResponse.
func (h handlerResponse) Respond(w http.ResponseWriter, r *http.Request) error {
	h.h.ServeHTTP(w, r)
	return nil
}

// Handler creates a response that delegates to a regular http.Handler.
func Handler(h http.Handler) internal.HttpResponse {
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
	return h.parent.Respond(w, r)
}

// WithHeader wraps a response and adds an additional header to be set.
func WithHeader(hr internal.HttpResponse, header, value string) internal.HttpResponse {
	ret := withHeaders{
		parent: hr,
		header: http.Header{},
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

type respondBytes struct {
	data []byte
}

// Respond implements convreq.HttpResponse.
func (b respondBytes) Respond(w http.ResponseWriter, r *http.Request) error {
	if w.Header().Get("Content-Length") == "" {
		w.Header().Set("Content-Length", strconv.Itoa(len(b.data)))
	}
	_, err := w.Write(b.data)
	return err
}

// String creates a response that sends a string back to the client.
func String(data string) internal.HttpResponse {
	return respondBytes{[]byte(data)}
}

// Bytes creates a response that sends bytes back to the client.
func Bytes(data []byte) internal.HttpResponse {
	return respondBytes{data}
}

// Printf creates a response that sends a formatted string back to the client.
func Printf(format string, v ...interface{}) internal.HttpResponse {
	return respondBytes{[]byte(fmt.Sprintf(format, v...))}
}
