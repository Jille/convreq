package respond

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strconv"

	internal "github.com/Jille/convreq/internal"
)

// TemplateProducer is a function that returns a template (like) object.
// This is a function to encourage you to implement hot reloading for development servers.
// For production servers we recommend parsing the template at startup and then have a closure that returns that.
type TemplateProducer func() interface {
	Execute(wr io.Writer, data interface{}) error
}

// RenderTemplate creates a response that will render the given template.
// See TemplateProducer.
func RenderTemplate(tplp TemplateProducer, data interface{}) internal.HttpResponse {
	var buf bytes.Buffer
	if err := tplp().Execute(&buf, data); err != nil {
		return InternalServerError(fmt.Errorf("failed to render template: %v", err))
	}
	return Reader(&buf)
}

type repondReader struct {
	r io.Reader
}

// Respond implements convreq.HttpResponse.
func (rr repondReader) Respond(w http.ResponseWriter, r *http.Request) error {
	if w.Header().Get("Content-Length") == "" {
		if lenner, ok := rr.r.(interface{ Len() int }); ok {
			w.Header().Set("Content-Length", strconv.Itoa(lenner.Len()))
		}
	}
	if _, err := io.Copy(w, rr.r); err != nil {
		return fmt.Errorf("failed to write response to client: %v", err)
	}
	if closer, ok := rr.r.(io.Closer); ok {
		closer.Close()
	}
	return nil
}

// Reader creates a response that copies everything from the reader to the client.
func Reader(r io.Reader) internal.HttpResponse {
	return repondReader{r}
}
