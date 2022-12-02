package http

import (
	"context"
	"github.com/gorilla/mux"
	"net/http"
	"net/url"
)

type Context interface {
	context.Context
	Vars() url.Values
	//Query() url.Values
	//Form() url.Values
	//Header() http.Header
	//Request() *http.Request
	//Response() http.ResponseWriter
	//Bind(interface{}) error
	//BindVars(interface{}) error
	//BindQuery(interface{}) error
	//BindForm(interface{}) error
	//Returns(interface{}, error) error
	Result(int, interface{}) error
	//JSON(int, interface{}) error
	//XML(int, interface{}) error
	//String(int, string) error
	//Blob(int, string, []byte) error
	//Stream(int, string, io.Reader) error
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

func (c *wrapper) Vars() url.Values {
	raws := mux.Vars(c.req)
	vars := make(url.Values, len(raws))
	for k, v := range raws {
		vars[k] = []string{v}
	}
	return vars
}
