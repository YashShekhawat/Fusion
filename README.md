# Fusion

Fusion is a provider-agnostic Go SDK for building applications with Large Language Models (LLMs).

Instead of writing provider-specific code for every AI service, Fusion provides a single, consistent API that works across multiple providers. Switching providers should require changing configuration—not rewriting your application.

## Features

* ✅ Provider-agnostic API
* ✅ Built-in Gemini support
* ✅ Streaming responses
* ✅ Middleware support
* ✅ Canonical error handling
* ✅ Extensible driver architecture
* ✅ Simple configuration using functional options

---

# Installation

```bash
go get github.com/YashShekhawat/fusion
```

---

# Quick Start

```go
package main

import (
	"context"
	"log"
	"os"

	"github.com/YashShekhawat/fusion/client"
	"github.com/YashShekhawat/fusion/models"
)

func main() {
	fusionClient, err := client.New(
		client.WithProvider(client.ProviderGemini),
		client.WithAPIKey(os.Getenv("GEMINI_API_KEY")),
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
					Role:    models.RoleUser,
					Content: "Explain Go interfaces.",
				},
			},
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(resp.Message.Content)
}
```

---

# Streaming

Fusion also supports streaming responses.

```go
stream, err := fusionClient.GenerateStream(ctx, req)
if err != nil {
	log.Fatal(err)
}
defer stream.Close()

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
```

`Recv()` blocks until the next chunk is available.

When the stream has finished, `Recv()` returns `io.EOF`.

Always call `Close()` after you're done reading from the stream.

---

# Configuration

Fusion uses functional options to configure the client.

```go
fusionClient, err := client.New(
	client.WithProvider(client.ProviderGemini),
	client.WithAPIKey(os.Getenv("GEMINI_API_KEY")),
)
```

### Available Options

#### Provider

```go
client.WithProvider(client.ProviderGemini)
```

#### API Key

```go
client.WithAPIKey(os.Getenv("GEMINI_API_KEY"))
```

#### Middleware

```go
client.WithMiddleware(
	middleware.Logging(logger),
)
```

Options can be combined:

```go
fusionClient, err := client.New(
	client.WithProvider(client.ProviderGemini),
	client.WithAPIKey(os.Getenv("GEMINI_API_KEY")),
	client.WithMiddleware(
		middleware.Logging(logger),
	),
)
```

---

# Middleware

Fusion supports middleware that wraps provider drivers.

Current middleware includes:

* Logging

Middleware applies to both:

* `Generate()`
* `GenerateStream()`

Example:

```go
logger := log.New(os.Stdout, "[Fusion] ", log.LstdFlags)

fusionClient, err := client.New(
	client.WithProvider(client.ProviderGemini),
	client.WithAPIKey(os.Getenv("GEMINI_API_KEY")),
	client.WithMiddleware(
		middleware.Logging(logger),
	),
)
```

---

# Project Structure

```
fusion/
├── client/
├── drivers/
│   ├── gemini/
│   └── stream.go
├── fusionerrors/
├── middleware/
├── models/
├── registry/
└── examples/
    ├── gemini/
    └── streaming/
```

---

# Architecture

Fusion follows a provider-agnostic architecture.

Normal request flow:

```
Application
      │
      ▼
Client
      │
      ▼
Middleware
      │
      ▼
Registry
      │
      ▼
Driver
      │
      ▼
Provider
```

Streaming request flow:

```
Application
      │
      ▼
Client.GenerateStream()
      │
      ▼
Middleware
      │
      ▼
StreamDriver
      │
      ▼
Stream
      │
      ▼
Provider
```

Provider-specific request and response formats are isolated inside each driver.

---

# Design Principles

Fusion is built around a few simple principles:

* Applications should use a single public API regardless of the underlying provider.
* Provider-specific logic belongs inside drivers.
* Middleware should operate on interfaces instead of concrete implementations.
* Streaming is an optional driver capability, allowing providers to support it independently.
* Common concepts such as requests, responses, and errors are normalized across providers.

---

# Examples

The repository includes complete working examples.

```
examples/
├── gemini/
└── streaming/
```

The Gemini example demonstrates standard text generation.

The streaming example demonstrates incremental response generation using `GenerateStream()`.

---

# Extending Fusion

Fusion is designed to support additional providers without changing the core client.

Advanced users can register custom drivers using:

```go
client.Register(driver)
```

This is primarily intended for:

* Custom providers
* Internal company integrations
* Testing
* Community-maintained drivers

Most applications will only need to configure the provider and API key using `client.New()`.

---

# Contributing

Contributions are welcome.

When contributing:

* Keep provider-specific logic inside the corresponding driver.
* Preserve the provider-agnostic public API.
* Add tests for new functionality.
* Keep middleware generic so it works across providers.

Before opening a pull request, ensure all examples compile and existing tests pass.

---

# License

MIT License.