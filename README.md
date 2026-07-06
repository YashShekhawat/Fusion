# Fusion

> A provider-agnostic AI SDK for Go.

Fusion provides a unified interface for working with multiple AI providers. Instead of learning and maintaining different SDKs for OpenAI, Gemini, Anthropic, and others, you write your application once using Fusion and choose the provider through configuration.

---

## Why Fusion?

Every AI provider exposes its own SDK, request format, response model, and error handling. Supporting multiple providers often results in duplicated code and vendor lock-in.

Fusion solves this by providing:

* A unified API across AI providers
* Common request and response models
* Standardized error handling
* Provider-independent application code
* Extensible architecture for adding new providers

Your application interacts with Fusion—not with individual provider SDKs.

---

## Features

* Unified API
* Provider-agnostic request and response models
* Standardized error handling
* Context-aware request execution
* Middleware support
* Modular driver architecture
* Easy provider switching
* Extensible design

### Planned

* Streaming responses
* Tool / Function Calling
* Embeddings
* Image generation
* Retry support
* Logging improvements
* Additional providers

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
	"fmt"
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

	fmt.Println(resp.Message.Content)
}
```

Switching providers only requires changing the client configuration.

```go
client.WithProvider(client.ProviderOpenAI)
```

Your application code remains unchanged.

---

# Supported Providers

| Provider  | Status |
| --------- | :----: |
| OpenAI    |    ✅   |
| Gemini    |    ✅   |
| Anthropic |   🚧   |
| Groq      |   📋   |
| Ollama    |   📋   |

---

# Project Structure

```text
fusion/
├── client/
├── drivers/
│   ├── gemini/
│   ├── openai/
│   └── ...
├── middleware/
├── models/
├── registry/
├── fusionerrors/
└── examples/
```

---

# Architecture

Fusion follows a layered architecture where every component has a single responsibility.

```text
Application
     │
     ▼
 Client
     │
     ▼
 Registry
     │
     ▼
 Driver
     │
 ┌───┼───────────────┐
 │   │               │
 ▼   ▼               ▼
Request Mapping   Transport   Response Mapping
     │               │               │
     └───────────────┴───────────────┘
                     │
                     ▼
               AI Provider
```

Each layer focuses on one responsibility, making Fusion easier to maintain, test, and extend.

---

# Package Responsibilities

## client

The public entry point of Fusion.

Applications should only interact with the client.

Responsibilities:

* Configure the selected provider
* Register the appropriate built-in driver
* Execute middleware
* Send requests
* Return provider-agnostic responses

Example:

```go
fusionClient, err := client.New(
	client.WithProvider(client.ProviderGemini),
	client.WithAPIKey(apiKey),
)
```

---

## registry

The registry stores all registered drivers.

Internally it behaves like:

```text
gemini  -> GeminiDriver
openai  -> OpenAIDriver
...
```

The client retrieves the appropriate driver from the registry instead of containing provider-specific request logic.

---

## drivers

Each provider implements the same Driver interface.

Example:

```text
drivers/
    gemini/
    openai/
    anthropic/
```

Each driver is independent.

Adding one provider should never require modifying another provider.

---

## models

Contains the common request and response structures used across the SDK.

Applications always work with Fusion models instead of provider-specific models.

---

## middleware

Middleware wraps drivers to provide reusable functionality such as:

* Logging
* Metrics
* Tracing
* Retry (planned)

without modifying provider implementations.

---

## fusionerrors

Fusion converts provider-specific errors into common SDK errors.

Instead of parsing provider messages:

```go
if strings.Contains(err.Error(), "Incorrect API Key") {
    ...
}
```

Applications can simply write:

```go
if errors.Is(err, fusionerrors.ErrUnauthorized) {
    ...
}
```

making error handling consistent across providers.

---

# Driver Structure

Each provider follows the same structure.

```text
drivers/
    provider/
        driver.go
        request.go
        response.go
        transport.go
        errors.go
```

### driver.go

Coordinates the request lifecycle.

Responsibilities:

* Validate input
* Convert Fusion requests
* Call transport
* Convert provider responses

---

### request.go

Converts Fusion models into provider-specific request models.

---

### response.go

Converts provider responses into Fusion responses.

Applications never parse provider JSON directly.

---

### transport.go

Responsible only for communication.

Responsibilities:

* Build HTTP requests
* Send requests
* Decode JSON responses

No business logic belongs here.

---

### errors.go

Maps provider-specific errors into Fusion errors.

Example:

```text
401 Unauthorized

↓

fusionerrors.ErrUnauthorized
```

---

# Design Principles

Fusion follows a few simple principles.

## Single Responsibility

Every package should have one responsibility.

## Separation of Concerns

| Component | Responsibility         |
| --------- | ---------------------- |
| Client    | Public API             |
| Registry  | Driver lookup          |
| Driver    | Provider orchestration |
| Request   | Request mapping        |
| Transport | HTTP communication     |
| Response  | Response mapping       |
| Errors    | Error mapping          |

## Provider Agnostic

Applications should not depend on provider-specific SDKs.

Changing providers should require configuration changes—not application changes.

## Extensibility

Adding a new provider should only require:

1. Creating a new driver.
2. Registering it.
3. Done.

Existing providers should not require modification.

---

# Error Handling

Fusion standardizes provider errors.

Examples include:

* ErrUnauthorized
* ErrTimeout
* ErrRateLimit
* ErrProviderUnavailable
* ErrInvalidResponse

This allows applications to use idiomatic Go error handling regardless of the selected provider.

---

# Adding a New Provider

Every provider should implement the Driver interface.

Recommended structure:

```text
drivers/
    provider/
        driver.go
        request.go
        response.go
        transport.go
        errors.go
```

Steps:

1. Implement the Driver interface.
2. Map Fusion requests.
3. Send requests.
4. Map responses.
5. Map provider errors.
6. Register the driver.

---

# Contributing

Contributions are welcome.

Please follow these guidelines:

## Do

* Keep packages focused.
* Follow the existing driver structure.
* Reuse common models and errors.
* Keep provider-specific logic inside provider packages.
* Write clear documentation.

## Don't

* Put HTTP code inside the client.
* Build provider requests inside transport.
* Parse responses inside the client.
* Couple one provider to another.
* Duplicate logic across providers.

---

# Roadmap

* Streaming
* Tool Calling
* Embeddings
* Image Generation
* Middleware enhancements
* Retry support
* More AI providers
* Additional examples
* Expanded documentation

---

# Guiding Principle

Before adding new code, ask:

> **Does this belong to Fusion, or does it belong to a specific provider?**

If every provider could use it, it probably belongs in Fusion.

If it only applies to one provider, keep it inside that provider's package.

This principle keeps Fusion modular, maintainable, and easy to extend.

---

# License

Fusion is released under the MIT License.
