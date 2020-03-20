package main

import (
	"log"
	"net/http"
	"os"

	//"bufio"
	"fmt"
	//"io"
	"io/ioutil"
	//"os"
)

//Variable for the current ampel colour

//Handler, just giving back the current colour of the ampel
func getcol(w http.ResponseWriter, r *http.Request) {
	dat, err := ioutil.ReadFile("src/Ampelcolour.txt")
	if err != nil {
		panic(err)
	}
	fmt.Print(string(dat) + "/n")

	http.ServeFile(w, r, "src/"+string(dat)+".html")

	return
}

/*Handler to set the ampelcolor, on a get request, a form is printed and when the form is submited
this creates a post reqest also handled by that handler which changes the Ampelfarbe var.
*/

func ping(w http.ResponseWriter, r *http.Request) {
	//Handle a post request to set the color
	if r.Method == "POST" {
		col := r.FormValue("col")
		//write the new ampel colour to the file
		file, err := os.Create("src/Ampelcolour.txt")

		if err != nil {
			panic(err)
		}
		defer file.Close()
		var Ampelfarbe = []byte(col)

		file.Write(Ampelfarbe)
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
	http.HandleFunc("/", getcol)
	log.Fatal(http.ListenAndServe(":8080", nil))

}
