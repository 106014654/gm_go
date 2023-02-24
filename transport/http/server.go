package http

import (
	"context"

	"errors"
	log "github.com/sirupsen/logrus"
	"gm_go/transport"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/gorilla/mux"
)

type ServerOption func(*Server)

type EncodeErrorFunc func(http.ResponseWriter, *http.Request, error)

type EncodeResponseFunc func(http.ResponseWriter, *http.Request, interface{}) error

type Server struct {
	*http.Server
	lis         net.Listener
	endpoint    *url.URL
	err         error
	network     string
	address     string
	timeout     time.Duration
	filters     []FilterFunc
	enefunc     EncodeErrorFunc
	enc         EncodeResponseFunc
	strictSlash bool
	router      *mux.Router
}

func NewServer(opts ...ServerOption) *Server {
	srv := &Server{
		network: "tcp",
		address: ":8080",
		timeout: 10 * time.Second,
		router:  mux.NewRouter(),
	}
	for _, o := range opts {
		o(srv)
	}
	srv.router.StrictSlash(srv.strictSlash)
	srv.router.NotFoundHandler = http.DefaultServeMux
	srv.router.MethodNotAllowedHandler = http.DefaultServeMux
	srv.router.Use(srv.filter())
	srv.Server = &http.Server{
		Handler: FilterChain(srv.filters...)(srv.router),
	}
	return srv
}

func Filter(filters ...FilterFunc) ServerOption {
	return func(o *Server) {
		o.filters = filters
	}
}

func (s *Server) Endpoint() (*url.URL, error) {
	if err := s.listenAndEndpoint(); err != nil {
		return nil, err
	}
	return s.endpoint, nil
}

func (s *Server) listenAndEndpoint() error {
	if s.lis == nil {
		lis, err := net.Listen(s.network, s.address)
		if err != nil {
			s.err = err
			return err
		}
		s.lis = lis
	}
	if s.endpoint == nil {
		addr, err := transport.Extract(s.address, s.lis)
		if err != nil {
			s.err = err
			return err
		}
		s.endpoint = transport.NewEndpoint(transport.Scheme("http", false), addr)
	}
	return s.err
}

func (s *Server) filter() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			var (
				ctx    context.Context
				cancel context.CancelFunc
			)
			if s.timeout > 0 {
				ctx, cancel = context.WithTimeout(req.Context(), s.timeout)
			} else {
				ctx, cancel = context.WithCancel(req.Context())
			}
			defer cancel()

			pathTemplate := req.URL.Path
			if route := mux.CurrentRoute(req); route != nil {
				// /path/123 -> /path/{id}
				pathTemplate, _ = route.GetPathTemplate()
			}

			tr := &Transport{
				operation:    pathTemplate,
				pathTemplate: pathTemplate,
				reqHeader:    headerCarrier(req.Header),
				replyHeader:  headerCarrier(w.Header()),
				request:      req,
			}
			if s.endpoint != nil {
				tr.endpoint = s.endpoint.String()
			}
			tr.request = req.WithContext(transport.NewServerContext(ctx, tr))
			next.ServeHTTP(w, tr.request)
		})
	}
}

func (s *Server) Start(ctx context.Context) error {
	if err := s.listenAndEndpoint(); err != nil {
		return err
	}
	s.BaseContext = func(net.Listener) context.Context {
		return ctx
	}
	log.Infof("[HTTP] server listening on: %s", s.lis.Addr().String())
	var err error

	err = s.Serve(s.lis)

	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	log.Info("[HTTP] server stopping")
	return s.Shutdown(ctx)
}

func (s *Server) Route(prefix string, filters ...FilterFunc) *Router {
	return newRouter(prefix, s, filters...)
}

func (s *Server) Handle(path string, h http.Handler) {
	s.router.Handle(path, h)
}

func (s *Server) HandleFunc(path string, h http.HandlerFunc) {
	s.router.HandleFunc(path, h)
}
