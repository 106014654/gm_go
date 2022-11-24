package main

import (
	"fmt"
	"gm_go/config"
	"gm_go/config/file"
)

func main() {
	c := config.New(config.WithSource(file.NewSource("./config.yaml")))
	if err := c.Load(); err != nil {
		panic(err)
	}
	/**
	service:
	  name: config
	  version: v1.0.0
	*/

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
