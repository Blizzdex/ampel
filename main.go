package main

import (
	"log"
	"net/http"
)

var Ampelfarbe []byte = []byte("Green")

func getcol(w http.ResponseWriter, r *http.Request) {
	w.Write(Ampelfarbe)
	return
}

func ping(w http.ResponseWriter, r *http.Request) {
	//Handle a post request to set the color
	if r.Method == "POST" {
		col := r.FormValue("col")
		Ampelfarbe = []byte(col)

		w.Write(Ampelfarbe)
		return
	}
	//If it is a get request on the /set, we return the form to fill out.
	http.ServeFile(w, r, "setform.html")
	return

}

func main() {
	http.HandleFunc("/set", ping)
	http.HandleFunc("/colour", getcol)
	log.Fatal(http.ListenAndServe(":8080", nil))

}
