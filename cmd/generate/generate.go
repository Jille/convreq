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

package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"regexp"
)

func main() {
	flag.Parse()
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

var (
	// vi: (
	defineRe = regexp.MustCompile(`(?m)^func (?:\([a-z]+ ([^)]+)\) )?([A-Za-z0-9]+)\(.*http\.Request`) // vi: )
)

type handler struct {
	name          string
	ontype        string
}

func run() error {
	var handlers []handler
	for _, fn := range flag.Args() {
		code, err := ioutil.ReadFile(fn)
		if err != nil {
			return err
		}
		for _, m := range defineRe.FindAllStringSubmatch(string(code), -1) {
			handlers = append(handlers, handler{
				name: m[2],
				ontype: m[1],
			})
		}
	}
	generate(os.Stdout, handlers)
	return nil
}

func generate(w io.Writer, handlers []handler) {
	fmt.Fprintf(w, "package main\n")
	fmt.Fprintf(w, "\n")
	fmt.Fprintf(w, "import (\n")
	fmt.Fprintf(w, "\t%q\n", "fmt")
	fmt.Fprintf(w, "\t%q\n", "net/http")
	fmt.Fprintf(w, "\t%q\n", "net/url")
	fmt.Fprintf(w, "\n")
	fmt.Fprintf(w, "\t%q\n", "github.com/Jille/convreq")
	fmt.Fprintf(w, "\t%q\n", "github.com/Jille/convreq/genapi")
	fmt.Fprintf(w, "\t%q\n", "github.com/Jille/convreq/internal")
	fmt.Fprintf(w, "\t%q\n", "github.com/Jille/convreq/respond")
	fmt.Fprintf(w, "\t%q\n", "github.com/gorilla/mux")
	fmt.Fprintf(w, ")\n")

	for _, h := range handlers {
		base := h.name
		prefix := ""
		onstruct := ""
		if h.ontype != "" {
			onstruct = "(t " + h.ontype + ") "
			prefix = "t."
		}
		fmt.Fprintf(w, "\n")
		fmt.Fprintf(w, "func %scrq%s(w http.ResponseWriter, r *http.Request) {\n", onstruct, base)
		fmt.Fprintf(w, "\thr := %s_crqInternal%s(w, r)\n", prefix, base)
		fmt.Fprintf(w, "\tinternal.DoRespond(w, r, hr)\n")
		fmt.Fprintf(w, "}\n")
		fmt.Fprintf(w, "\n")
		fmt.Fprintf(w, "func %s_crqInternal%s(w http.ResponseWriter, r *http.Request) (resp convreq.HttpResponse) {\n", onstruct, base)
		fmt.Fprintf(w, "\tdefer genapi.PanicHandler(&resp)()\n")
		fmt.Fprintf(w, "\tvar get %sGet\n", base)
		fmt.Fprintf(w, "\tif err := internal.DecodeGet(r, &get); err != nil {\n")
		fmt.Fprintf(w, "\t\treturn respond.BadRequest(err.Error())\n")
		fmt.Fprintf(w, "\t}\n")
		fmt.Fprintf(w, "\tvar postptr *%sPost\n", base)
		fmt.Fprintf(w, "\tif r.Method == %q {\n", "POST")
		fmt.Fprintf(w, "\t\tvar post %sPost\n", base)
		fmt.Fprintf(w, "\t\tif err := internal.DecodePost(r, &post); err != nil {\n")
		fmt.Fprintf(w, "\t\t\treturn respond.BadRequest(err.Error())\n")
		fmt.Fprintf(w, "\t\t}\n")
		fmt.Fprintf(w, "\t\tpostptr = &post\n")
		fmt.Fprintf(w, "\t}\n")
		fmt.Fprintf(w, "\treturn %s%s(r.Context(), r, get, postptr)\n", prefix, base)
		fmt.Fprintf(w, "}\n")
		fmt.Fprintf(w, "\n")
		fmt.Fprintf(w, "func (g %sGet) URL(r *mux.Route) *url.URL {\n", base)
		fmt.Fprintf(w, "\tu, err := convreq.URL(r, g)\n")
		fmt.Fprintf(w, "\tif err != nil {\n")
		fmt.Fprintf(w, "\t\tpanic(fmt.Errorf(\"failed to construct URL: %%v\", err))\n")
		fmt.Fprintf(w, "\t}\n")
		fmt.Fprintf(w, "\treturn u\n")
		fmt.Fprintf(w, "}\n")
	}
}
