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

/*
Streaming response

[Fusion] 2026/07/07 01:26:48 [gemini] Stream started
Streaming response:

Chunk received (203 chars)
Go interfaces are a fundamental and powerful concept for achieving polymorphism and decoupling in Go. They define a **contract** based on a set of method signatures, without specifying any implementation

Chunk received (222 chars)
**Definition:** An interface is declared with the `interface` keyword, listing the methods it requires:

    type Greeter interface {
        Greet() string
    }

Chunk received (246 chars)
**Implicit Implementation:** The key difference from many other languages is that types **implicitly** implement an interface. A concrete type satisfies an interface if it provides *all* the methods defined by that interface, with the correct signatures. There is no explicit `implements` keyword.

    type Person struct { Name string }
    func (p Person) Greet() string { return "Hello, " + p.Name }

Chunk received (181 chars)
    type Robot struct { Model string }
    func (r Robot) Greet() string { return "Greetings, human from " + r.Model }

    // Both Person and Robot implicitly implement Greeter

Chunk received (232 chars)
**Polymorphism:** This implicit implementation allows you to write functions that accept an interface type. These functions can then operate on *any* concrete type that satisfies that interface, promoting highly reusable and flexible code:

    func sayHello(g Greeter) {
        fmt.Println(g.Greet())
    }

Chunk received (185 chars)
    // Usage:
    p := Person{Name: "Alice"}
    r := Robot{Model: "Unit 734"}
    sayHello(p) // Output: Hello, Alice
    sayHello(r) // Output: Greetings, human from Unit 734

Chunk received (194 chars)
**The Empty Interface (`interface{}` or `any`):**
The empty interface `interface{}` (aliased as `any` since Go 1.18) defines no methods. Consequently, *every* concrete type implicitly implements it. It's used when you need to accept values of any type, often requiring type assertions or type switches to discover the underlying concrete type.

Chunk received (232 chars)
In essence, Go interfaces enable you to define *what an object can do* rather than *what it is*, leading to cleaner, more modular, and easily testable code.

[Fusion] 2026/07/07 01:26:54 [gemini] Stream completed in 6.2209504s
*/
