package main

import "net/http"

//enum for the ampelcolours with methods to use it.
type col_t string

const (
	green         col_t = "green"
	yellow        col_t = "yellow"
	red           col_t = "red"
	invalidFormat col_t = "iv"
)

func (c col_t) file() string {
	return string("src/" + c + ".html")
}

func toCol(s string) col_t {
	var c col_t
	switch s {
	case "green":
		return green
	case "yellow":
		return yellow
	case "red":
		return red
	default:
		return invalidFormat
	}
	return c
}

//Handler, just giving back the current colour of the ampel
func getcol(w http.ResponseWriter, r *http.Request) {
	//read out the colour from the db
	sqlStatement := `SELECT colour FROM colour`
	var res string
	_ = db.QueryRow(sqlStatement).Scan(&res)
	var col = toCol(res)
	if col == invalidFormat {
		w.Write([]byte("Current ampel colour invalid, could not display."))
		return
	}
	//and print the colour to the website.
	http.ServeFile(w, r, col.file())

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
			UPDATE colour
			SET colour = $1
			WHERE id=1`
		_, err := db.Exec(sqlStatement, col)
		if err != nil {
			connectDB()
			w.Write([]byte("Could not change Ampelcolour, connection to DB failed. Retry to set colour!"))
			return
		}
		//Write out the new colour to the webpage
		w.Write([]byte(col))
		return
	}

	//If it is a get request on the /set, we return the form to fill out.
	if r.Method == "GET" {
		http.ServeFile(w, r, "src/setform.html")
	}

	return

}
