package main

import (
	"net/http"
)

func main() {
	http.HandleFunc("/home", func(writer http.ResponseWriter, request *http.Request) {})
	http.HandleFunc("/index", func(writer http.ResponseWriter, request *http.Request) {})

	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
