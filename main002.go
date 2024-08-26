package main

import (
	"log"
	"net/http"
)

func main002() {
	dir := http.Dir("./static")
	fs := http.FileServer(dir)
	mux := http.NewServeMux()
	mux.Handle("/", fs)
	err := http.ListenAndServe(":5555", mux)
	if err != nil {
		log.Fatal(err)
	}
}
