package router

import ( 
	"net/http"
	"strconv"
	"database/sql"
	"crudBooks/internal/pkg/web"
	"io/ioutil"
	"errors"
)

type RestFul interface{
	//TODO errors should have a custom struct to include the http response code (default 400) and the error message, then the main router could handle all error responses
	Init(CollectionUrl string, db *sql.DB)
	HandleGet(id int, w http.ResponseWriter) //bad requests handled prior to call + 404 handled in call = no error to return
	HandlePost(w http.ResponseWriter, body []byte) (error)
	HandlePut(id int, w http.ResponseWriter, body []byte) (error)
	HandlePatch(id int, w http.ResponseWriter, body []byte) (error)
	HandleDelete(id int, w http.ResponseWriter) (error)
}



func RestFulSplitter(path string, db *sql.DB, collection RestFul) (func(http.ResponseWriter, *http.Request)){
	return func(w http.ResponseWriter, r *http.Request) {
		id64, err := strconv.ParseInt(r.URL.Path[len(path):], 0, 32)
		if err != nil || id64 < 0 {
			id64 = 0
		}
		id := int(id64) // for functions to be called with system int def
		if id==0 && (r.Method == http.MethodPut || r.Method == http.MethodPatch || r.Method == http.MethodDelete ){
			// Forbid writes that are to the entire collection
			web.RespondWithError(w, 405, "Requested command (" + r.Method + ") not supported on collection. Positive integer item ID required (got '" + r.URL.Path[len(path):] + "').")
		} else {
			// id is non-zero or not required
			// while it won't be used in all subsequent scenarios, it will be in almost all and passing a pointer isn't all that demanding
			collection.Init(path, db) // TODO send w in there too, rather than in every following call
			body, err := ioutil.ReadAll(r.Body)
			if err == nil {
				//Finally! do something specific to each method/verb
				switch r.Method {
					case http.MethodGet:
						// Get all books or a book
						collection.HandleGet(id, w)
					case http.MethodPost:
						if id>0{
							// id is autoincrement, creator cannot set id
							web.RespondWithError(w, 405, "Cannot control ID of newly created records. Nothing done.")
						} else {
							// Create a new book.
							err = collection.HandlePost(w, body)
						}
					case http.MethodPut:
						// Replace an existing book.
						err = collection.HandlePut(id, w, body)
					case http.MethodPatch:
						// Modify an existing book.
						err = collection.HandlePatch(id, w, body)
					case http.MethodDelete:
						// Delete a book.
						err = collection.HandleDelete(id, w)
					default:
						err = errors.New("Requested command (" + r.Method + ") not supported.")
				}
			}

		}
		if (err != nil){
			// Bad request
			web.RespondWithError(w, 400, err.Error());
		}
	}
}