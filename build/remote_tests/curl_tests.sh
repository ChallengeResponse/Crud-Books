#!/bin/bash

if [ $# -lt 1 ]; then
  echo $'\n'"Usage:"
  echo "$0 targetHostOrIp"
  echo "or"
  echo "$0 targetHostOrIp -q"
  echo "the -q will report a simple result of a checksum against what was expected"
  echo "note: both commands may destroy data"
  exit 1
fi

host=$1
target="$host:8080/books"

curl="`which curl` -sD -"

function myRun(){
    echo $'\n'"==========================================="
    echo "==========================================="
    # Remove host ip for md5 sums to match up between hosts
    echo $@ | sed "s/$host/HOST/"
    echo "==========================================="
    # Remove date for md5 sums to match up between runs
    $@ | grep -v '^Date'
    echo $'\n'"==========================================="
}

if [ "$2" != "--no-truncate" ]; then
  echo $'\n'"Enter MySql password for go database to clear its 'books' table"
  echo $'\n'"aka enter your password to destroy data"
  echo "Truncate table books" | mysql -p go
fi

if [ "$2" == "-q" ]; then
    currMd5=`$0 $1 --no-truncate | md5sum`
    [ "$currMd5" == "248484bcc34158c768d60403ecac89ff  -" ] && echo "Tests passed" || echo "Tests differ, got $currMd5, run without -q"
else
    #empty list of books
    myRun $curl $target/

    #404 on book 1
    myRun $curl $target/1

    #fail to create book
    myRun $curl -X POST -H "Content-Type: application/json" -d @incomplete_patch.json $target/

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

    echo "Run $0 with -q to silently check against a checksum of the tests here"
fi