package book

type RestBooksStore struct{
	bookDbConn *sql.DB
}

func (r RestBooksStore) LinkDb(db *sql.DB){
	r.bookDbConn = db
}


func (r RestBooksStore) HandleGet(id int, w http.ResponseWriter){
	
}

func (r RestBooksStore) HandlePost(w http.ResponseWriter, req *http.Request){

}


func (r RestBooksStore) HandlePut(id int, w http.ResponseWriter, req *http.Request){
}


func (r RestBooksStore) HandlePatch(id int, w http.ResponseWriter, req *http.Request){
}


func (r RestBooksStore) HandleDelete(id int, w http.ResponseWriter){
}