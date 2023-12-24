package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	handlerFunc := func(writer http.ResponseWriter, request *http.Request) {}
	handlerWith404BodyFunc := func(writer http.ResponseWriter, request *http.Request) {
		pageName := request.URL.Path
		// respond with 200 intentionally
		writer.WriteHeader(http.StatusOK)

		_, err := writer.Write([]byte(fmt.Sprintf("404: page %s was not found", pageName)))
		if err != nil {
			log.Printf("Error writing response: %v", err)
		}
	}
	handlerWith403BodyFunc := func(writer http.ResponseWriter, request *http.Request) {
		// respond with 200 intentionally
		writer.WriteHeader(http.StatusOK)
		_, err := writer.Write([]byte("403: forbidden"))
		if err != nil {
			log.Printf("Error writing response: %v", err)
		}
	}

	http.HandleFunc("/noHome", handlerWith404BodyFunc)
	http.HandleFunc("/forbiddenHome", handlerWith403BodyFunc)
	http.HandleFunc("/home", handlerFunc)
	http.HandleFunc("/index", handlerFunc)
	http.HandleFunc("/index/home", handlerFunc)

	if err := http.ListenAndServe("127.0.0.1:7999", nil); err != nil {
		panic(err)
	}
}
