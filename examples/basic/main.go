package main

import (
	"context"
	"fmt"

	"github.com/YashShekhawat/fusion/client"
	"github.com/YashShekhawat/fusion/drivers/mock"
	"github.com/YashShekhawat/fusion/models"
)

func main() {
	client := client.New()

	client.Register(mock.New())

	fmt.Println("Application: Calling Client.Generate()")

	resp, err := client.Generate(context.TODO(), "mocks", models.GenerateRequest{
		Messages: []models.Message{
			{
				Role:    models.RoleUser,
				Content: "Hello, how are you?",
			},
		},
	})

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("%+v\n", resp)
}
