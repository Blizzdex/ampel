package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	_ "gitlab.ch/ampel2/grpc"
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

//Main function, the webpage responds on /set and / get request and on /set post requests.
func main() {
	//connect to the DB
	connectDB()

	//handle the requests
	fmt.Println("Listening")
	http.HandleFunc("/set", setcol)
	http.HandleFunc("/", getcol)
	log.Fatal(http.ListenAndServe(":80", nil))

}

//func to connect to the Database if connection fails, the program will panic
func connectDB() {
	//set up postgresql db.
	//create connection string
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+"password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	//set up connection and save the db in global variable
	var err error
	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Connection to db successful")
}
