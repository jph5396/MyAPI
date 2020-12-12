package myapi

import (
	"fmt"
	"net/http"
	"testing"
)

func TestRouter(t *testing.T) {
	testrouter := NewSubRouter("/v1")

	route := NewRoute("/test", http.MethodPost,
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Println("hello")
		}))
	testrouter.AddRoute(route)

	if len(testrouter.routes) != 1 {
		t.Errorf("len of test router is %d; should be 1", len(testrouter.routes))
	}
}
