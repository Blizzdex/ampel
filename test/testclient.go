package main

import (
	"context"
	"flag"
	"fmt"
	"google.golang.org/grpc/metadata"
	"log"
	"time"

	"github.com/golang/protobuf/ptypes/empty"
	pb "gitlab.ethz.ch/vis/cat/ampel2/servis/vseth/vis/ampel"
	"google.golang.org/grpc"
)

var serverAddr = flag.String("server_addr", "localhost:7777", "addr of ampel server")

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

	//try to change the ampel color
	var update = pb.UpdateColorRequest{
		Color:                3,
		XXX_NoUnkeyedLiteral: struct{}{},
		XXX_unrecognized:     nil,
		XXX_sizecache:        0,
	}
	var token = "eyJhbGciOiJSUzI1NiIsInR5cCIgOiAiSldUIiwia2lkIiA6ICIyS08wb3pQTWdBV2F2MEZwMXNaYjljajE0SGl1dWk4ek44dkpQMXVsT2k4In0.eyJleHAiOjE2Mzc3NjU4MjIsImlhdCI6MTYzNzc2MjIyMiwianRpIjoiYWFkYjhhZDktOWM0Ny00NjEyLWI3ZjYtMTI0M2E1Mjk2ZGM2IiwiaXNzIjoiaHR0cDovL2xvY2FsaG9zdDo4MTgwL2F1dGgvcmVhbG1zL1ZTRVRIIiwiYXVkIjpbImxvY2FsLWFtcGVsIiwiYWNjb3VudCJdLCJzdWIiOiI0NWNmODVhZi04NzQxLTRhY2EtYmI3MS1hMjg4MDc5OTYxMzkiLCJ0eXAiOiJCZWFyZXIiLCJhenAiOiJsb2NhbC1hbXBlbCIsInNlc3Npb25fc3RhdGUiOiJlNTc5ODhiMS1lMzFlLTRlNTAtOGFjMi1mYTkxNWQwYTg5NzkiLCJhY3IiOiIxIiwiYWxsb3dlZC1vcmlnaW5zIjpbImh0dHA6Ly9sb2NhbGhvc3QqIiwiaHR0cHM6Ly9sb2NhbGhvc3QqIl0sInJlYWxtX2FjY2VzcyI6eyJyb2xlcyI6WyJkZWZhdWx0LXJvbGVzLXZzZXRoIiwib2ZmbGluZV9hY2Nlc3MiLCJ1bWFfYXV0aG9yaXphdGlvbiJdfSwicmVzb3VyY2VfYWNjZXNzIjp7ImFjY291bnQiOnsicm9sZXMiOlsibWFuYWdlLWFjY291bnQiLCJtYW5hZ2UtYWNjb3VudC1saW5rcyIsInZpZXctcHJvZmlsZSJdfX0sInNjb3BlIjoicHJvZmlsZSIsInNpZCI6ImU1Nzk4OGIxLWUzMWUtNGU1MC04YWMyLWZhOTE1ZDBhODk3OSIsIm5hbWUiOiJBZGRpc29uIE1pbGxlciIsInByZWZlcnJlZF91c2VybmFtZSI6ImFkZGlzb25fbWlsbGVyIiwiZ2l2ZW5fbmFtZSI6IkFkZGlzb24iLCJmYW1pbHlfbmFtZSI6Ik1pbGxlciJ9.LgSg8GLJeabdKyMP4Bmz05G4zMQ4tBQ7EPsQBw9BKRbpukfwNbggRuztBZZo1kcG2nVv20fV7gT7JuzgTa9pEIphbk0yP0zt339UYfINMJZsTno8_hj6XAHwFSknMIvmroP4gdwyM1RGNx3hChGisMDleF5gZyHOhfSHRPevkw2COww9iBE7-sYhrTQhYdmyQJ6ypi8UaHbUu2CxQNoSSXeyca6Q00iHatJ5WH28BkjMVDJKJtBja8jNQIWuEBdEkFUnpdLvDNJCt-CWDUADnhrix4f7tY8EDj9iqb-XWuJJojcmYDqUB5f2-KbpGkly-JdRqy94sY2l2BC3vn2yAw"
	var ctx2 = metadata.AppendToOutgoingContext(ctx, "authorization", "bearer "+token)
	var meta3, _ = metadata.FromOutgoingContext(ctx2)
	println(meta3)
	_, err = client.UpdateColor(ctx2, &update)

	if err != nil {
		log.Fatal("failed to update color")
	}

	time.Sleep(1000)
}
