/*
   Restfool-go

   Copyright (C) 2018 Carsten Seeger

   This program is free software: you can redistribute it and/or modify
   it under the terms of the GNU General Public License as published by
   the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.

   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU General Public License for more details.

   You should have received a copy of the GNU General Public License
   along with this program.  If not, see <http://www.gnu.org/licenses/>.

   @author Carsten Seeger
   @copyright Copyright (C) 2018 Carsten Seeger
   @license http://www.gnu.org/licenses/gpl-3.0 GNU General Public License 3
   @link https://github.com/cseeger-epages/rest-api-go-skeleton
*/

package restfool

import (
	"net/http"

	"github.com/gorilla/mux"
	"gopkg.in/throttled/throttled.v2"
	"gopkg.in/throttled/throttled.v2/store/memstore"
)

var prefixList []pathPrefix

// NewRouter is the router constructor
func (a RestAPI) NewRouter() *mux.Router {

	router := mux.NewRouter().StrictSlash(true)

	// latest
	a.AddRoutes(router)

	//v1
	a.AddV1Routes(router.PathPrefix("/v1").Subrouter())

	//v2 only dummy yet
	a.AddV2Routes(router.PathPrefix("/v2").Subrouter())

	// add additional prefix handler
	a.addPrefixList(router)

	return router
}

// AddRoutes add default handler, routing and ratelimit
func (a RestAPI) AddRoutes(router *mux.Router) {
	store, err := memstore.New(65536)
	Error("ROUTES: could not create memstore", err)

	// rate limiter
	quota := throttled.RateQuota{
		MaxRate:  throttled.PerMin(a.Conf.RateLimit.Limit),
		MaxBurst: a.Conf.RateLimit.Burst,
	}
	rateLimiter, err := throttled.NewGCRARateLimiter(store, quota)
	Error("ROUTES: error in ratelimiting", err)

	httpRateLimiter := throttled.HTTPRateLimiter{
		RateLimiter: rateLimiter,
		VaryBy:      &throttled.VaryBy{Path: true},
	}

	for _, route := range routes {
		var handler http.Handler

		handler = route.HandlerFunc
		handler = Logger(handler, route.Name)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(use(handler, a.addDefaultHeader, a.basicAuthHandler, httpRateLimiter.RateLimit))
	}
}

// Middleware chainer
func use(h http.Handler, middleware ...func(http.Handler) http.Handler) http.Handler {
	for _, m := range middleware {
		h = m(h)
	}

	return h
}

// AddV1Routes version 1 routes
func (a RestAPI) AddV1Routes(router *mux.Router) {
	a.AddRoutes(router)
}

// AddV2Routes dummy for version 2 routes
func (a RestAPI) AddV2Routes(router *mux.Router) {
	a.AddRoutes(router)
}

// AddPathPrefix to add additional Path Prefix handler
func (a RestAPI) AddPathPrefix(prefix string, handler http.Handler) {
	prefixList = append(prefixList, pathPrefix{prefix, handler})
}

func (a RestAPI) addPrefixList(router *mux.Router) {
	for _, v := range prefixList {
		router.PathPrefix(v.Prefix).Handler(v.Handler)
	}
}
