package main

import (
	"log"
	"net/http"
)

const (
	addr = "localhost:8008"
	dir  = ".."
)

func main() {
	log.Fatal(http.ListenAndServe(addr, http.FileServer(http.Dir(dir))))
}
