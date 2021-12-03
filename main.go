package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"html/template"
	"net"
	"net/http"
	"strings"

	oidc "github.com/coreos/go-oidc"
	"github.com/golang/protobuf/ptypes/empty"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	_ "github.com/lib/pq"
	migrate "github.com/rubenv/sql-migrate"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	pb "gitlab.ethz.ch/vseth/0403-isg/libraries/protostub-golang/vseth/vis/ampel"
)

//set up the ampel server variables
var (
	portgrpc       = flag.Int("port", 7777, "Port for grpc ampel requests")
	postgresURL    = flag.String("postgres-url", "", "(required) example: myuser:mypass@172.17.0.2:5432/drinks_registry?sslmode=disable")
	oidc_client_id = flag.String("client-id", "", "sets the client id of our ampel")
	l              *log.Logger
)

type server struct {
	pb.UnimplementedAmpelServer
	db *sql.DB //pointer to the postgresdb
	t  *template.Template
}

//function parses the args coming from cinit and checks for empty args.
func checkArgs() {
	flag.Parse()
	for k, v := range map[string]string{"postgres-url": *postgresURL, "client-id": *oidc_client_id} {
		if strings.HasPrefix("{{", v) {
			l.Fatalf("missing required argument %v:\n", k)
		}
	}
}

//Main function, the webpage responds on /set and / get request and on /set post requests.
func main() {
	//define server struct
	var serv = server{}

	//read program args
	checkArgs()

	//Set up the logger
	l = log.New()
	l.SetReportCaller(false)
	l.SetFormatter(&log.JSONFormatter{})
	//connect to the postgres DB
	var err error
	serv.db, err = connectDB()
	if err != nil {
		l.WithError(err).Fatal("failed to connect to db")
	}
	l.Println("Connection to db successful")

	//apply the migrations
	var migrations = migrate.FileMigrationSource{Dir: "migrations"}
	migCount, err := migrate.Exec(serv.db, "postgres", migrations, migrate.Up)
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
	http.HandleFunc("/set", serv.setColor)
	http.HandleFunc("/", serv.getColor)
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

//function used to extract color from db
func (s server) DbGetColor() (int, error) {
	sqlStatement := `SELECT color FROM color`
	var color int
	var err = s.db.QueryRow(sqlStatement).Scan(&color)
	if color < 1 || color > 3 {
		log.Warn("Failed to get valid AmpelColor.")
		return 0, err
	}
	return color, err
}

//check validity of new color and update color if it is valid
func (s server) DbSetColor(color int) error {
	if color < 1 || color > 3 {
		return nil
	}

	sqlStatement := `
			UPDATE color
			SET color = $1
			WHERE id=1`
	_, err := s.db.Exec(sqlStatement, color)
	if err != nil {
		l.Warn("Failed to set colour")
	}
	return err
}

//The handlers for grpc requests.
func (s *server) GetColor(ctx context.Context, _ *empty.Empty) (*pb.GetColorResponse, error) {
	//Create the right color elem.
	var color, err = (*s).DbGetColor()
	return &pb.GetColorResponse{Color: pb.Color(color)}, err
}

func (s *server) UpdateColor(ctx context.Context, req *pb.UpdateColorRequest) (*empty.Empty, error) {
	//read out authentication token from context
	rawIDToken, err := grpc_auth.AuthFromMD(ctx, "bearer")
	if err != nil {
		l.Warn("Authentication token missing.")
		return &empty.Empty{}, err
	}
	//check request authenticiy
	//TODO save provider + verifier in server (move code to main)
	provider, err := oidc.NewProvider(ctx, "http://172.18.0.1:8180/auth/realms/VSETH") //localhost statt dessen http://172.18.0.1:8180/auth/realms/VSETH
	if err != nil {
		l.Warn("Failed to fetch public key to authenticate grpc UpdateColor.")
		l.Warn(err.Error())
		return &empty.Empty{}, err
	}

	//l.Warn(*oidc_client_id)
	var verifier = provider.Verifier(&oidc.Config{ClientID: *oidc_client_id, SkipIssuerCheck: true, SkipClientIDCheck: true})

	idToken, err := verifier.Verify(ctx, rawIDToken)
	if err != nil {
		l.Warn("Failed to verify access token.")
		l.Warn(err.Error())
		return &empty.Empty{}, err
	}
	//now need to extract roles from claims in the token and then check if contains admin.
	var t Token
	t.token = idToken
	var role = "admin"

	_, hasRole := Find(t.Roles(*oidc_client_id), role)
	if !hasRole {
		l.Warn("Failed due to insufficient permissions of user.")
		return &empty.Empty{}, nil
	}

	var color = int(req.Color)
	var error = s.DbSetColor(color)
	return &empty.Empty{}, error

}
