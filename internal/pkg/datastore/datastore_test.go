package datastore

import (
	"testing"
	"fmt"
)

func TestMain(t *testing.T) {
	db := Main("127.0.0.1")
	fmt.Print("Tested DB config & connection!")
	db.Close()
}