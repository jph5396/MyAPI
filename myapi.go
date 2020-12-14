package myapi

import (
	"net/http"
	"net/http/httptest"

	"github.com/dimfeld/httptreemux/v5"
)

//MyAPI ...
type MyAPI struct {
	Name             string
	subrouters       []SubRouter
	Routes           []Route
	globalmiddleware []Middleware
	managedRouter    httptreemux.ContextMux
	Port             string
}

//NewMyAPI creates a new MyAPI instance.
func NewMyAPI(name string, port string) MyAPI {
	mux := httptreemux.NewContextMux()
	return MyAPI{
		Name:          name,
		Port:          port,
		managedRouter: *mux,
	}
}

//UseMiddleware applies the middleware provided to the router.
// this middleware is used for all routes.
func (m *MyAPI) UseMiddleware(mw Middleware) {
	m.managedRouter.UseHandler(mw.handler)
	m.globalmiddleware = append(m.globalmiddleware, mw)
}

//UseSubrouter applies a SubRouter to the api.
func (m *MyAPI) UseSubrouter(sr SubRouter) {

	sub := m.managedRouter.NewContextGroup(sr.prefix)
	//apply subrouters middleware.
	for _, mw := range sr.middleware {
		sub.UseHandler(mw.handler)
	}

	// add all routes to the subrouter.
	for _, route := range sr.routes {
		sub.Handle(route.method, route.path, route.handler)
	}
	m.subrouters = append(m.subrouters, sr)
}

//UseRoute applies an individual route to the api that is
// not part of any subrouter.
func (m *MyAPI) UseRoute(r Route) {
	m.managedRouter.Handle(r.method, r.path, r.handler)
}

//StartServer starts the server using the port and managed routers.
func (m *MyAPI) StartServer() error {
	return http.ListenAndServe(m.Port, m.managedRouter)
}

//StartTestServer starts the server using httptest.NewServer() instead of
// http.ListenAndServe for testing purposes.
func (m *MyAPI) StartTestServer() *httptest.Server {
	return httptest.NewServer(m.managedRouter)
}
