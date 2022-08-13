package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	apiv1 "github.com/daichimukai/x/syakyo/proglog/api/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	addr := flag.String("addr", ":8400", "service address")
	flag.Parse()

	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	client := apiv1.NewLogClient(conn)

	ctx := context.Background()
	res, err := client.GetServers(ctx, &apiv1.GetServersRequest{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("servers:")
	for _, server := range res.Servers {
		fmt.Printf("- %v\n", server)
	}
}
