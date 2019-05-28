package main

import (
	"fmt"
	"crudBooks/internal/pkg/datastore"
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
	var tableName string
	// container parent address = host.docker.internal
	db := datastore.Main("127.0.0.1")
	defer db.Close()
	rows, err := db.Query("show tables")
	if err != nil {
		panic(err.Error())
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&tableName)
		if err != nil {
			panic(err.Error())
		}
		fmt.Println(tableName)
	}
	err = rows.Err()
	if err != nil {
		panic(err.Error())
	}
	
}