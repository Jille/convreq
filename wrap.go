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

type extractor func(w http.ResponseWriter, r *http.Request) (reflect.Value, HttpResponse)

var extractorMap = map[reflect.Type]extractor{
	reflect.TypeOf((*context.Context)(nil)).Elem():     getContext,
	reflect.TypeOf(&http.Request{}):                    getRequest,
	reflect.TypeOf((*http.ResponseWriter)(nil)).Elem(): getResponseWriter,
}

var handlerMap = map[reflect.Type]reflect.Value{
	reflect.TypeOf((*error)(nil)).Elem():        reflect.ValueOf(handleError),
	reflect.TypeOf((*HttpResponse)(nil)).Elem(): reflect.ValueOf(internal.DoRespond),
}

type wrapOptions struct{
	extractors map[reflect.Type]extractor
	handlers map[reflect.Type]reflect.Value
}

type WrapOption func(wo *wrapOptions)

func WithParameterType(t reflect.Type, e extractor) WrapOption {
	return func(wo *wrapOptions) {
		wo.extractors[t] = e
	}
}

func WithReturnType(t reflect.Type, f reflect.Value) WrapOption {
	return func(wo *wrapOptions) {
		wo.handlers[t] = f
	}
}

// Wrap takes a request handler function and returns a http.HandlerFunc for use with net/http.
// The requested handler is expected to take arguments like context.Context, *http.Request and return a convreq.HttpResponse or an error.
func Wrap(f interface{}, opts ...WrapOption) http.HandlerFunc {
	t := reflect.TypeOf(f)
	v := reflect.ValueOf(f)

	wo := wrapOptions{
		extractors: map[reflect.Type]extractor{},
		handlers: map[reflect.Type]reflect.Value{},
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
			ins[i] = createGet(t.In(i))
		} else if t.In(i).Kind() == reflect.Ptr && strings.HasSuffix(t.In(i).Elem().Name(), "Post") {
			ins[i] = createPost(t.In(i))
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

func createGet(t reflect.Type) func(w http.ResponseWriter, r *http.Request) (reflect.Value, HttpResponse) {
	return func(w http.ResponseWriter, r *http.Request) (reflect.Value, HttpResponse) {
		// TODO(quis): Consider putting v in a sync.Pool.
		v := reflect.New(t)
		if err := internal.DecodeGet(r, v.Interface()); err != nil {
			return reflect.Value{}, respond.BadRequest(err)
		}
		return v.Elem(), nil
	}
}

func createPost(pt reflect.Type) func(w http.ResponseWriter, r *http.Request) (reflect.Value, HttpResponse) {
	nilptr := reflect.New(pt).Elem()
	t := pt.Elem()
	return func(w http.ResponseWriter, r *http.Request) (reflect.Value, HttpResponse) {
		if r.Method != "POST" {
			return nilptr, nil
		}
		// TODO(quis): Consider putting v in a sync.Pool.
		v := reflect.New(t)
		if err := internal.DecodePost(r, v.Interface()); err != nil {
			return reflect.Value{}, respond.BadRequest(err)
		}
		return v, nil
	}
}

// === Below are some functions that can handle the return value of a request handler.

func handleVoid(w http.ResponseWriter, r *http.Request) {
}

func handleError(w http.ResponseWriter, r *http.Request, err error) {
	if err != nil {
		internal.DoRespond(w, r, respond.InternalServerError(err))
	}
}
