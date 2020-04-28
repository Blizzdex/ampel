package main

//!!PORT NUMBERS are wrong for testing.

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	_ "github.com/lib/pq"
	pb "gitlab.ch/ampel2/grpc"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
)

//set up the grpc server variables
var portgrpc = flag.Int("port", 8083, "Port for grpc requests")

type ampel2Server struct {
	pb.UnimplementedAmpel2Server
}

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
	//connect to the postgres DB
	connectDB()
	//set up the grpc server
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *portgrpc))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	var grpcServer = grpc.NewServer()
	var serv = ampel2Server{}
	pb.RegisterAmpel2Server(grpcServer, &serv)
	go func() {
		fmt.Println("grpc up")
		grpcServer.Serve(lis)
	}()

	//handle http requests
	fmt.Println("Listening")
	http.HandleFunc("/set", setcol)
	http.HandleFunc("/", getcol)
	log.Fatal(http.ListenAndServe(":8080", nil))

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

func (*ampel2Server) GetColour(ctx context.Context, req *pb.Null) (*pb.Colour, error) {
	//Create the right colour elem.
	sqlStatement := `SELECT colour FROM colour`
	var res string
	_ = db.QueryRow(sqlStatement).Scan(&res)
	var farbe pb.AmpelColour
	var err error
	switch res {
	case "green":
		farbe = 0
	case "yellow":
		farbe = 1
	case "red":
		farbe = 2
	default:
		return nil, err
	}
	return &pb.Colour{Colour: farbe}, nil
}

func (*ampel2Server) SetColour(ctx context.Context, req *pb.Colour) (*pb.Ack, error) {
	var col = req.Colour
	var str = col.String()
	var ack pb.Ack
	ack.Success = true
	sqlStatement := `
			UPDATE colour
			SET colour = $1
			WHERE id=1`
	_, err := db.Exec(sqlStatement, str)
	if err != nil {
		ack.Success = false
		connectDB()
	}
	return &ack, nil

}
