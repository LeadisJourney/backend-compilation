package main

import (
	"net/http"
	"html"
	"fmt"
	"log"
)

func main() {
	http.HandleFunc("/v0.1/ce/status", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("received")
		fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	})

	log.Fatal(http.ListenAndServe(":2222", nil))
}
