package config

import (
	"flag"
	"gm_go/config/file"
)

var flagconf string

func init() {
	flag.StringVar(&flagconf, "conf", "config.yaml", "config path, eg: -conf config.yaml")
}

func main() {
	flag.Parse()

	c := config.New(file.NewSource(flagconf))
	data, err := c.Load()

	println(data, err)

	var v struct {
		Service struct {
			Name    string `json:"name"`
			Version string `json:"version"`
		} `json:"service"`
	}

	// Unmarshal the config to struct
	if err := c.Scan(&v); err != nil {
		panic(err)
	}
	log.Printf("config: %+v", v)
}
