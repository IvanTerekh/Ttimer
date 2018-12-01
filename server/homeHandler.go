package server

import "net/http"

// HomeHandler serves home html file.
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "views/index.html")
}
