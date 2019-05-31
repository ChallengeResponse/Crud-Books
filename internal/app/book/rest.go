package book

import ( 
	"strconv"
	"errors"
	"net/http"
	"database/sql"
	"crudBooks/internal/pkg/web"
)

type RestBooksStore struct{
	db *sql.DB
	urlPath string
}

func (r *RestBooksStore) Init(CollectionUrlPath string, db *sql.DB){
	r.db = db
	r.urlPath = CollectionUrlPath
}

func (r RestBooksStore) UrlPath() (string){
	return r.urlPath
}

func (r RestBooksStore) loadOr404(id int, w http.ResponseWriter) (book *BookInfo){
	book = nil
	var bookLoader BookInfo
	err := bookLoader.FromDb(r.db,id)
	if err == sql.ErrNoRows{
		web.RespondWithError(w, 404, "Requested book (" + strconv.Itoa(id) + ") not found.")
	} else if err != nil{
		panic(err.Error())
	} else {
		book = &bookLoader
	}
	return
}

func (r RestBooksStore) HandleGet(id int, w http.ResponseWriter){
	// Bad requests (error 400) should have been filtered out already, but 404 may happen for some books
	books := make([]BookInfo, 0)
	if (id > 0){
		// Select by Id
		book := r.loadOr404(id, w);
		if book == nil {
			return
		}
		books = append(books,*book)
	} else {
		//TODO pagination / limit for larger collections
		rows, err := r.db.Query("select * from " + bookTable)
		defer rows.Close() //TODO necessary with no rows or other error?
		if err != sql.ErrNoRows{
			if err != nil{
				panic(err.Error())
			} else { // err == nil
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
			}
		} // else  err == sql.ErrNoRows, but books is already an empty slice
	}
	// 400s handled, and books either has a book or the collection (possibly empty) of all books
	web.RespondwithJSON(w, 200, books)
}

// create new, then return a 201 with a location header that points at the new resource
// per parent router, non-nil error returns will be converted into a response 400 with message
func (r RestBooksStore) HandlePost(body []byte, w http.ResponseWriter) (error){
	var book BookInfo
	err := book.FromJson(body)
	if err != nil{
		return err
	}
	id, err := book.SaveToDb(r.db)
	if err != nil{
		return err
	}
	w.Header().Set("Location", r.urlPath + strconv.Itoa(id))
	w.WriteHeader(201)
	return nil
}

// replace an existing resource.  404 if it does not exist. Return an error if the request is badly formed.
func (r RestBooksStore) HandlePut(id int, body []byte, w http.ResponseWriter) (error){
	var newInfo BookInfo
	// first check the request is valid before hitting the database
	err := newInfo.FromJson(body)
	if err != nil{
		return err
	}
	oldInfo := r.loadOr404(id, w);
	if oldInfo != nil {
		if (newInfo.Id != oldInfo.Id){
			return errors.New("Resource Id mismatch between json and url.")
		}
		_, err := newInfo.SaveToDb(r.db)
		if err != nil{
			return err
		}
		w.WriteHeader(204)
	}
	return nil
}


// Supports a send what you need concept rather than deal with application/json-patch+json
func (r RestBooksStore) HandlePatch(id int, body []byte, w http.ResponseWriter) (error){
	book := r.loadOr404(id, w)
	if book != nil {
		// Unmarshal json onto the already loaded book...
		err := book.FromJson(body)
		if err != nil{
			return err
		}
		if book.Id != id{
			return errors.New("Cannot change ID. Request included change from " + strconv.Itoa(id) + "(url) to " + strconv.Itoa(book.Id) + "(json).")
		}
		_, err = book.SaveToDb(r.db)
		if err != nil{
			return err
		}
		w.WriteHeader(204)
	}
	return nil
}


func (r RestBooksStore) HandleDelete(id int, w http.ResponseWriter) (error){
	// id is an integer, it cannot contain any sql commands/injection
	res, err := r.db.Exec("DELETE FROM " + bookTable + " WHERE id = " + strconv.Itoa(id))
	if err != nil {
		return err
	}
	rowCnt, err := res.RowsAffected()
	if err != nil {
		web.RespondWithError(w, 500, "Internal error after deleting book (" + strconv.Itoa(id) + ") " + err.Error())
	} else if rowCnt == 0 {
		web.RespondWithError(w, 404, "Requested book (" + strconv.Itoa(id) + ") not found.")
	} else {
		w.WriteHeader(204)
	}
	return nil
}