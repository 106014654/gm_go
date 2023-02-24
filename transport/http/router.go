package http

import (
	"net/http"
	"path"
	"sync"
)

type HandlerFunc func(Context) error

type Router struct {
	prefix  string
	pool    sync.Pool
	srv     *Server
	filters []FilterFunc
}

func newRouter(prefix string, srv *Server, filters ...FilterFunc) *Router {
	r := &Router{
		prefix:  prefix,
		srv:     srv,
		filters: filters,
	}
	r.pool.New = func() interface{} {
		return &wrapper{router: r}
	}
	return r
}

func (r *Router) Group(prefix string, filters ...FilterFunc) *Router {
	var newFilters []FilterFunc
	newFilters = append(newFilters, r.filters...)
	newFilters = append(newFilters, filters...)
	return newRouter(path.Join(r.prefix, prefix), r.srv, newFilters...)
}

func (r *Router) Handle(method, relativePath string, h HandlerFunc, filters ...FilterFunc) {
	next := http.Handler(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		ctx := r.pool.Get().(Context)
		ctx.Reset(res, req)
		if err := h(ctx); err != nil {
			r.srv.enefunc(res, req, err)
		}
		ctx.Reset(nil, nil)
		r.pool.Put(ctx)
	}))
	next = FilterChain(filters...)(next)
	next = FilterChain(r.filters...)(next)
	r.srv.router.Handle(path.Join(r.prefix, relativePath), next).Methods(method)
}

func (r *Router) GET(path string, h HandlerFunc, m ...FilterFunc) {
	r.Handle(http.MethodGet, path, h, m...)
}

func (r *Router) POST(path string, h HandlerFunc, m ...FilterFunc) {
	r.Handle(http.MethodPost, path, h, m...)
}
