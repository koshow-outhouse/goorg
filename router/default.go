package router

import (
	"net/http"
	"path"

	"github.com/julienschmidt/httprouter"
	"github.com/organization-service/goorg/logger"
)

type (
	defaultRouter struct {
		*httprouter.Router
		group string
	}
)

func newDefault() IRouter {
	router := &defaultRouter{
		Router: httprouter.New(),
	}
	router.PanicHandler = panicHandler
	return router
}

// DELETE replaces httprouter.Router.DELETE.
func (r *defaultRouter) DELETE(path string, h interface{}) {
	r.Router.DELETE(joinURL(r, path), logger.Log(h))
}

// GET replaces httprouter.Router.GET.
func (r *defaultRouter) GET(path string, h interface{}) {
	r.Router.GET(joinURL(r, path), logger.Log(h))
}

// HEAD replaces httprouter.Router.HEAD.
func (r *defaultRouter) HEAD(path string, h interface{}) {
	r.Router.HEAD(joinURL(r, path), logger.Log(h))
}

// OPTIONS replaces httprouter.Router.OPTIONS.
func (r *defaultRouter) OPTIONS(path string, h interface{}) {
	r.Router.OPTIONS(joinURL(r, path), logger.Log(h))
}

// PATCH replaces httprouter.Router.PATCH.
func (r *defaultRouter) PATCH(path string, h interface{}) {
	r.Router.PATCH(joinURL(r, path), logger.Log(h))
}

// POST replaces httprouter.Router.POST.
func (r *defaultRouter) POST(path string, h interface{}) {
	r.Router.POST(joinURL(r, path), logger.Log(h))
}

// PUT replaces httprouter.Router.PUT.
func (r *defaultRouter) PUT(path string, h interface{}) {
	r.Router.PUT(joinURL(r, path), logger.Log(h))
}

// Handle replaces httprouter.Router.Handle.
func (r *defaultRouter) Handle(method, path string, h interface{}) {
	r.Router.Handle(method, joinURL(r, path), logger.Log(h))
}

// Handler replaces httprouter.Router.Handler.
func (r *defaultRouter) Handler(method, path string, handler http.Handler) {
	h := logger.Log(handler)
	r.Router.Handler(method, joinURL(r, path), http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		p := httprouter.ParamsFromContext(r.Context())
		h(rw, r, p)
	}))
}

// HandlerFunc replaces httprouter.Router.HandlerFunc.
func (r *defaultRouter) HandlerFunc(method, path string, handler http.HandlerFunc) {
	h := logger.Log(handler)
	r.Router.HandlerFunc(method, joinURL(r, path), func(rw http.ResponseWriter, r *http.Request) {
		p := httprouter.ParamsFromContext(r.Context())
		h(rw, r, p)
	})
}

// ServeHTTP replaces httprouter.Router.ServeHTTP.
func (r *defaultRouter) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.Router.ServeHTTP(w, req)
}

func (r *defaultRouter) ServeFiles(path string, fileSystem http.FileSystem) {
	r.Router.ServeFiles(joinURL(r, path), fileSystem)
}

func (r *defaultRouter) GlobalOPTIONS(h http.HandlerFunc) {
	r.Router.GlobalOPTIONS = h
}

func (r *defaultRouter) Group(url string) IRouter {
	return &defaultRouter{
		Router: r.Router,
		group:  path.Join(r.group, url),
	}
}

func (r *defaultRouter) getGroup() string {
	return r.group
}
