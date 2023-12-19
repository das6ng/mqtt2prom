package main

import (
	"os"

	"github.com/das6ng/mqtt2prom"
)

func main() {
	if err := mqtt2prom.NewApp().Run(os.Args); err != nil {
		os.Exit(1)
	}
}
