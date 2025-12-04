package server

import (
	"context"
	"errors"
	"exc8/pb"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type GRPCService struct {
	pb.UnimplementedOrderServiceServer
	drinks *pb.Drinks      //save drinks and orders in-memory
	orders map[int32]int32 //save drink-id and amount, where drink-id is the key
}

func StartGrpcServer() error {
	// Create a new gRPC server.
	srv := grpc.NewServer()
	// Create grpc service
	grpcService := &GRPCService{
		drinks: &pb.Drinks{Drinks: make([]*pb.Drink, 3)},
		orders: make(map[int32]int32),
	}
	// Prepopulate Drinks
	grpcService.drinks.Drinks[0] = &pb.Drink{Id: 0, Name: "Spritzer", Price: 2, Description: "Wine with soda"}
	grpcService.drinks.Drinks[1] = &pb.Drink{Id: 1, Name: "Beer", Price: 3, Description: "Hagenberger Gold"}
	grpcService.drinks.Drinks[2] = &pb.Drink{Id: 2, Name: "Coffee", Price: 0, Description: "Mifare isn't that secure ;)"}
	// Register our service implementation with the gRPC server.
	pb.RegisterOrderServiceServer(srv, grpcService)
	// Serve gRPC server on port 4000.
	lis, err := net.Listen("tcp", ":4000")
	if err != nil {
		return err
	}
	err = srv.Serve(lis)
	if err != nil {
		return err
	}
	return nil
}

func (s *GRPCService) OrderDrink(ctx context.Context, order *pb.Order) (*wrapperspb.BoolValue, error) {
	// check if the drink id is valid
	if order.Drink.Id >= int32(len(s.drinks.Drinks)) || order.Drink.Id < 0 {
		return &wrapperspb.BoolValue{Value: false}, errors.New("Invalid Drink Id!")
	}
	s.orders[order.Drink.Id] += order.Amount
	return wrapperspb.Bool(true), nil
}

func (s *GRPCService) GetDrinks(ctx context.Context, empty *emptypb.Empty) (*pb.Drinks, error) {
	if s.drinks == nil {
		return nil, errors.New("No Drinks!")
	}
	return s.drinks, nil
}

func (s *GRPCService) GetOrders(ctx context.Context, empty *emptypb.Empty) (*pb.Orders, error) {
	if s.orders == nil {
		return nil, errors.New("No Orders!")
	}
	order_list := pb.Orders{}
	for key, value := range s.orders {
		//append orders to the list
		order_list.Orders = append(order_list.Orders, &pb.Order{Amount: value, Drink: s.drinks.Drinks[key]})
	}

	return &order_list, nil
}
