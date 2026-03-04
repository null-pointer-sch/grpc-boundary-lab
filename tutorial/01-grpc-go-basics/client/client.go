package main

import (
	"context"
	"io"
	"log"
	"time"

	pb "github.com/null-pointer-sch/grpc-boundary-lab/tutorial/01-grpc-go-basics/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient(":9001", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewCoffeeShopClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	menuStream, err := client.GetMenu(ctx, &pb.MenuRequest{})
	if err != nil {
		log.Fatalf("failed to get menu: %v", err)
	}

	done := make(chan bool)
	var items []*pb.Item

	go func() {
		for {
			resp, err := menuStream.Recv()
			if err == io.EOF {
				done <- true
				break
			}
			if err != nil {
				log.Fatalf("failed to receive menu: %v", err)
			}
			items = resp.Items
			log.Printf("Menu: %v", resp)
		}
	}()
	<-done

	receipt, err := client.PlaceOrder(ctx, &pb.Order{Items: items})
	log.Printf("Receipt: %v", receipt)

	status, err := client.GetOrderStatus(ctx, receipt)
	log.Printf("Status: %v", status)
}
