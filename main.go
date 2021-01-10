package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"net"
	"net/http"
	"strings"

	"github.com/golang/protobuf/ptypes/empty"
	_ "github.com/lib/pq"
	migrate "github.com/rubenv/sql-migrate"
	log "github.com/sirupsen/logrus"
	pb "gitlab.ethz.ch/vis/cat/ampel2/servis/vseth/vis/ampel"
	"google.golang.org/grpc"
)

//set up the ampel server variables
var (
	portgrpc    = flag.Int("port", 7777, "Port for grpc ampel requests")
	postgresURL = flag.String("postgres-url", "", "(required) example: myuser:mypass@172.17.0.2:5432/drinks_registry?sslmode=disable")
	db          *sql.DB //pointer to the postgresdb
	l           *log.Logger
)

type server struct {
}

func checkArgs() {
	flag.Parse()
	for k, v := range map[string]string{"postgres-url": *postgresURL} {
		if strings.HasPrefix("{{", v) {
			l.Fatalf("missing required argument %v:\n", k)
		}
	}
}

//Main function, the webpage responds on /set and / get request and on /set post requests.
func main() {
	checkArgs()

	//Set up the logger
	l = log.New()
	l.SetReportCaller(false)
	l.SetFormatter(&log.JSONFormatter{})
	//connect to the postgres DB
	var err error
	db, err = connectDB()
	if err != nil {
		l.WithError(err).Fatal("failed to connect to db")
	}
	l.Println("Connection to db successful")

	//apply the migrations
	var migrations = migrate.FileMigrationSource{Dir: "migrations"}
	migCount, err := migrate.Exec(db, "postgres", migrations, migrate.Up)
	if err != nil {
		l.WithError(err).Fatal("failed to migrate")
	}
	l.Printf("applied %v migrations\n", migCount)

	//set up the ampel server
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *portgrpc))
	if err != nil {
		l.WithError(err).Fatal("failed to listen")
	}
	var grpcServer = grpc.NewServer()
	var serv = server{}
	pb.RegisterAmpelServer(grpcServer, &serv)
	go func() {
		l.Info("ampel up")
		l.Fatal(grpcServer.Serve(lis))
	}()

	//set up file server
	var fileServer = http.FileServer(http.Dir("."))
	http.Handle("/resources/", http.StripPrefix("/resources", fileServer))
	//handle http requests
	l.Println("Listening")
	http.HandleFunc("/set", setcol)
	http.HandleFunc("/", getcol)
	l.Fatal(http.ListenAndServe(":80", nil))

}

//func to connect to the Database if connection fails, the program will log or return the error
func connectDB() (*sql.DB, error) {
	//set up postgresql db.
	var dbp, err = sql.Open("postgres", fmt.Sprintf("postgres://%v", *postgresURL))
	if err != nil {
		return dbp, err
	}
	err = dbp.Ping()
	return dbp, err
}

//The handlers for grpc requests.
func (*server) GetColor(ctx context.Context, req *empty.Empty) (*pb.GetColorResponse, error) {
	//Create the right colour elem.
	sqlStatement := `SELECT color FROM color`
	var color int
	var err = db.QueryRow(sqlStatement).Scan(&color)

	return &pb.GetColorResponse{Color: pb.Color(color)}, err
}

func (*server) UpdateColor(ctx context.Context, req *pb.UpdateColorRequest) (*pb.Ack, error) {
	var col = req.Color
	var ack pb.Ack
	ack.Success = true
	sqlStatement := `
			UPDATE color
			SET color = $1
			WHERE id=1`
	_, err := db.Exec(sqlStatement, col)
	if err != nil {
		ack.Success = false
		l.Warn("Failed to set colour")
		connectDB()
	}
	return &ack, nil

}
