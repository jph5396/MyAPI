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
	routeProps       map[string]PropertyGroup
}

//NewMyAPI creates a new MyAPI instance.
func NewMyAPI(name string, port string) MyAPI {
	myapi := MyAPI{
		Name: name,
		Port: port,
	}
	mux := mux.NewRouter()
	myapi.managedRouter = *mux
	return myapi
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
func (m *MyAPI) StartTestServer() (*httptest.Server, error) {
	return httptest.NewServer(&m.managedRouter), nil
}

//myapiMiddleware middleware function used to implement all type checking and rule
// validation on endpoints.
func (m *MyAPI) myapiMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path, err := mux.CurrentRoute(r).GetPathTemplate()

		//err should only be nonnil when there isnt a matching route.
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(http.StatusText(http.StatusNotFound)))
			return
		}
		props, ok := m.routeProps[path]
		if !ok {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(http.StatusText(http.StatusInternalServerError)))
		}

		//TODO: finish
		next.ServeHTTP(w, r)
	})
}

//build adds anything to myapi that is needed for it to run, but needs to be set
//after all routes and middleware have been added.
func (m *MyAPI) build() error {
	//TODO: add code to fill routeProps
	m.managedRouter.Use(m.myapiMiddleware)
	return nil
}
