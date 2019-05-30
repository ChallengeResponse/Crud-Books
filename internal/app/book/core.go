package book

import (
	"database/sql"
	"encoding/json"
	"strconv"
	"errors"
)

var bookTable = "books"

type BookInfo struct{
	Id int
	Title, Author, Publisher, PublishDate string
	Rating byte
	Status bool
}

func (b *BookInfo) FromDb(db *sql.DB, id int) (error){
	b.Id = id
	//row :=  db.QueryRow("select title, author, publisher, publishDate, rating, status from " + bookTable + " where id = ?", b.Id)
	//return row.Scan(&b.Title, &b.Author, &b.Publisher, &b.PublishDate, &b.Rating, &b.Status)
	return b.FromDbRow(db.QueryRow("select title, author, publisher, publishDate, rating, status from " + bookTable + " where id = ?", b.Id))
}

func (b *BookInfo) FromDbRow(r interface{Scan(dest ...interface{}) error}) (error){
	return r.Scan(&b.Title, &b.Author, &b.Publisher, &b.PublishDate, &b.Rating, &b.Status)
}

func (b *BookInfo) FromJson(reqBody []byte) (error){
	err := json.Unmarshal(reqBody, &b)
	if (err != nil){
		return err
	}
	if b.IsValid(){
		return nil
	}
	return errors.New("Valid Json, but incomplete or otherwise not within spec for a book.")
}

func (b *BookInfo) SaveToDb(db *sql.DB) (int, error){
	if b.IsValid(){
		var id int
		var sqlStr string
		if (b.Id > 0){
			// If there is an id, saving to the db is assumed to be an update
			// as b.Id is an integer, it cannot contain any sql commands/injection and
			// that lets the insert and update commands have the same number of arguments
			sqlStr = "UPDATE " + bookTable + " SET title = ?, author = ?, publisher = ?, publishDate = ?, rating = ?, status = ? WHERE id = " + strconv.Itoa(b.Id)
			id = b.Id
		} else {
			sqlStr = "INSERT INTO " + bookTable + "(title, author, publisher, publishDate, rating, status) VALUES(?,?,?,?,?,?)"
		}
		stmt, err := db.Prepare(sqlStr)
		if err != nil {
			return 0, err
		}
		// TODO Verify, rather than simply assuming, concurrency/thread safety of core db object when looking for LastInsertId
		res, err := stmt.Exec(b.Title, b.Author, b.Publisher, b.PublishDate, b.Rating, b.Status)
		if err != nil {
			return 0, err
		}
		if (b.Id == 0){
			id64, err := res.LastInsertId()
			if err != nil {
				return 0, err
			}
			// TODO make sure id not > max of current system's int type... not sure of the go constants
			b.Id = int(id64)
		}
		return id, nil
	}
	return 0, errors.New("Book does not meet spec or is incomplete.")
}

//Make sure there's decent/enough info on non-id fields, leave id requirement to other functions
//TODO return error(s) specific to failure reason
func (b *BookInfo) IsValid() (bool){
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