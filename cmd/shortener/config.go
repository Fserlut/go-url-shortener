package main

import "flag"

var (
	serverAddress = ":8080"
	baseReturnUrl = "http://localhost:8080"
)

func parseFlags() {
	flag.StringVar(&serverAddress, "a", serverAddress, "address and port to run server")
	flag.StringVar(&baseReturnUrl, "b", baseReturnUrl, "address return url")

	flag.Parse()
}
