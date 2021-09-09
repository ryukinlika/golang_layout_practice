package main

import (
	"golang_layout/internal/handler/page_handler"
	"log"
	"net/http"
)

func main() {
	page_handler.CreateHandlers() //create http handlers for all web directory

	log.Fatal(http.ListenAndServe(":8080", nil)) //serve till fatal error or ctrl^c
}

//initialize configs, dll
