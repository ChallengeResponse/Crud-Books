package datastore

import (
	"testing"
	"fmt"
)

func TestMain(t *testing.T) {
	db := Main()
	fmt.Print("Tested DB config & connection!")
	db.Close()
}