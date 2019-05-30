#!/bin/bash

target="127.0.0.1:8080/books"
curl="`which curl` -sD -"

function myRun(){
    echo $'\n'"==========================================="
    echo "==========================================="
    echo $@
    echo "==========================================="
    $@ | grep -v '^Date'
    echo $'\n'"==========================================="
}

echo "truncate table books" | mysql -p go

#empty list of books
myRun $curl $target/

#404 on book 1
myRun $curl $target/1

#fail to create book
myRun $curl -X POST -H "Content-Type: application/json" -d @incomplete_patch.json $target/1

#fail to create book (trying to demand a specific ID be set)
myRun $curl -X POST -H "Content-Type: application/json" -d @replace_book.json $target/1

#404 on book 1
myRun $curl $target/1

#create first book
myRun $curl -X POST -H "Content-Type: application/json" -d @book.json $target/

#view book 1
myRun $curl $target/1

#modify book 1
myRun $curl -X PATCH -H "Content-Type: application/json" -d @incomplete_book.json $target/1

#list of books (including modified book 1)
myRun $curl $target/

#replace book 1
myRun $curl -X PUT -H "Content-Type: application/json" -d @replace_book.json $target/1

#view book 1
myRun $curl $target/1

#delete book 1
myRun $curl -X DELETE $target/1

#empty list of books
myRun $curl $target/