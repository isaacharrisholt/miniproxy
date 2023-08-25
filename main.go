package main

import (
	"log"

	"github.com/isaacharrisholt/miniproxy/proxy"
)

func main() {
	settings, err := proxy.NewSettings("miniproxy.json")
	if err != nil {
		log.Fatal(err)
	}

	p, err := proxy.NewProxy(settings)
	if err != nil {
		log.Fatal(err)
	}
	p.Start()
}
