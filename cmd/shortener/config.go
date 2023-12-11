package main

import "flag"

var (
	serverAddress = ":8080"
	baseReturnURL = "http://localhost:8080"
)

func parseFlags() {
	flag.StringVar(&serverAddress, "a", serverAddress, "address and port to run server")
	flag.StringVar(&baseReturnURL, "b", baseReturnURL, "address return url")

	flag.Parse()
}
