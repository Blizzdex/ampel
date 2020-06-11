package main

import (
	log "github.com/sirupsen/logrus"
	"net/http"
)

//enum for the ampelcolours with methods to use it.
type color int32 //why would I need a col_t type?

const (
	GREEN  color = 1
	YELLOW color = 2
	RED    color = 3
)

var color_name = map[color]string{
	0: "invalidFormat",
	1: "green",
	2: "yellow",
	3: "red",
}

func (c color) file() string {
	return string("src/" + color_name[c] + ".html")
}

//Handler, just giving back the current colour of the ampel
func getcol(w http.ResponseWriter, r *http.Request) {
	//read out the colour from the db
	sqlStatement := `SELECT color FROM color`
	var res color
	_ = db.QueryRow(sqlStatement).Scan(&res)
	var col = color_name[res]
	if col == "invalidFormat" {
		w.Write([]byte("Could not display."))
		log.Warn("Failed to get valid AmpelColor.")
		return
	}
	//and print the colour to the website.
	http.ServeFile(w, r, res.file())

	return
}

/*Handler to set the ampelcolor, on a get request, a form is printed and when the form is submited
this creates a post reqest also handled by that handler which changes the Ampelfarbe var.
*/
func setcol(w http.ResponseWriter, r *http.Request) {
	//Handle a post request to set the color
	if r.Method == "POST" {
		col := r.FormValue("col")
		if col == "" {
			return
		}
		//Write the new colour into the db
		sqlStatement := `
			UPDATE color
			SET color = $1
			WHERE id=1`
		_, err := db.Exec(sqlStatement, col)
		if err != nil {
			w.Write([]byte("Could not change Ampelcolour. Retry to set colour!"))
			log.Warn("Could not change Ampelcolour, connection to DB failed. Retrying to connect to DB!")
			connectDB() //can cause program to panic if the connection fails.
			return
		}
		//Write out the new colour to the webpage

		getcol(w, r)
		return
	}

	//If it is a get request on the /set, we return the form to fill out.
	if r.Method == "GET" {
		http.ServeFile(w, r, "src/setform.html")
	}

	return

}
