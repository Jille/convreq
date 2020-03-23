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

package convreq_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/Jille/convreq"
	"github.com/Jille/convreq/respond"
)

type myStruct struct{}

func createMyStruct(w http.ResponseWriter, r *http.Request) (reflect.Value, convreq.HttpResponse) {
	return reflect.ValueOf(myStruct{}), nil
}

func TestWithParameterType(t *testing.T) {
	respRecorder := httptest.NewRecorder()
	var handler http.Handler = convreq.Wrap(func(ms myStruct) {}, convreq.WithParameterType(reflect.TypeOf(myStruct{}), createMyStruct))
	handler.ServeHTTP(respRecorder, httptest.NewRequest("GET", "/", nil))
	if respRecorder.Code != 200 {
		t.Errorf("got code %d; want %d", respRecorder.Code, 200)
	}
}

func handleCode(w http.ResponseWriter, r *http.Request, code int) {
	w.WriteHeader(code)
}

func TestWithReturnType(t *testing.T) {
	respRecorder := httptest.NewRecorder()
	var handler http.Handler = convreq.Wrap(func() int { return 204 }, convreq.WithReturnType(reflect.TypeOf(204), handleCode))
	handler.ServeHTTP(respRecorder, httptest.NewRequest("GET", "/", nil))
	if respRecorder.Code != 204 {
		t.Errorf("got code %d; want %d", respRecorder.Code, 204)
	}
}

func TestWithContextWrapper(t *testing.T) {
	respRecorder := httptest.NewRecorder()
	step := 0
	ctxCancel := func() {
		if step != 2 {
			t.Fatal("ordering issue: context destruction should be called last")
		}
		step++
	}
	ctxCreate := func(ctx context.Context) (context.Context, func()) {
		if step != 0 {
			t.Fatal("ordering issue: context creation should be called first")
		}
		step++
		return ctx, ctxCancel
	}
	h := func(ctx context.Context) {
		if step < 1 {
			t.Fatal("ordering issue: handler called before context creation")
		}
		if step > 1 {
			t.Fatal("ordering issue: handler called after context destruction")
		}
		step++
	}
	var handler http.Handler = convreq.Wrap(h, convreq.WithContextWrapper(ctxCreate))
	handler.ServeHTTP(respRecorder, httptest.NewRequest("GET", "/", nil))
	if respRecorder.Code != 200 {
		t.Errorf("got code %d; want %d", respRecorder.Code, 200)
	}
	if step != 3 {
		t.Fatal("ordering issue: context didn't get destroyed?")
	}
}

func TestWithErrorHandler(t *testing.T) {
	respRecorder := httptest.NewRecorder()
	eh := func(code int, msg string, r *http.Request) convreq.HttpResponse {
		if code != 500 {
			t.Errorf("Error handler got code %d; want %d", code, 500)
		}
		if msg != "test" {
			t.Errorf("Error handler got error %q; want %q", msg, "test")
		}
		return respond.OverrideResponseCode(respond.String("Bummer"), 501)
	}
	var handler http.Handler = convreq.Wrap(func() error { return errors.New("test") }, convreq.WithErrorHandler(eh))
	handler.ServeHTTP(respRecorder, httptest.NewRequest("GET", "/", nil))
	if respRecorder.Code != 501 {
		t.Errorf("got code %d; want %d", respRecorder.Code, 501)
	}
}
