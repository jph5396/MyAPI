package myapi

import (
	"net/http"
	"testing"
)

func TestMyAPI(t *testing.T) {
	testserver := buildtestserver()
	srv := testserver.StartTestServer()
	defer srv.Close()
}

func buildtestserver() MyAPI {

	testserver := NewMyAPI("test", ":8080")
	route := NewRoute("/plainroute", http.MethodGet, func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("plainroute"))
	})
	testserver.UseRoute(route)

	globalmw := NewMiddleware("testglobalmw", func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("global-header", "found")
			next.ServeHTTP(w, r)
		})
	})
	testserver.UseMiddleware(globalmw)

	sub := NewSubRouter("/subroute")
	submw := NewMiddleware("testsubmw", func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("sub-header", "found")
			next.ServeHTTP(w, r)
		})
	})
	sub.UseMiddleware(submw)
	sub.AddRoute(route)
	testserver.UseSubrouter(sub)

	return testserver
}
