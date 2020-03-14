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
