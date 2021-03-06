package router

import ( 
	"net/http"
	"strconv"
	"crudBooks/internal/pkg/web"
	"io/ioutil"
	"errors"
)

type RestFul interface{
	//TODO errors should have a custom struct to include the http response code (default 400) and the error message, then the main router could handle all error responses
	UrlPath() (string)
	HandleGet(id int, w http.ResponseWriter) //bad requests handled prior to call + 404 handled in call = no error to return
	HandlePost(body []byte, w http.ResponseWriter) (error)
	HandlePut(id int, body []byte, w http.ResponseWriter) (error)
	HandlePatch(id int, body []byte, w http.ResponseWriter) (error)
	HandleDelete(id int, w http.ResponseWriter) (error)
}

func RestFulSplitter(collection RestFul) (func(http.ResponseWriter, *http.Request)){
	return func(w http.ResponseWriter, r *http.Request) {
		var WebErr error
		id64, err := strconv.ParseInt(r.URL.Path[len(collection.UrlPath()):], 0, 32)
		if err != nil || id64 < 0 {
			id64 = 0
		}
		id := int(id64) // for functions to be called with system int def
		if id==0 && (r.Method == http.MethodPut || r.Method == http.MethodPatch || r.Method == http.MethodDelete ){
			// Forbid writes that are to the entire collection
			web.RespondWithError(w, 405, "Requested command (" + r.Method + ") not supported on collection. Positive integer item ID required (got '" + r.URL.Path[len(collection.UrlPath()):] + "').")
			return
		} else {
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
							WebErr = collection.HandlePost(body, w)
						}
					case http.MethodPut:
						// Replace an existing book.
						WebErr = collection.HandlePut(id, body, w)
					case http.MethodPatch:
						// Modify an existing book.
						WebErr = collection.HandlePatch(id, body, w)
					case http.MethodDelete:
						// Delete a book.
						WebErr = collection.HandleDelete(id, w)
					default:
						WebErr = errors.New("Requested command (" + r.Method + ") not supported.")
				}
			}
		}
		if (WebErr != nil){
			// Bad request
			web.RespondWithError(w, 400, err.Error());
		}
	}
}