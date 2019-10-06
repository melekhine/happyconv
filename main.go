package main

import (
	"log"
	"net/http"
)

func main() {

	srv := NewServer()
	log.Fatal(http.ListenAndServe("localhost:8080", srv))

}
