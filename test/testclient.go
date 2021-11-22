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
	var token = "eyJhbGciOiJSUzI1NiIsInR5cCIgOiAiSldUIiwia2lkIiA6ICIyQnpsT1lpOWZzRTNlX2MzTUw3OTluVEszZVUteFpZZjBGUWk0b0t3bThjIn0.eyJleHAiOjE2Mzc1MjA1NzEsImlhdCI6MTYzNzUxNjk3MSwianRpIjoiMWFhMDliODYtNTI5ZS00YzBiLWJlYTAtODdhZGIyMzk4MjYxIiwiaXNzIjoiaHR0cDovL2xvY2FsaG9zdDo4MTgwL2F1dGgvcmVhbG1zL1ZTRVRIIiwiYXVkIjoiYWNjb3VudCIsInN1YiI6ImRiY2Q3Y2E2LTU5YjctNDc4ZC1hNDUwLTk1MjQ2OWUyMThhMyIsInR5cCI6IkJlYXJlciIsImF6cCI6ImxvY2FsLWFtcGVsIiwic2Vzc2lvbl9zdGF0ZSI6IjdjY2IxY2NjLTM1ZjUtNDQyZi05OWM3LWVhNzE2M2RiNDZiNiIsImFjciI6IjEiLCJhbGxvd2VkLW9yaWdpbnMiOlsiaHR0cDovL2xvY2FsaG9zdCoiLCJodHRwczovL2xvY2FsaG9zdCoiXSwicmVhbG1fYWNjZXNzIjp7InJvbGVzIjpbImRlZmF1bHQtcm9sZXMtdnNldGgiLCJvZmZsaW5lX2FjY2VzcyIsInVtYV9hdXRob3JpemF0aW9uIl19LCJyZXNvdXJjZV9hY2Nlc3MiOnsiYWNjb3VudCI6eyJyb2xlcyI6WyJtYW5hZ2UtYWNjb3VudCIsIm1hbmFnZS1hY2NvdW50LWxpbmtzIiwidmlldy1wcm9maWxlIl19fSwic2NvcGUiOiJwcm9maWxlIiwic2lkIjoiN2NjYjFjY2MtMzVmNS00NDJmLTk5YzctZWE3MTYzZGI0NmI2IiwibmFtZSI6IkFiaWdhaWwgU21pdGgiLCJwcmVmZXJyZWRfdXNlcm5hbWUiOiJhYmlnYWlsX3NtaXRoIiwiZ2l2ZW5fbmFtZSI6IkFiaWdhaWwiLCJmYW1pbHlfbmFtZSI6IlNtaXRoIn0.cfugYSz7eTDOCViJ2rPvtEHSH8mOfYnnaR0FoYITDiWMV0sSfcffqKcSGpCMgAc8IvHTl_UcGcQQgmBgRxfpFhave89DCM6M_Ip465ZHBDQJoExL6GrRbjayaB2B-zzzmHgpd6r8j6HUzAgtfiakS8nikkIgHFNSHonWxA7vqrz7pt6YHKBqVPa5r0BGM_3NKZ12InbykAUusX8XqyK6WG71htQu5G0HQJ0oJzdzWKfq8kqKFWb2b2ZNdYt2GDcxXzTMYZp96R_Y_wzFJxGWT23-calFAKFxpTbpC6JUZsDNcQ_6_nyccT_2mimWUDHLgTxVvBG9P5I8XAbDDWD5-w"
	var ctx2 = metadata.AppendToOutgoingContext(ctx, "authorization", "bearer "+token)
	var meta3, _ = metadata.FromOutgoingContext(ctx2)
	println(meta3)
	_, err = client.UpdateColor(ctx2, &update)

	if err != nil {
		log.Fatal("failed to update color")
	}

	time.Sleep(1000)
}
