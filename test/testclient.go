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
		Color:                1,
		XXX_NoUnkeyedLiteral: struct{}{},
		XXX_unrecognized:     nil,
		XXX_sizecache:        0,
	}
	var token = "eyJhbGciOiJSUzI1NiIsInR5cCIgOiAiSldUIiwia2lkIiA6ICIyS08wb3pQTWdBV2F2MEZwMXNaYjljajE0SGl1dWk4ek44dkpQMXVsT2k4In0.eyJleHAiOjE2Mzc4Mzg5NTMsImlhdCI6MTYzNzgzNTM1MywianRpIjoiMTRjOTY1ZGMtZjUwMC00NzA0LTkxNzQtOGU3NmVlMzE0NDIyIiwiaXNzIjoiaHR0cDovL2xvY2FsaG9zdDo4MTgwL2F1dGgvcmVhbG1zL1ZTRVRIIiwiYXVkIjpbImxvY2FsLWFtcGVsIiwiYWNjb3VudCJdLCJzdWIiOiI1OTY2NDA4NS0yZjg2LTRjMjktODQxMy0xMGU0OWI4MTZkYzIiLCJ0eXAiOiJCZWFyZXIiLCJhenAiOiJsb2NhbC1hbXBlbCIsInNlc3Npb25fc3RhdGUiOiI0MzZiN2U3Zi1lZTgxLTRjMjEtYTFlYi1jZDBiNDA3MThjMTYiLCJhY3IiOiIxIiwiYWxsb3dlZC1vcmlnaW5zIjpbImh0dHA6Ly9sb2NhbGhvc3QqIiwiaHR0cHM6Ly9sb2NhbGhvc3QqIl0sInJlYWxtX2FjY2VzcyI6eyJyb2xlcyI6WyJkZWZhdWx0LXJvbGVzLXZzZXRoIiwib2ZmbGluZV9hY2Nlc3MiLCJ1bWFfYXV0aG9yaXphdGlvbiJdfSwicmVzb3VyY2VfYWNjZXNzIjp7ImxvY2FsLWFtcGVsIjp7InJvbGVzIjpbImFkbWluIl19LCJhY2NvdW50Ijp7InJvbGVzIjpbIm1hbmFnZS1hY2NvdW50IiwibWFuYWdlLWFjY291bnQtbGlua3MiLCJ2aWV3LXByb2ZpbGUiXX19LCJzY29wZSI6InByb2ZpbGUiLCJzaWQiOiI0MzZiN2U3Zi1lZTgxLTRjMjEtYTFlYi1jZDBiNDA3MThjMTYiLCJuYW1lIjoiQWJpZ2FpbCBTbWl0aCIsInByZWZlcnJlZF91c2VybmFtZSI6ImFiaWdhaWxfc21pdGgiLCJnaXZlbl9uYW1lIjoiQWJpZ2FpbCIsImZhbWlseV9uYW1lIjoiU21pdGgifQ.fmreZgn5b89T4QBGwJXc6TDBo3jeNP8widdiG7MbCjXqgUCH_8X-IyrvR4yb8g2Joi1FY-_wcns5fIB2G5DC3L2tFI91kecSl12W6BJ0ln8tvJjnrQIRkxDGZtvD5xCl6B2_pqggoTr4q28Y41d7abpdNDO2ai8G-33TPsqE7qMCw8aIRyhuhQh4E_VLJM8Lf3e9oK8ViQu94yGiZUNtENVCs6YsVZhhgWLoWxEywtupml0fkM54r7_G02uo9p1igQG5KN0HFLODmzUWbihY0aD0tq-1zi-2Bc_zM5SpRP6nUOMYoo6IXVZLHEZQG-GqMEdmJrirFbr9LLOfKL3ZEQ"
	var ctx2 = metadata.AppendToOutgoingContext(ctx, "authorization", "bearer "+token)
	var meta3, _ = metadata.FromOutgoingContext(ctx2)
	println(meta3)
	_, err = client.UpdateColor(ctx2, &update)

	if err != nil {
		log.Fatal("failed to update color")
	}

	time.Sleep(1000)
}
