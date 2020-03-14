package convreq_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Jille/convreq"
	"github.com/Jille/convreq/respond"
)

type ArticlesCategoryHandlerGet struct {
	Category string
	Id       int64
}

type ArticlesCategoryHandlerPost struct {
	NewName string
}

func ArticlesCategoryHandler(ctx context.Context, r *http.Request, get ArticlesCategoryHandlerGet, post *ArticlesCategoryHandlerPost) convreq.HttpResponse {
	if get.Category == "unimplemented" {
		return respond.InternalServerError(fmt.Errorf("not yet implemented"))
	}
	if post != nil {
		return respond.String(fmt.Sprintf("I like post. NewName=%s", post.NewName))
	}
	return respond.String(fmt.Sprintf("Hello world. Id=%d", get.Id))
}

func TestStuff(t *testing.T) {
	tests := []struct {
		req      *http.Request
		handler  interface{}
		wantCode int
		wantBody string
	}{
		{
			req:      httptest.NewRequest("GET", "/?category=unimplemented&id=1", nil),
			handler:  ArticlesCategoryHandler,
			wantCode: 500,
			wantBody: "not yet implemented\n",
		},
		{
			req:      httptest.NewRequest("GET", "/?category=test&id=7", nil),
			handler:  ArticlesCategoryHandler,
			wantCode: 200,
			wantBody: "Hello world. Id=7",
		},
		{
			req:      httptest.NewRequest("POST", "/?category=test&id=7", strings.NewReader("newname=dude")),
			handler:  ArticlesCategoryHandler,
			wantCode: 200,
			wantBody: "I like post. NewName=dude",
		},
	}

	for _, tc := range tests {
		t.Run(tc.req.URL.String(), func(t *testing.T) {
			if tc.req.Method == "POST" {
				tc.req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
			}
			respRecorder := httptest.NewRecorder()
			var handler http.Handler = convreq.Wrap(tc.handler)
			handler.ServeHTTP(respRecorder, tc.req)
			if respRecorder.Code != tc.wantCode {
				t.Errorf("Request returned code %d (want %d)", respRecorder.Code, tc.wantCode)
			}
			body := respRecorder.Body.String()
			if body != tc.wantBody {
				t.Errorf("Requested return unexpected body: %q; want %q", body, tc.wantBody)
			}
		})
	}
}
