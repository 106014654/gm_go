package main

import (
	"context"
	"fmt"
	"gm_go/config"
	"gm_go/config/file"
	trs "gm_go/transport/http"
	"log"
	"net/http"
	"time"
)

const appJSONStr = "application/json"

type User struct {
	Name string `json:"name"`
}

func corsFilter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			log.Println("cors:", r.Method, r.RequestURI)
			w.Header().Set("Access-Control-Allow-Methods", r.Method)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func authFilter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		log.Println("auth:", r.Method, r.RequestURI)
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}

func loggingFilter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		log.Println("logging:", r.Method, r.RequestURI)
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}

func main() {
	configDemo()
}

func configDemo() {
	c := config.New(config.WithSource(file.NewSource("src/gm_go/config.yaml")))
	if err := c.Load(); err != nil {
		panic(err)
	}

	var v struct {
		Service struct {
			Name    string `json:"name"`
			Version string `json:"version"`
		} `json:"service"`
	}

	if err := c.Scan(&v); err != nil {
		panic(err)
	}
	fmt.Println(v.Service.Version)

	name, err := c.Value("service.name").String()
	if err != nil {
		panic(err)
	}
	fmt.Println(name)
}

func httpRouterDemo() {
	ctx := context.Background()
	srv := trs.NewServer(
		trs.Filter(corsFilter, loggingFilter),
	)

	//srv.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
	//	writer.WriteHeader(http.StatusOK)
	//	fmt.Fprintf(writer, "hello world xxxxxxxxxxxxxxxxxxx")
	//})

	//
	//route := srv.Route("/v1")
	srv.Route("/v1").GET("/users/{name}", func(ctx trs.Context) error {
		u := new(User)
		u.Name = ctx.Vars().Get("name")
		//fmt.Println(ctx.Vars().Get("name"))

		_ = ctx.String(200, u.Name)
		return nil
		//return ctx.Result(200, u)

	})

	http.ListenAndServe(":8080", srv.Handler)

	if e, err := srv.Endpoint(); err != nil || e == nil {
		log.Fatal(e, err)
	}
	go func() {
		if err := srv.Start(ctx); err != nil {
			panic(err)
		}
	}()

	time.Sleep(30 * time.Second)
	_ = srv.Stop(ctx)
}
