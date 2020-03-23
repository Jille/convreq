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
	"fmt"
	"net/url"
	"strings"

	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
)

var encoder = schema.NewEncoder()

func URL(r *mux.Route, strct interface{}) (*url.URL, error) {
	values := url.Values{}
	if err := encoder.Encode(strct, values); err != nil {
		return nil, err
	}
	queries, err := r.GetQueriesTemplates()
	if err != nil && !strings.Contains(err.Error(), "route doesn't have queries") {
		return nil, err
	}
	params := make([]string, 0, len(queries)*2)
	for _, q := range queries {
		switch len(values[q][0]) {
		case 0:
			return nil, fmt.Errorf("field %q required for path missing", q)
		case 1:
			params = append(params, q, values[q][0])
			delete(values, q)
		default:
			return nil, fmt.Errorf("multiple values for field %q encountered", q)
		}
	}
	u, err := r.URL(params...)
	if err != nil {
		return nil, err
	}
	if u.RawQuery == "" {
		u.RawQuery = values.Encode()
	} else {
		u.RawQuery += "&" + values.Encode()
	}
	return u, nil
}
