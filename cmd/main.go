package main

import (
	"fmt"

	"crypto-to-fiat-converter/intenral/dep_container"
)

func main() {
	container, err := dep_container.New()
	if err != nil {
		panic(fmt.Sprintf("error initializing DI container: %+v", err))
	}

	container.RunGrpcServer()
}
