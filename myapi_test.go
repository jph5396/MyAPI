package myapi

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"testing"
)

func TestMyAPI(t *testing.T) {
	testserver := buildtestserver()
	srv, err := testserver.StartTestServer()
	if err != nil {
		t.Fatalf("test server could not start. reason: %v", err.Error())
	}
	defer srv.Close()

	testplainroute := func() error {
		req, err := http.NewRequest(http.MethodGet, srv.URL+"/plainroute", nil)
		if err != nil {
			return err
		}
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			return err
		}
		defer res.Body.Close()
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
		if err != nil {
			return err
		}
		defer res.Body.Close()
		if res.StatusCode != 200 {
			return fmt.Errorf("wrong status code: got %v want %v", res.StatusCode, 200)
		}
		if head := res.Header.Get("sub-header"); head == "" {
			return errors.New("middleware fail: sub-header missing")
		}
		return nil
	}

	testproproute := func() error {
		reqbody, err := json.Marshal(map[string]string{
			"Test": "Test",
		})
		if err != nil {
			return errors.New("could not marshal json body for testproproute")
		}

		req, err := http.NewRequest(http.MethodPost, srv.URL+"/subroute/proptest", bytes.NewBuffer(reqbody))
		if err != nil {
			return err
		}

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			return err
		}
		defer res.Body.Close()
		if res.StatusCode != 200 {
			return fmt.Errorf("wrong status code: got %v want %v", res.StatusCode, 200)
		}

		return nil
	}

	var tests = []struct {
		testname string
		testfunc func() error
	}{
		{"hit plainroute", testplainroute},
		{"hit subroute", testsubrouter},
		{"hit route with prop", testproproute},
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
	route := NewRoute("/plainroute", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("plainroute"))
	}, http.MethodGet)
	err := testserver.UseRoute(route)
	if err != nil {
		panic(err)
	}

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

	//create Properties and a route to test them on.
	prop := NewProperty("Test", String)
	propRouteTest := NewRoute("/proptest", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("proptest"))
	}, http.MethodPost)

	err = propRouteTest.AddProperty(prop)
	if err != nil {
		panic(err)
	}
	sub.UseMiddleware(submw)
	sub.AddRoute(route)
	sub.AddRoute(propRouteTest)
	err = testserver.UseSubrouter(sub)
	if err != nil {
		panic(err)
	}

	return testserver
}
