package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/YashShekhawat/fusion/client"
	"github.com/YashShekhawat/fusion/middleware"
	"github.com/YashShekhawat/fusion/models"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()
	logger := log.New(
		os.Stdout,
		"[Fusion] ",
		log.LstdFlags,
	)

	fusionClient, err := client.New(
		client.WithProvider(client.ProviderGemini),
		client.WithAPIKey(os.Getenv("GEMINI_API_KEY")),
		client.WithMiddleware(
			middleware.Logging(logger),
		),
	)
	if err != nil {
		log.Fatal(err)
	}

	stream, err := fusionClient.GenerateStream(
		context.Background(),
		models.GenerateRequest{
			Model: "gemini-2.5-flash",
			Messages: []models.Message{
				{
					Role:    models.RoleUser,
					Content: "Explain Go interfaces in 300 words.",
				},
			},
		},
	)
	if err != nil {
		log.Fatal(err)
	}
	defer stream.Close()

	fmt.Println("Streaming response:\n")

	for {

		chunk, err := stream.Recv()

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatal(err)
		}

		fmt.Print(chunk.Content)
	}

	fmt.Println()
}
