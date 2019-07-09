package main

import (
	"net/http"
)

func main() {
	handlerFunc := func(writer http.ResponseWriter, request *http.Request) {}

	http.HandleFunc("/home", handlerFunc)
	http.HandleFunc("/index", handlerFunc)
	http.HandleFunc("/index/home", handlerFunc)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
