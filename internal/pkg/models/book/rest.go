package book

import ( 
	str "strconv"
	"net/http"
	"database/sql"
	"crudBooks/internal/pkg/web"
)

type RestBooksStore struct{
	bookDbConn *sql.DB
}

func (r RestBooksStore) LinkDb(db *sql.DB){
	r.bookDbConn = db
}


func (r RestBooksStore) HandleGet(id int, w http.ResponseWriter){
	// Bad requests (error 400) should have been filtered out already, but 404 may happen for some books
	books := make([]BookInfo, 0)
	if (id > 0){
		// Select by Id
		var book BookInfo
		err := book.FromDb(r.bookDbConn,id)
		if err == sql.ErrNoRows{
			web.RespondWithError(w, 404, "Requested book (" + str.Itoa(id) + ") not found.")
			return
		} else if err != nil{
			panic(err.Error())
		} else {
			books = append(books,book)
		}
	} else {
		//TODO pagination / limit for larger collections
		rows, err := r.bookDbConn.Query("select * from " + bookTable, 1)
		defer rows.Close() //TODO necessary with no rows or other error?		
		if err != nil  && err != sql.ErrNoRows{
			panic(err.Error())
		} else if err != sql.ErrNoRows{
			for rows.Next() {
				var book BookInfo
				err := book.FromDbRow(rows)
				if err != nil {
					panic(err.Error())
				}
				books = append(books,book)
			}
			err = rows.Err()
			if err != nil {
				panic(err.Error())
			}
		} // else  err == sql.ErrNoRows, but books is already an empty slice
	}
	// 400s handled, and books either has a book or the collection (possibly empty) of all books
	web.RespondwithJSON(w, 200, books)
}

// create new, then return a 201 with a location header that points at the new resource
// per parent router, non-nil error returns will be converted into a response 400 with message
func (r RestBooksStore) HandlePost(w http.ResponseWriter, body []byte) (error){
	var book BookInfo
	err := book.FromJson(body)
	if err != nil{
		return err
	}
	err := book.SaveToDb(r.bookDbConn)
	if err != nil{
		return err
	}
	// TODO set location header and respond with 201 + empty body
	return nil
}


// replace an existing resource.  404 if it does not exist. Return an error if the request is badly formed.
func (r RestBooksStore) HandlePut(id int, w http.ResponseWriter, body []byte) (error){
}


func (r RestBooksStore) HandlePatch(id int, w http.ResponseWriter, body []byte) (error){
}


func (r RestBooksStore) HandleDelete(id int, w http.ResponseWriter) (error){
}