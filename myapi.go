package myapi

import (
	"net/http"
	"net/http/httptest"

	"github.com/gorilla/mux"
)

//MyAPI ...
type MyAPI struct {
	Name             string
	subrouters       []SubRouter
	Routes           []Route
	globalmiddleware []Middleware
	managedRouter    mux.Router
	Port             string
}

//NewMyAPI creates a new MyAPI instance.
func NewMyAPI(name string, port string) MyAPI {
	mux := mux.NewRouter()
	return MyAPI{
		Name:          name,
		Port:          port,
		managedRouter: *mux,
	}
}

//UseMiddleware applies the middleware provided to the router.
// this middleware is used for all routes.
func (m *MyAPI) UseMiddleware(mw Middleware) {
	m.managedRouter.Use(mw.handler)
	m.globalmiddleware = append(m.globalmiddleware, mw)
}

//UseSubrouter applies a SubRouter to the api.
func (m *MyAPI) UseSubrouter(sr SubRouter) {

	sub := m.managedRouter.PathPrefix(sr.prefix).Subrouter()

	// add all routes to the subrouter.
	for _, route := range sr.routes {
		sub.Handle(route.path, route.handler).Methods(route.method)
	}
	//apply subrouters middleware.
	for _, mw := range sr.middleware {
		sub.Use(mw.handler)
	}
	m.subrouters = append(m.subrouters, sr)
}

//UseRoute applies an individual route to the api that is
// not part of any subrouter.
func (m *MyAPI) UseRoute(r Route) {
	m.managedRouter.Handle(r.path, r.handler).Methods(r.method)
}

//StartServer starts the server using the port and managed routers.
func (m *MyAPI) StartServer() error {
	return http.ListenAndServe(m.Port, &m.managedRouter)
}

//StartTestServer starts the server using httptest.NewServer() instead of
// http.ListenAndServe for testing purposes.
func (m *MyAPI) StartTestServer() *httptest.Server {
	return httptest.NewServer(&m.managedRouter)
}
