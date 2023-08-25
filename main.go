package main

import (
	"log"

	"github.com/isaacharrisholt/tinyproxy/proxy"
)

func main() {
	settings, err := proxy.NewSettings("tinyproxy.json")
	if err != nil {
		log.Fatal(err)
	}

	p, err := proxy.NewProxy(settings)
	if err != nil {
		log.Fatal(err)
	}
	p.Start()
}
