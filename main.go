package main

import (
	"log"
	"net/http"
	"strings"
)

var Ampelfarbe []byte = []byte("Green")


func getcol(w http.ResponseWriter, r *http.Request) {
	w.Write(Ampelfarbe)
	return
}


func ping(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		col := r.FormValue("col")
		boo := r.FormValue("newcol")
		if boo == "yes" {
			Ampelfarbe = []byte(col)
		}
		w.Write(Ampelfarbe)
		return
	}
	//um zu schauen wie so ein html interpretiert wird.
	http.ServeFile(w, r, "form.html")
	return

	//Alte version mit get und parametern.
	//checke die params ab, falls newcol=true, setzt Ampelfarbe auf die farbe 1=grün, 2=gelb, 3=rot
	newcol, _ := r.URL.Query()["newcol"]
	newcolstr := strings.Join(newcol, "")
	col, _ := r.URL.Query()["col"]
	colstr:= strings.Join(col,"")


	if string(newcolstr)=="true" {
		switch colstr {
		case "1":
			Ampelfarbe= []byte("Green")
		case "2":
			Ampelfarbe= []byte("Yellow")
		case "3":
			Ampelfarbe= []byte("Red")
		}
	}


	w.Write(Ampelfarbe)

}

func getFarbe() []byte {
	return Ampelfarbe
}

func setFarbe(farbe []byte) {
	Ampelfarbe = farbe
}

func main() {
	http.HandleFunc("/set", ping)
	http.HandleFunc("/colour", getcol)
	log.Fatal(http.ListenAndServe(":8080", nil))

}

