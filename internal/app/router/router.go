package router

import ( 
	"net/http"
	"strconv"
	"crudBooks/internal/pkg/web"
)

type RestFul interface{
	HandleGet(id int, w http.ResponseWriter)
	HandlePost(w http.ResponseWriter, r *http.Request)
	HandlePut(id int, w http.ResponseWriter, r *http.Request)
	HandlePatch(id int, w http.ResponseWriter, r *http.Request)
	HandleDelete(id int, w http.ResponseWriter)
}



func RestFulSplitter(string path, collection RestFul){
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseInt(r.URL.Path[len(path):], 0, 64)
		if err != nil || id < 0 {
			id = 0
		}
		if id==0 && (r.Method == http.MethodPut || r.Method == http.MethodPatch || r.Method == http.MethodDelete ){
			// Forbid writes that are to the entire collection
			respondWithError(w, 405, "Requested command (" + r.Method + ") not supported on collection. Positive integer item ID required (got '" + r.URL.Path[len(path):] + "').")
		} else {
			// id is non-zero or not required
			switch r.Method {
				case http.MethodGet:
					// Get all books or a book
					collection.HandleGet(id, w)
				case http.MethodPost:
					if id>0{
						// id is autoincrement, creator cannot set id
						respondWithError(w, 405, "Cannot control ID of newly created records. Nothing done.")
					} else {
						// Create a new book.
						HandlePost(w, r)
					}
				case http.MethodPut:
					// Replace an existing book.
					HandlePut(id, w, r)
				case http.MethodPatch:
					// Modify an existing book.
					HandlePatch(id, w, r)
				case http.MethodDelete:
					// Delete a book.
					HandleDelete(id, w)
				default:
					respondWithError(w, 400, "Requested command (" + r.Method + ") not supported.")
			}
		}
	}
}