package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/golang/protobuf/ptypes/empty"
	pb "gitlab.ch/ampel2/ampel"
	"google.golang.org/grpc"
	"log"
	"time"
)

var serverAddr = flag.String("server_addr", "localhost:8083", "addr of ampel server")

//set up ampel client to interact with the server.
func main() {
	flag.Parse()
	conn, err := grpc.Dial(*serverAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	client := pb.NewAmpel2Client(conn)
	time.Sleep(1000)
	//request the colour now.
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	time.Sleep(1000)
	var null = empty.Empty{}
	colour, err := client.GetColor(ctx, &null)
	if err != nil {
		log.Fatalf("%v.GetFeatures(_) = _, %v: ", client, err)
	}
	var col = colour.Color
	var c = pb.Color_name[int32(col)]
	fmt.Println("Colour is: " + c)
	time.Sleep(1000)
}
