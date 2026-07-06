package main

import (
	"context"
	"fmt"
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
		client.WithProvider(client.Gemini),
		client.WithAPIKey(os.Getenv("GEMINI_API_KEY")),
		client.WithMiddleware(
			middleware.Logging(logger),
		),
	)
	if err != nil {
		log.Fatal(err)
	}

	resp, err := fusionClient.Generate(
		context.Background(),
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

	fmt.Println("Provider :", client.ProviderGemini)
	fmt.Println("Role     :", resp.Message.Role)
	fmt.Println("Response :", resp.Message.Content)
	fmt.Printf(
		"Usage    : Prompt=%d Completion=%d Total=%d\n",
		resp.Usage.PromptTokens,
		resp.Usage.CompletionTokens,
		resp.Usage.TotalTokens,
	)
}
