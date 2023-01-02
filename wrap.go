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

package convreq

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/Jille/convreq/internal"
	"github.com/Jille/convreq/respond"
)

// extractor is a function that extracts one specific type from a ResponseWriter or Request.
// An example of an extractor is getContext, that retrieves the context from Request.
type extractor func(w http.ResponseWriter, r *http.Request) (reflect.Value, HttpResponse)

var extractorMap = map[reflect.Type]extractor{
	// The best way to get a reflect.Type of an interface seems to be to get a pointer to it and dereference it.
	reflect.TypeOf((*context.Context)(nil)).Elem():     getContext,
	reflect.TypeOf(&http.Request{}):                    getRequest,
	reflect.TypeOf((*http.ResponseWriter)(nil)).Elem(): getResponseWriter,
}

// handlers are function that can handle the return value from a request handler.
var handlerMap = map[reflect.Type]reflect.Value{
	reflect.TypeOf((*error)(nil)).Elem():        reflect.ValueOf(handleError),
	reflect.TypeOf((*HttpResponse)(nil)).Elem(): reflect.ValueOf(internal.DoRespond),
}

type wrapOptions struct {
	extractors      map[reflect.Type]extractor
	handlers        map[reflect.Type]reflect.Value
	contextWrappers []func(ctx context.Context) (context.Context, func())
}

// WrapOption can be given to Wrap to modify behavior.
type WrapOption func(wo *wrapOptions)

// WithParameterType makes Wrap() understand an extra type for request handler signatures.
// The extractor is a function that can derive the requested type from a ResponseWriter and Request.
func WithParameterType(t reflect.Type, e extractor) WrapOption {
	return func(wo *wrapOptions) {
		wo.extractors[t] = e
	}
}

// WithReturnType makes Wrap() understand an extra return type for request handler signatures.
// You should pass in a function that takes `http.ResponseWriter, *http.Request, T` where T is the type you've passed in as t.
func WithReturnType(t reflect.Type, f interface{}) WrapOption {
	return func(wo *wrapOptions) {
		wo.handlers[t] = reflect.ValueOf(f)
	}
}

// WithContextWrapper allows you to replace the context for the request.
// f is called just before the request gets handled, and the cancel function is called after the request is finished.
// The cancel function may be nil.
func WithContextWrapper(f func(ctx context.Context) (context.Context, func())) WrapOption {
	return func(wo *wrapOptions) {
		wo.contextWrappers = append(wo.contextWrappers, f)
	}
}

// WithErrorHandler can be passed on Wrap() to set an ErrorHandler for requests.
func WithErrorHandler(f ErrorHandler) WrapOption {
	return WithContextWrapper(func(ctx context.Context) (context.Context, func()) {
		return ContextWithErrorHandler(ctx, f), nil
	})
}

// Wrap takes a request handler function and returns a http.HandlerFunc for use with net/http.
// The given handler is expected to take arguments like context.Context, *http.Request and return a convreq.HttpResponse or an error.
func Wrap(f interface{}, opts ...WrapOption) http.HandlerFunc {
	t := reflect.TypeOf(f)
	v := reflect.ValueOf(f)

	wo := wrapOptions{
		extractors: map[reflect.Type]extractor{},
		handlers:   map[reflect.Type]reflect.Value{},
	}
	for t, e := range extractorMap {
		wo.extractors[t] = e
	}
	for t, f := range handlerMap {
		wo.handlers[t] = f
	}
	for _, o := range opts {
		o(&wo)
	}

	if t.Kind() != reflect.Func {
		panic(fmt.Errorf("convreq: %s: not a function", v.String()))
	}

	// Look up all input parameters, and look up how we can create that type based.
	ins := make([]extractor, t.NumIn())
	for i := 0; t.NumIn() > i; i++ {
		if e, ok := wo.extractors[t.In(i)]; ok {
			ins[i] = e
		} else if strings.HasSuffix(t.In(i).Name(), "Get") {
			ins[i] = createGetInput(t.In(i))
		} else if t.In(i).Kind() == reflect.Ptr && strings.HasSuffix(t.In(i).Elem().Name(), "Post") {
			ins[i] = createPostInput(t.In(i))
		} else if t.In(i).Kind() == reflect.Ptr && strings.HasSuffix(t.In(i).Elem().Name(), "JSON") {
			ins[i] = createJSONInput(t.In(i), true)
		} else if strings.HasSuffix(t.In(i).Name(), "JSON") {
			ins[i] = createJSONInput(t.In(i), false)
		}
		if ins[i] == nil {
			panic(fmt.Errorf("convreq: %s: don't know how to produce %s", v.String(), t.In(i).String()))
		}
	}
	if t.IsVariadic() {
		panic(fmt.Errorf("convreq: %s: can't use variadic functions", v.String()))
	}

	// Figure out how to handle return values.
	var handler reflect.Value
	switch t.NumOut() {
	case 0:
		handler = reflect.ValueOf(handleVoid)
	case 1:
		handler = wo.handlers[t.Out(0)]
	}
	if !handler.IsValid() {
		panic(fmt.Errorf("convreq: %s: don't know how to handle return type(s)", v.String()))
	}

	// We've done all the prework we can. We try to minimize the things on the request path.

	return func(w http.ResponseWriter, r *http.Request) {
		if len(wo.contextWrappers) > 0 {
			ctx := r.Context()
			var cancel func()
			for _, cw := range wo.contextWrappers {
				ctx, cancel = cw(ctx)
				if cancel != nil {
					defer cancel()
				}
			}
			r = r.WithContext(ctx)
		}
		var hr HttpResponse
		// Now that we're called, extract all input parameters from w and r.
		in := make([]reflect.Value, len(ins))
		for i, e := range ins {
			in[i], hr = e(w, r)
			if hr != nil {
				internal.DoRespond(w, r, hr)
				return
			}
		}
		// Call the user's handler function.
		outs := v.Call(in)
		// Handle the return value.
		handler.Call(append([]reflect.Value{reflect.ValueOf(w), reflect.ValueOf(r)}, outs...))
	}
}

// === Below are some functions that extract something from the http.Request and return a reflect.Value with that value.
// === Their results will be passed into request handlers.

func getContext(w http.ResponseWriter, r *http.Request) (reflect.Value, HttpResponse) {
	return reflect.ValueOf(r.Context()), nil
}

func getRequest(w http.ResponseWriter, r *http.Request) (reflect.Value, HttpResponse) {
	return reflect.ValueOf(r), nil
}

func getResponseWriter(w http.ResponseWriter, r *http.Request) (reflect.Value, HttpResponse) {
	return reflect.ValueOf(w), nil
}

func createGetInput(t reflect.Type) func(w http.ResponseWriter, r *http.Request) (reflect.Value, HttpResponse) {
	return func(w http.ResponseWriter, r *http.Request) (reflect.Value, HttpResponse) {
		// TODO(quis): Consider putting v in a sync.Pool.
		v := reflect.New(t)
		if err := internal.DecodeGet(r, v.Interface()); err != nil {
			return reflect.Value{}, respond.BadRequest(err.Error())
		}
		return v.Elem(), nil
	}
}

func createPostInput(pt reflect.Type) func(w http.ResponseWriter, r *http.Request) (reflect.Value, HttpResponse) {
	nilptr := reflect.New(pt).Elem()
	t := pt.Elem()
	return func(w http.ResponseWriter, r *http.Request) (reflect.Value, HttpResponse) {
		if r.Method != "POST" {
			return nilptr, nil
		}
		// TODO(quis): Consider putting v in a sync.Pool.
		v := reflect.New(t)
		if err := internal.DecodePost(r, v.Interface()); err != nil {
			return reflect.Value{}, respond.BadRequest(err.Error())
		}
		return v, nil
	}
}

func createJSONInput(pt reflect.Type, isPtr bool) func(w http.ResponseWriter, r *http.Request) (reflect.Value, HttpResponse) {
	t := pt
	if isPtr {
		t = pt.Elem()
	}
	return func(w http.ResponseWriter, r *http.Request) (reflect.Value, HttpResponse) {
		v := reflect.New(t)
		if err := internal.DecodeJSON(r, v.Interface()); err != nil {
			return reflect.Value{}, respond.BadRequest(err.Error())
		}
		if isPtr {
			return v, nil
		}
		return v.Elem(), nil
	}
}

// === Below are some functions that can handle the return value of a request handler.

func handleVoid(w http.ResponseWriter, r *http.Request) {
}

func handleError(w http.ResponseWriter, r *http.Request, err error) {
	if err != nil {
		internal.DoRespond(w, r, respond.Error(err))
	}
}
