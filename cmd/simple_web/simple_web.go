package main

import (
	"golang_layout/internal/handler/page_handler"
	"golang_layout/internal/repo/wiki_db"
	"log"
	"net/http"
)

func main() {
	wiki := wiki_db.WikiRepo{}
	wiki.Open()
	page_handler.CreateHandlers() //create http handlers for all web directory

	defer wiki.Close()

	log.Fatal(http.ListenAndServe(":8080", nil)) //serve till fatal error or ctrl^c
}

//initialize configs, dll
