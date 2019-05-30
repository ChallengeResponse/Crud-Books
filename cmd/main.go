package main

import (
	"log"
	"net/http"
	"crudBooks/internal/pkg/datastore"
	"crudBooks/internal/pkg/router"
	"crudBooks/internal/app/book"
)

func main(){
	/*
		POST /books
		GET /books
		GET /books/{id}
		PUT /books/{id}
		PATCH /books/{id}
		DELETE /books/{id}

		where a book has:
		- Title
		- Author
		- Publisher
		- Publish Date
		- Rating (1-3)
		- Status (CheckedIn, CheckedOut)
	*/
	// container parent address = host.docker.internal
	db := datastore.Main("127.0.0.1")
	defer db.Close()
	path := "/books/"
	collection := new(book.RestBooksStore)
	collection.Init(path, db)
	http.HandleFunc(path, router.RestFulSplitter(collection))
    log.Fatal(http.ListenAndServe(":8080", nil))	
}