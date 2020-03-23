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

//go:generate sh -c "go run ../cmd/generate/generate.go -- example.go > example_generated.go"
package main

import (
	"context"
	"log"
	"net/http"

	"github.com/Jille/convreq"
	"github.com/Jille/convreq/respond"
	"github.com/gorilla/mux"
)

type server struct{}

func main() {
	s := &server{}
	_ = s
	r := mux.NewRouter()
	// r.HandleFunc("/products/{key}", crqProductHandler)
	r.HandleFunc("/articles/{category}/", s.crqArticlesCategoryHandler)
	r.HandleFunc("/articles/{category}/{id:[0-9]+}", convreq.Wrap(ArticleHandler))

	srv := &http.Server{
		Handler: r,
		Addr:    ":8007",
	}

	log.Fatal(srv.ListenAndServe())
}

type ArticlesCategoryHandlerGet struct {
	Category string
	Id       int64
}

type ArticlesCategoryHandlerPost struct {
	NewName string
}

func (s *server) ArticlesCategoryHandler(ctx context.Context, r *http.Request, get ArticlesCategoryHandlerGet, post *ArticlesCategoryHandlerPost) convreq.HttpResponse {
	return respond.InternalServerError("not yet implemented")
}

type ArticleHandlerGet struct {
	Category string
	Id       int64
	Piet     string
}

type ArticleHandlerPost struct {
	NewName string
}

func ArticleHandler(ctx context.Context, r *http.Request, get ArticleHandlerGet, post *ArticleHandlerPost) convreq.HttpResponse {
	return respond.InternalServerError("also not yet implemented")
}
