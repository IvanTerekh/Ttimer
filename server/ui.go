package server

import "net/http"

func HomeHandler(w http.ResponseWriter, r *http.Request){
	http.ServeFile(w, r, "views/index.html")
}