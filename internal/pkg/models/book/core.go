package book

import "encoding/json"

bookTable =: "books"

type BookInfo struct{
	Id int
	Title, Author, Publisher, PublishDate string
	Rating byte
	Status bool
}

func (b BookInfo) FromDb(db *sql.DB, id int) (error){
	b.id = id;
	return db.QueryRow(
		"select title, author, publisher, publishDate, rating, status from " + bookTable + " where id = ?", 
		b.id
		).Scan(&b.Title, &b.Author, &b.Publisher, &b.PublishDate, &b.Rating, &b.Status)
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
	if (b.rating > 3 || b.rading < 1){
		return false
	}
	if (strlen(b.author)<1 || strlen(b.title)<1 || strlen(b.publisher) < 1 || strlen(b.publishDate) < 8){
		return false
	}
	return true
}