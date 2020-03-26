package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"net/http"
)

//pointer to the postgresdb running in the background
var db *sql.DB

//information about the postgresql db running in the background.
const (
	host     = "192.168.1.21"
	port     = 5432
	user     = "postgres"
	password = "wallah123"
	dbname   = "ampeldb"
)

//Handler, just giving back the current colour of the ampel
func getcol(w http.ResponseWriter, r *http.Request) {
	//read out the colour from the db
	sqlStatement := `SELECT colour FROM colour`
	res := ""
	_ = db.QueryRow(sqlStatement).Scan(&res)

	//and print the colour to the website.
	http.ServeFile(w, r, "src/"+res+".html")

	return
}

/*Handler to set the ampelcolor, on a get request, a form is printed and when the form is submited
this creates a post reqest also handled by that handler which changes the Ampelfarbe var.
*/
func ping(w http.ResponseWriter, r *http.Request) {
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
			panic(err)
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

//Main function, the webpage responds on /set and / get request and on /set post requests.
func main() {
	//set up postgresql db.
	//create connection string
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+"password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	//set up connection and save the db in global variable
	var err error
	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Connection to db successful")

	//handle the requests
	http.HandleFunc("/set", ping)
	http.HandleFunc("/", getcol)
	log.Fatal(http.ListenAndServe(":8080", nil))

}
