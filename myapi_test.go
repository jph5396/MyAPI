package myapi

import (
	"errors"
	"fmt"
	"net/http"
	"testing"
)

func TestMyAPI(t *testing.T) {
	testserver := buildtestserver()
	srv := testserver.StartTestServer()
	defer srv.Close()

	testplainroute := func() error {
		req, err := http.NewRequest(http.MethodGet, srv.URL+"/plainroute", nil)
		if err != nil {
			return err
		}
		res, err := http.DefaultClient.Do(req)
		if res.StatusCode != 200 {
			return fmt.Errorf("wrong status code: got %v want %v", res.StatusCode, 200)
		}
		if head := res.Header.Get("global"); head == "" {
			return errors.New("middleware fail: global-header missing")
		}
		return nil
	}

	testsubrouter := func() error {
		req, err := http.NewRequest(http.MethodGet, srv.URL+"/subroute/plainroute", nil)
		if err != nil {
			return err
		}
		res, err := http.DefaultClient.Do(req)
		if res.StatusCode != 200 {
			return fmt.Errorf("wrong status code: got %v want %v", res.StatusCode, 200)
		}
		if head := res.Header.Get("sub-header"); head == "" {
			return errors.New("middleware fail: sub-header missing")
		}
		return nil
	}

	var tests = []struct {
		testname string
		testfunc func() error
	}{
		{"hit plainroute", testplainroute},
		{"hit subroute", testsubrouter},
	}

	for _, test := range tests {
		t.Run(test.testname, func(t *testing.T) {
			err := test.testfunc()
			if err != nil {
				t.Error(err)
			}
		})
	}
}
func returnnil() error {
	return nil
}
func buildtestserver() MyAPI {

	testserver := NewMyAPI("test", ":8080")
	route := NewRoute("/plainroute", http.MethodGet, func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("plainroute"))
		return
	})
	testserver.UseRoute(route)

	globalmw := NewMiddleware("testglobalmw", func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("global", "found")
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
