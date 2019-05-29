package book

import (
	"database/sql"
)

var bookTable = "books"

type BookInfo struct{
	Id int
	Title, Author, Publisher, PublishDate string
	Rating byte
	Status bool
}

func (b BookInfo) FromDb(db *sql.DB, id int) (error){
	b.Id = id;
	return b.FromDbRow(db.QueryRow("select title, author, publisher, publishDate, rating, status from " + bookTable + " where id = ?", b.Id))
}

func (b BookInfo) FromDbRow(r interface{Scan(dest ...interface{}) error}) (error){
	return r.Scan(&b.Title, &b.Author, &b.Publisher, &b.PublishDate, &b.Rating, &b.Status)
}

func (b BookInfo) SaveToDb(db *sql.DB) (error){
}

//Make sure there's decent/enough info on non-id fields, leave id requirement to other functions
//TODO return error(s) specific to failure reason
func (b BookInfo) IsValid() (bool){
	// allowed: a 0 Id (uninitialized) in case an insert is intended
	// not allowed: sub-zero Id
	if (b.Id < 0){
		return false
	}
	if (b.Rating > 3 || b.Rating < 1){
		return false
	}
	if (len(b.Author)<1 || len(b.Title)<1 || len(b.Publisher) < 1 || len(b.PublishDate) < 8){
		return false
	}
	return true
}