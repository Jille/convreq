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
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/Jille/convreq/internal"
)

// Template is a template-like object, like text/template.Template.
type Template interface {
	Execute(wr io.Writer, data interface{}) error
}

// RenderTemplate creates a response that will render the given template.
func RenderTemplate(tpl Template, data interface{}) internal.HttpResponse {
	var buf bytes.Buffer
	if err := tpl.Execute(&buf, data); err != nil {
		return Error(fmt.Errorf("failed to render template: %v", err))
	}
	return Reader(&buf)
}

type respondReader struct {
	r io.Reader
}

// Respond implements convreq.HttpResponse.
func (rr respondReader) Respond(w http.ResponseWriter, r *http.Request) error {
	if closer, ok := rr.r.(io.Closer); ok {
		defer closer.Close()
	}
	if w.Header().Get("Content-Length") == "" {
		if lenner, ok := rr.r.(interface{ Len() int }); ok {
			w.Header().Set("Content-Length", strconv.Itoa(lenner.Len()))
		}
	}
	if _, err := io.Copy(w, rr.r); err != nil {
		return fmt.Errorf("failed to write response to client: %v", err)
	}
	return nil
}

// Reader creates a response that copies everything from the reader to the client.
// If the reader is also an io.Closer, r will be closed.
// If the reader has a method Len() int, it will be sent as the Content-Length if that header isn't set yet.
func Reader(r io.Reader) internal.HttpResponse {
	return respondReader{r}
}
