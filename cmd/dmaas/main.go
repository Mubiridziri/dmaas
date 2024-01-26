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

	app := dmaas.New(":" + bindAddr)

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
