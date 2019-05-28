package datastore

import (
//    "log"
//    "net/http"
	"crudBooks/configs/dbConfig"
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
)

func Main() (db *sql.DB){
	//TODO accept an argument to switch between testing/staging/live databases
	conf := dbConfig.Main()
	//TODO verify conf.driver is a supported (imported here/initialized) driver
	db, err := sql.Open(conf.Driver, conf.User+":"+conf.Password+"@tcp("+conf.HostAndPort+")/"+conf.DbName)
	if err != nil {
		panic(err.Error())
		//log.Fatal(err)
	}
	//TODO if using the db only in this function and once per app run: defer db.Close()
	// Test/open connection
	err = db.Ping()
	if err != nil {
		panic(err.Error())
		//log.Fatal(err)
	}
	//TODO do something with database
	//TODO return something from database
	return db
}