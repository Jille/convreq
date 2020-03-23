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

type modifyingResponseWriter struct {
	w           http.ResponseWriter
	code        int
	codeWritten bool
}

// Header implements http.ResponseWriter.
func (w *modifyingResponseWriter) Header() http.Header {
	return w.w.Header()
}

// Write implements http.ResponseWriter.
func (w *modifyingResponseWriter) Write(b []byte) (int, error) {
	if !w.codeWritten {
		w.WriteHeader(w.code)
	}
	return w.w.Write(b)
}

// WriteHeader implements http.ResponseWriter.
func (w *modifyingResponseWriter) WriteHeader(statusCode int) {
	w.w.WriteHeader(w.code)
	w.codeWritten = true
}

type modifyingResponse struct {
	parent internal.HttpResponse
	code   int
}

// Respond implements convreq.HttpResponse.
func (m modifyingResponse) Respond(w http.ResponseWriter, r *http.Request) error {
	mw := &modifyingResponseWriter{
		w:    w,
		code: m.code,
	}
	return m.parent.Respond(mw, r)
}

// OverrideResponseCode wraps a response to override the status code with the one given to this function.
// Note that this replaces the original ResponseWriter with an internal one, possibly not implementing interfaces like Flusher.
func OverrideResponseCode(hr internal.HttpResponse, code int) internal.HttpResponse {
	return modifyingResponse{
		parent: hr,
		code:   code,
	}
}
