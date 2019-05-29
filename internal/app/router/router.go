package router

import ( 
	"net/http"
	"strconv"
	"database/sql"
	"crudBooks/internal/pkg/web"
)

type RestFul interface{
	LinkDb(db *sql.DB)
	HandleGet(id int, w http.ResponseWriter)
	HandlePost(w http.ResponseWriter, r *http.Request)
	HandlePut(id int, w http.ResponseWriter, r *http.Request)
	HandlePatch(id int, w http.ResponseWriter, r *http.Request)
	HandleDelete(id int, w http.ResponseWriter)
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
			collection.LinkDb(db)

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
						collection.HandlePost(w, r)
					}
				case http.MethodPut:
					// Replace an existing book.
					collection.HandlePut(id, w, r)
				case http.MethodPatch:
					// Modify an existing book.
					collection.HandlePatch(id, w, r)
				case http.MethodDelete:
					// Delete a book.
					collection.HandleDelete(id, w)
				default:
					web.RespondWithError(w, 400, "Requested command (" + r.Method + ") not supported.")
			}
		}
	}
}