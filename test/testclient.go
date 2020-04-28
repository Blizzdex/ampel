package main

import (
	"context"
	"flag"
	"fmt"
	pb "gitlab.ch/ampel2/grpc"
	"google.golang.org/grpc"
	"log"
	"time"
)

var serverAddr = flag.String("server_addr", "localhost:8083", "addr of grpc server")

//set up grpc client to interact with the server.
func main() {
	flag.Parse()
	conn, err := grpc.Dial(*serverAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()
	client := pb.NewAmpel2Client(conn)
	//request the colour now.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var null = pb.Null{}
	colour, err := client.GetColour(ctx, &null)
	var col = colour.Colour
	var m = pb.AmpelColour_name
	var c = m[int32(col)]
	if err != nil {
		log.Fatalf("%v.GetFeatures(_) = _, %v: ", client, err)
	}
	fmt.Println("Colour is: " + c)
}
