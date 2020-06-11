package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"github.com/golang/protobuf/ptypes/empty"
	_ "github.com/lib/pq"
	migrate "github.com/rubenv/sql-migrate"
	log "github.com/sirupsen/logrus"
	pb "gitlab.ch/ampel2/ampel"
	"google.golang.org/grpc"
	"net"
	"net/http"
	"strings"
)

//set up the ampel server variables
var (
	portgrpc    = flag.Int("port", 8083, "Port for ampel requests")
	postgresURL = flag.String("postgres-url", "", "(required) example: myuser:mypass@172.17.0.2:5432/drinks_registry?sslmode=disable")
	db          *sql.DB //pointer to the postgresdb
	setup       bool    = false
	l           *log.Logger
)

type ampel2Server struct {
	pb.UnimplementedAmpel2Server
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
	l.SetReportCaller(true)
	l.SetFormatter(&log.JSONFormatter{})
	//connect to the postgres DB
	connectDB()
	//apply the migrations
	var migrations = migrate.FileMigrationSource{Dir: "migrations"}
	migCount, err2 := migrate.Exec(db, "postgres", migrations, migrate.Up)
	if err2 != nil {
		l.Fatalf("failed to migrate: %v\n", err2)
	}
	l.Println("applied %v migrations\n", migCount)

	//set up the ampel server
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *portgrpc))
	if err != nil {
		l.Fatalf("failed to listen: %v", err)
	}
	var grpcServer = grpc.NewServer()
	var serv = ampel2Server{}
	pb.RegisterAmpel2Server(grpcServer, &serv)
	go func() {
		l.Info("ampel up")
		grpcServer.Serve(lis)
	}()

	//handle http requests
	l.Println("Listening")
	http.HandleFunc("/set", setcol)
	http.HandleFunc("/", getcol)
	l.Fatal(http.ListenAndServe(":80", nil))

}

//func to connect to the Database if connection fails, the program will panic
func connectDB() {
	//set up postgresql db.
	var err error
	db, err = sql.Open("postgres", fmt.Sprintf("postgres://%v", *postgresURL))
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	l.Println("Connection to db successful")
}

func (*ampel2Server) GetColor(ctx context.Context, req *empty.Empty) (*pb.Col, error) {
	//Create the right colour elem.
	sqlStatement := `SELECT color FROM color`
	var farbe int
	_ = db.QueryRow(sqlStatement).Scan(&farbe)

	return &pb.Col{Color: pb.Color(farbe)}, nil
}

func (*ampel2Server) SetColor(ctx context.Context, req *pb.Col) (*pb.Ack, error) {
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
