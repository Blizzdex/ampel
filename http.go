package main

import (
	"html/template"
	"net/http"
	"strconv"

	log "github.com/sirupsen/logrus"
)

//enum for the ampelcolours with methods to use it.
type color int32 //why would I need a col_t type?

const (
	GREEN  int = 1
	YELLOW int = 2
	RED    int = 3
)

type col4Temp struct {
	Col string
}

var colorName = map[int]string{
	0: "invalidFormat",
	1: "green",
	2: "yellow",
	3: "red",
}

//Handler, just giving back the current colour of the ampel
func (s server) getColor(w http.ResponseWriter, r *http.Request) {
	//read out the colour from the db
	var res, err = s.DbGetColor()
	//check if color is valid
	if res == 0 {
		w.Write([]byte("Could not display."))
		log.Warn("failed to get color, invalid color.")
		return
	}
	if err != nil {
		w.Write([]byte("Could not display."))
		log.WithError(err).Warn("failed to get ampelcolor")
		return
	}
	var color = colorName[res]
	//and print the colour to the website.
	var p = col4Temp{Col: color}

	//create the template if that has not been done yet.
	if s.t == nil {
		var e error
		s.t, e = template.ParseFiles("src/colTemplate.html")
		if e != nil {
			l.Fatalf("Failed to parse Template")
		}
	}

	s.t.Execute(w, p)
	return
}

/*Handler to set the ampelcolor, on a get request, a form is printed and when the form is submited
this creates a post reqest also handled by that handler which changes the Ampelfarbe var.
*/
func (s server) setColor(w http.ResponseWriter, r *http.Request) {
	//Handle a post request to set the color
	if r.Method == "POST" {
		//get the color from the form
		var col = r.FormValue("col")
		if col == "" {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		var color, err = strconv.ParseInt(col, 10, 32)
		if err != nil {
			log.Warn("Could not change Ampelcolor, invalid input.")
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		//set the color in the db
		s.DbSetColor(int(color))

		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	//If it is a get request on the /set, we return the form to fill out.
	if r.Method == "GET" {
		http.ServeFile(w, r, "src/setform.html")
	}

	return

}
