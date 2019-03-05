package main

import (
	"context"
	"flag"
	"log"
	"time"

	"google.golang.org/grpc"

	"github.com/stelo/blackmore/pkg/api/v1"
)

const (
	// apiVersion is version of API is provided by server
	apiVersion = "v1"
)

func main() {
	// get configuration
	address := flag.String("server", "", "gRPC server in format host:port")
	flag.Parse()

	// Set up a connection to the server.
	conn, err := grpc.Dial(*address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := v1.NewLinksellerServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Call Create
	req1 := v1.CreateRequest{
		Api: apiVersion,
		Linkseller: &v1.Linkseller{
			Person:  &v1.Person{Type: "PF", Document: "0203939"},
			Machine: &v1.Machine{Modelcode: 123, Seriesnumber: "123", Value: 123.45, Model: "123", Chip: "123"},
			Order:   &v1.Order{Ordercode: 123},
		},
	}
	res1, err := c.Create(ctx, &req1)
	if err != nil {
		log.Fatalf("Create failed: %v", err)
	}
	log.Printf("Create result: <%+v>\n\n", res1)

	id := res1.Id

	log.Printf("Create Id: <%+v>\n\n", id)

}
