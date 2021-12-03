package main

import (
	"context"
	"flag"
	"fmt"
	"google.golang.org/grpc/metadata"
	"log"
	"time"

	"github.com/golang/protobuf/ptypes/empty"
	pb "gitlab.ethz.ch/vseth/0403-isg/libraries/protostub-golang/vseth/vis/ampel"
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
		Color: 2,
	}
	var token = "eyJhbGciOiJSUzI1NiIsInR5cCIgOiAiSldUIiwia2lkIiA6ICIyS08wb3pQTWdBV2F2MEZwMXNaYjljajE0SGl1dWk4ek44dkpQMXVsT2k4In0.eyJleHAiOjE2Mzg1NzA5MjksImlhdCI6MTYzODU2NzMyOSwianRpIjoiNjEyMDk1Y2QtZTc3Zi00YzQ5LTg0NDktYzg5NGQzMTJhYjgyIiwiaXNzIjoiaHR0cDovL2xvY2FsaG9zdDo4MTgwL2F1dGgvcmVhbG1zL1ZTRVRIIiwiYXVkIjpbImxvY2FsLWFtcGVsIiwiYWNjb3VudCJdLCJzdWIiOiI1OTY2NDA4NS0yZjg2LTRjMjktODQxMy0xMGU0OWI4MTZkYzIiLCJ0eXAiOiJCZWFyZXIiLCJhenAiOiJsb2NhbC1hbXBlbCIsInNlc3Npb25fc3RhdGUiOiJiZWRjYTcxZi0zY2ViLTQ2Y2EtYWFkMi0zNTgzNTE4MGEzN2IiLCJhY3IiOiIxIiwiYWxsb3dlZC1vcmlnaW5zIjpbImh0dHA6Ly9sb2NhbGhvc3QqIiwiaHR0cHM6Ly9sb2NhbGhvc3QqIl0sInJlYWxtX2FjY2VzcyI6eyJyb2xlcyI6WyJkZWZhdWx0LXJvbGVzLXZzZXRoIiwib2ZmbGluZV9hY2Nlc3MiLCJ1bWFfYXV0aG9yaXphdGlvbiJdfSwicmVzb3VyY2VfYWNjZXNzIjp7ImxvY2FsLWFtcGVsIjp7InJvbGVzIjpbImFkbWluIl19LCJhY2NvdW50Ijp7InJvbGVzIjpbIm1hbmFnZS1hY2NvdW50IiwibWFuYWdlLWFjY291bnQtbGlua3MiLCJ2aWV3LXByb2ZpbGUiXX19LCJzY29wZSI6InByb2ZpbGUiLCJzaWQiOiJiZWRjYTcxZi0zY2ViLTQ2Y2EtYWFkMi0zNTgzNTE4MGEzN2IiLCJuYW1lIjoiQWJpZ2FpbCBTbWl0aCIsInByZWZlcnJlZF91c2VybmFtZSI6ImFiaWdhaWxfc21pdGgiLCJnaXZlbl9uYW1lIjoiQWJpZ2FpbCIsImZhbWlseV9uYW1lIjoiU21pdGgifQ.QZxpkbnQfRTwYXyhmioeOJBWIQvl_Oe8OFP5D60hQGs5Z6FKAw5kh3ZWON3Tw-nX7Tj1AKyzwTCcOH0HflQmKYuvfc__6bu4OpoOuJAE0taFxqn-Fia6EJ8fbCexh7iDXm-HzDtIwtF1pjWZjk2f3PIvdNsKzu_hZJNmcfg83vU0GGILPoP2aCPuPgIzBeOXEnghiOvnB-dFCTWiHTsbq5GxMwSiACgr9grMbAEsPkaLVckqK-Z_XASopPyf5uTWlks2freTNvPVwLUSAgubNV8k8wt09kh2bFyZAaReJ4UUhIQ3-eXIEDCA0t52Xek04SFPxSApHmzFw5kaTVXYmQ"
	var ctx2 = metadata.AppendToOutgoingContext(ctx, "authorization", "bearer "+token)
	var meta3, _ = metadata.FromOutgoingContext(ctx2)
	println(meta3)
	_, err = client.UpdateColor(ctx2, &update)

	if err != nil {
		log.Fatal("failed to update color")
	}

	time.Sleep(1000)
}
