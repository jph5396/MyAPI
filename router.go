package myapi

//SubRouter a grouping of routes. all routes will have a prefix
// applied to them at 8
type SubRouter struct {
	prefix     string
	routes     []Route
	middleware []Middleware
}

//NewSubRouter returns an empty subrouter with the provided prefix
func NewSubRouter(p string) SubRouter {
	return SubRouter{prefix: p}
}

//AddRoute adds the route provided to the router.
func (sr *SubRouter) AddRoute(r Route) {
	sr.routes = append(sr.routes, r)
}

//UseMiddleware applies the middleware to the provided
func (sr *SubRouter) UseMiddleware(mw Middleware) {
	sr.middleware = append(sr.middleware, mw)
}
