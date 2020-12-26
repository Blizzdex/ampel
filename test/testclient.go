package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/golang/protobuf/ptypes/empty"
	pb "gitlab.ethz.ch/vis/cat/ampel2/servis/vseth/vis/ampel"
	"google.golang.org/grpc"
)

var serverAddr = flag.String("server_addr", "localhost:8083", "addr of ampel server")

//set up ampel client to interact with the server.
func main() {
	flag.Parse()
	conn, err := grpc.Dial(*serverAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()
	client := pb.NewAmpelClient(conn)
	//request the colour now.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var null = empty.Empty{}
	colour, err := client.GetColor(ctx, &null)
	var col = colour.Color
	var c = pb.Color_name[int32(col)]
	if err != nil {
		log.Fatalf("%v.GetFeatures(_) = _, %v: ", client, err)
	}
	fmt.Println("Colour is: " + c)
}
