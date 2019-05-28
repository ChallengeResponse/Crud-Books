package main

import (
	"fmt"
	"crudBooks/internal/pkg/datastore"
)

func main(){
	var tableName string
	db := datastore.Main()
	defer db.Close()
	rows, err := db.Query("show tables")
	if err != nil {
		panic(err.Error())
		//log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&tableName)
		if err != nil {
			panic(err.Error())
			//log.Fatal(err)
		}
		fmt.Println(tableName)
	}
	err = rows.Err()
	if err != nil {
		panic(err.Error())
		//log.Fatal(err)
	}
}