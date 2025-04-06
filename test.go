package main

import (
	"fmt"
	"net/http"
)

func HomeHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Println("Hello world")
	w.Write([]byte("Hello world "))
}

func main() {
	fmt.Println("STARTED")
	http.HandleFunc("/", HomeHandler)
	http.ListenAndServe("0.0.0.0:8000", nil)
}
