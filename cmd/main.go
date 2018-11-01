package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func main() {
	addressF := flag.String("address", "127.0.0.1:3030", "Address listen (host and port).")

	flag.Parse()

	log.Println("Listen address: ", *addressF)
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(*addressF, nil))
}
