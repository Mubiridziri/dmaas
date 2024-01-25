package main

import (
	"dmaas/internal/app/dmaas"
	"flag"
	"log"
)

var (
	bindAddr string
)

func init() {
	flag.StringVar(&bindAddr, "port", "8080", "listen server port")
}

func main() {
	flag.Parse()

	a := dmaas.New(":" + bindAddr)

	if err := a.Run(); err != nil {
		log.Fatal(err)
	}
}
