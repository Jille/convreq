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
	"io"
	"net/http"
	"time"

	"github.com/Jille/convreq/internal"
)

type respondFile struct {
	name string
}

// Respond implements convreq.HttpResponse.
func (rf respondFile) Respond(w http.ResponseWriter, r *http.Request) error {
	http.ServeFile(w, r, rf.name)
	return nil
}

// ServeFile uses http.ServeFile() to serve a file.
func ServeFile(name string) internal.HttpResponse {
	return respondFile{name}
}

type respondContent struct {
	name    string
	modtime time.Time
	content io.ReadSeeker
}

// Respond implements convreq.HttpResponse.
func (rc respondContent) Respond(w http.ResponseWriter, r *http.Request) error {
	http.ServeContent(w, r, rc.name, rc.modtime, rc.content)
	return nil
}

// ServeContent uses http.ServeContent() to serve content.
func ServeContent(name string, modtime time.Time, content io.ReadSeeker) internal.HttpResponse {
	return respondContent{name, modtime, content}
}
