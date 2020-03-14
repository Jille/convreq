package convreq_test

import (
	"context"
	"errors"
	"fmt"
	"text/template"
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
	NewName string `schema:"newname,required"`
}

func ArticlesCategoryHandler(ctx context.Context, r *http.Request, get ArticlesCategoryHandlerGet, post *ArticlesCategoryHandlerPost) convreq.HttpResponse {
	if get.Category == "unimplemented" {
		return respond.InternalServerError("not yet implemented")
	}
	if post != nil {
		return respond.String(fmt.Sprintf("I like post. NewName=%s", post.NewName))
	}
	return respond.String(fmt.Sprintf("Hello world. Id=%d", get.Id))
}

func TestStuff(t *testing.T) {
	tests := []struct {
		req         *http.Request
		handler     interface{}
		wantCode    int
		wantHeaders map[string]string
		wantBody    string
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
		{
			req:      httptest.NewRequest("POST", "/?category=test&id=not-a-number", strings.NewReader("newname=dude")),
			handler:  ArticlesCategoryHandler,
			wantCode: 400,
			wantBody: "failed to parse url/query: schema: error converting value for \"id\"\n",
		},
		{
			req:      httptest.NewRequest("POST", "/?category=test&id=7", strings.NewReader("newname=")),
			handler:  ArticlesCategoryHandler,
			wantCode: 400,
			wantBody: "failed to parse form input: newname is empty\n",
		},
		{
			// Test return value error.
			req: httptest.NewRequest("GET", "/", nil),
			handler: func() error {
				return errors.New("test")
			},
			wantCode: 500,
			wantBody: "test\n",
		},
		{
			req: httptest.NewRequest("GET", "/", nil),
			// Test no return value.
			handler: func(w http.ResponseWriter, r *http.Request) {
				http.NotFoundHandler().ServeHTTP(w, r)
			},
			wantCode: 404,
			wantBody: "404 page not found\n",
		},
		{
			req: httptest.NewRequest("GET", "/", nil),
			handler: func() convreq.HttpResponse {
				return respond.BadRequest("test")
			},
			wantCode: 400,
			wantBody: "test\n",
		},
		{
			req: httptest.NewRequest("GET", "/", nil),
			handler: func() convreq.HttpResponse {
				return respond.PermissionDenied("test")
			},
			wantCode: 403,
			wantBody: "test\n",
		},
		{
			req: httptest.NewRequest("GET", "/", nil),
			handler: func() convreq.HttpResponse {
				return respond.NotFound("test")
			},
			wantCode: 404,
			wantBody: "test\n",
		},
		{
			req: httptest.NewRequest("GET", "/", nil),
			handler: func() convreq.HttpResponse {
				return respond.Handler(http.NotFoundHandler())
			},
			wantCode: 404,
			wantBody: "404 page not found\n",
		},
		{
			req: httptest.NewRequest("GET", "/", nil),
			handler: func() convreq.HttpResponse {
				return respond.OverrideResponseCode(respond.Handler(http.NotFoundHandler()), 406)
			},
			wantCode: 406,
			wantBody: "404 page not found\n",
		},
		{
			req: httptest.NewRequest("GET", "/", nil),
			handler: func() convreq.HttpResponse {
				return respond.Redirect(302, "/login/")
			},
			wantCode:    302,
			wantHeaders: map[string]string{"Location": "/login/"},
			wantBody:    "",
		},
		{
			req: httptest.NewRequest("GET", "/", nil),
			handler: func() convreq.HttpResponse {
				return respond.WithHeader(respond.Handler(http.NotFoundHandler()), "A", "B")
			},
			wantCode:    404,
			wantHeaders: map[string]string{"A": "B"},
			wantBody:    "404 page not found\n",
		},
		{
			req: httptest.NewRequest("GET", "/", nil),
			handler: func() convreq.HttpResponse {
				return respond.WithHeaders(respond.Handler(http.NotFoundHandler()), http.Header{"A": []string{"B"}, "C": []string{"D"}})
			},
			wantCode:    404,
			wantHeaders: map[string]string{"A": "B", "C": "D"},
			wantBody:    "404 page not found\n",
		},
		{
			req: httptest.NewRequest("GET", "/", nil),
			handler: func() convreq.HttpResponse {
				t, err := template.New("test").Parse("{{.}} is a number")
				if err != nil {
					return respond.Error(err)
				}
				return respond.RenderTemplate(t, 7)
			},
			wantCode: 200,
			wantBody: "7 is a number",
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
			for h, v := range tc.wantHeaders {
				if respRecorder.Header().Get(h) != v {
					t.Errorf("Request returned %q for header %q (want %q)", respRecorder.Header().Get(h), h, v)
				}
			}
			body := respRecorder.Body.String()
			if body != tc.wantBody {
				t.Errorf("Requested return unexpected body: %q; want %q", body, tc.wantBody)
			}
		})
	}
}
