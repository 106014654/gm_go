package http

import (
	"context"
	"net/http"
)

type Context interface {
	context.Context
	Request() *http.Request
	Response() http.ResponseWriter
	Reset(http.ResponseWriter, *http.Request)
}

type responseWriter struct {
	code int
	w    http.ResponseWriter
}

type wrapper struct {
	router *Router
	req    *http.Request
	res    http.ResponseWriter
	w      responseWriter
}
