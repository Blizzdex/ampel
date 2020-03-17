package main

import (
	"log"
	"net/http"
)

//Variable for the current ampel colour
var Ampelfarbe []byte = []byte("Green")

//Handler, just giving back the current colour of the ampel
func getcol(w http.ResponseWriter, r *http.Request) {
	w.Write(Ampelfarbe)
	return
}

/*Handler to set the ampelcolor, on a get request, a form is printed and when the form is submited
this creates a post reqest also handled by that handler which changes the Ampelfarbe var.
*/

func ping(w http.ResponseWriter, r *http.Request) {
	//Handle a post request to set the color
	if r.Method == "POST" {
		col := r.FormValue("col")
		Ampelfarbe = []byte(col)

		w.Write(Ampelfarbe)
		return
	}

	//If it is a get request on the /set, we return the form to fill out.
	http.ServeFile(w, r, "src/setform.html")
	return

}

//Main function, the webpage responds on /set and /colour get request and on /set post requests.
func main() {
	http.HandleFunc("/set", ping)
	http.HandleFunc("/colour", getcol)
	log.Fatal(http.ListenAndServe(":8080", nil))

}
