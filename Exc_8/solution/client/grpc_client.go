package client

import (
	"context"
	"exc8/pb"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
)

type GrpcClient struct {
	client pb.OrderServiceClient
}

func NewGrpcClient() (*GrpcClient, error) {
	conn, err := grpc.NewClient(":4000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	client := pb.NewOrderServiceClient(conn)
	return &GrpcClient{client: client}, nil
}

func (c *GrpcClient) Run() error {
	ctx := context.Background()
	// 1. List drinks
	fmt.Println("Requesting drinks 🍹🍺☕")
	drinks, err := c.client.GetDrinks(ctx, &emptypb.Empty{})
	if err != nil {
		return err
	}
	fmt.Println("Available drinks:")
	for _, drink := range drinks.Drinks {
		fmt.Println("\t > id:", drink.Id, " name:", drink.Name, "\t price:", drink.Price, "\t description:", drink.Description)
	}

	// 2. Order a few drinks
	fmt.Println("Ordering drinks 👨‍🍳⏱️🍻🍻 ")
	fmt.Println("\t > Ordering 2 x", drinks.Drinks[2].Name)
	_, err = c.client.OrderDrink(ctx, &pb.Order{Drink: drinks.Drinks[2], Amount: 2})
	if err != nil {
		return err
	}
	fmt.Println("\t > Ordering 1 x", drinks.Drinks[0].Name)
	_, err = c.client.OrderDrink(ctx, &pb.Order{Drink: drinks.Drinks[0], Amount: 1})
	if err != nil {
		return err
	}
	fmt.Println("\t > Ordering 2 x", drinks.Drinks[1].Name)
	_, err = c.client.OrderDrink(ctx, &pb.Order{Drink: drinks.Drinks[1], Amount: 2})
	if err != nil {
		return err
	}

	// 3. Order more drinks
	fmt.Println("Ordering another round of drinks 👨‍🍳⏱️🍻🍻")
	fmt.Println("\t > Ordering 5 x", drinks.Drinks[2].Name)
	_, err = c.client.OrderDrink(ctx, &pb.Order{Drink: drinks.Drinks[2], Amount: 5})
	if err != nil {
		return err
	}
	fmt.Println("\t > Ordering 5 x", drinks.Drinks[0].Name)
	_, err = c.client.OrderDrink(ctx, &pb.Order{Drink: drinks.Drinks[0], Amount: 5})
	if err != nil {
		return err
	}
	fmt.Println("\t > Ordering 2 x", drinks.Drinks[1].Name)
	_, err = c.client.OrderDrink(ctx, &pb.Order{Drink: drinks.Drinks[1], Amount: 2})
	if err != nil {
		return err
	}

	// 4. Get order total
	fmt.Println("Getting the bill 💹💹💹 ")
	orders, err := c.client.GetOrders(ctx, &emptypb.Empty{})
	if err != nil {
		return err
	}
	price := float32(0)
	for _, order := range orders.Orders {
		fmt.Println("\t > Total:", order.Amount, "x", order.Drink.Name)
		price += float32(order.Amount) * order.Drink.Price
	}
	fmt.Println("Orders complete!")

	return nil
}
