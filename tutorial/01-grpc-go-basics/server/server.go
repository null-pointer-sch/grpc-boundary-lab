package main

import (
	"context"
	"log"
	"net"

	pb "github.com/null-pointer-sch/grpc-boundary-lab/tutorial/01-grpc-go-basics/proto"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedCoffeeShopServer
}

func (s *server) GetMenu(menuRequest *pb.MenuRequest, srv pb.CoffeeShop_GetMenuServer) error {
	items := []*pb.Item{
		{Id: "1", Name: "Black Coffee"},
		{Id: "2", Name: "Americano"},
		{Id: "3", Name: "Vanilla Soy Chai Latte"},
	}

	for i := range items {
		srv.Send(&pb.Menu{
			Items: items[0 : i+1],
		})
	}
	return nil
}

func (s *server) PlaceOrder(ctx context.Context, order *pb.Order) (*pb.Receipt, error) {
	return &pb.Receipt{
		Id: "ABC123",
	}, nil
}

func (s *server) GetOrderStatus(ctx context.Context, receipt *pb.Receipt) (*pb.OrderStatus, error) {
	return &pb.OrderStatus{
		OrderId: receipt.Id,
		Status:  "IN PROGRESS",
	}, nil
}

func main() {
	// setup listener on port 9001
	lis, err := net.Listen("tcp", ":9001")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// create gRPC server
	grpcServer := grpc.NewServer()
	pb.RegisterCoffeeShopServer(grpcServer, &server{}) // most interesting line in this file

	// start server
	log.Println("Server listening on port 9001")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
