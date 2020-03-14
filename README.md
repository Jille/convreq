# convreq

[![GoDoc](https://godoc.org/github.com/Jille/convreq?status.svg)](https://godoc.org/github.com/Jille/convreq)
[![Build Status](https://travis-ci.org/Jille/convreq.png)](https://travis-ci.org/Jille/convreq)

Experimental project to make writing webservers more convenient.

The core principle of the library is that while the `func(w http.ResponseWriter, r *http.Request)` signature is very powerful, often a more convenient interface would be preferable.

# Signature

I propose a signature that:

* passes a context directly (as the first parameter)
* simplifies error handling by allowing a return value
* passes in type casted input (GET+POST) to the request handler
* defers responding until the function has returned so later errors can still be handled and you're not caught with errors halfway through your response

So I arrived at `func MyHandler(ctx context.Context, r *http.Request, get MyHandlerGet, post *MyHandlerPost) convreq.HttpResponse`.

`get` and `post` are structs defined by the application that hold input fields. github.com/gorilla/schema is used to decode the GET/POST values into these structs.

`get` also contains any URL parameters for github.com/gorilla/mux.

# Dispatchers

I have implemented two different dispatchers: one with code generation, and one with reflect. Advantages of each:

* The reflect dispatcher doesn't enforce an argument order. It looks at the types of each parameter and constructs the value required. This will be convenient for e.g. github.com/gorilla/sessions as handlers that require sessions can trivially add a parameter to receive a session.
* The codegen dispatcher does enforce an argument order, thus promoting consistency.
* The codegen dispatcher is most likely much faster as it does far less magic at runtime.
* The codegen dispatcher will catch more problems at compile time. The reflect dispatcher tries its best to catch errors at initialization time, but some might only be noticed at request time.
* The reflect dispatcher is easier to use, doesn't depend on go-generate or requiring the programmer to know when to regenerate.

Maybe I'll make the codegen dispatcher smarter to also allow for more freedom in method signature.

I'm considering to keep both dispatchers, encourage the reflect dispatcher for development and the codegen dispatcher for prod deployments.

# Responding

Rather than allowing the code to respond pieces whenever it wants, I enforce that all responding happens after returning from the request handler. An interface HttpResponse can be returned and one can implement it to do anything you want.

I'll implement a couple of basic responders:

* Rendering errors easily. Setting the correct HTTP response code and optionally allowing to configure an error page renderer.
* Template rendering. This should be as easy as `return convreq.RenderTemplate(myTemplate, myData)`. This'll make it possible to automatically reload templates for development servers too.
* Redirection is as easy as `return convreq.Redirect(302, "/home")`.
* Setting response headers.
