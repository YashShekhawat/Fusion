package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/YashShekhawat/fusion/client"
	"github.com/YashShekhawat/fusion/drivers/gemini"
	"github.com/YashShekhawat/fusion/middleware"
	"github.com/YashShekhawat/fusion/models"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()
	// Create Fusion client
	logger := log.New(
		os.Stdout,
		"[Fusion] ",
		log.LstdFlags,
	)

	// Create Fusion client with logging middleware.
	fusionClient, err := client.New(
		client.WithMiddleware(
			middleware.Logging(logger),
		),
	)
	if err != nil {
		log.Fatal(err)
	}

	// Read them with os.Getenv
	apiKey := os.Getenv("GEMINI_API_KEY")
	fmt.Println("API Key:", apiKey)

	// Create Gemini driver
	geminiDriver, err := gemini.New(gemini.Config{
		APIKey: apiKey,
	})
	if err != nil {
		log.Fatal(err)
	}

	// Register the driver
	if err := fusionClient.Register(geminiDriver); err != nil {
		log.Fatal(err)
	}

	// Generate a response
	resp, err := fusionClient.Generate(
		context.Background(),
		"gemini",
		models.GenerateRequest{
			Model: "gemini-2.5-flash",
			Messages: []models.Message{
				{
					Role:    models.RoleSystem,
					Content: "You are a helpful AI assistant.",
				},
				{
					Role:    models.RoleUser,
					Content: "Explain Go generics in 1 word",
				},
			},
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Provider :", geminiDriver.Name())
	fmt.Println("Role     :", resp.Message.Role)
	fmt.Println("Response :", resp.Message.Content)
	fmt.Printf(
		"Usage    : Prompt=%d Completion=%d Total=%d\n",
		resp.Usage.PromptTokens,
		resp.Usage.CompletionTokens,
		resp.Usage.TotalTokens,
	)
}
