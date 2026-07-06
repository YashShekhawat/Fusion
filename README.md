# Fusion

> A provider-agnostic AI SDK for Go.

Fusion lets you integrate multiple AI providers through a single, consistent API. Instead of learning and maintaining different SDKs for every provider, write your application once and switch providers with minimal or no application changes.

---

## Why Fusion?

Every AI provider exposes a different SDK, request format, response structure, and error model. Supporting multiple providers often leads to duplicated code and vendor lock-in.

Fusion solves this by providing:

* A unified API across AI providers
* Provider-independent request and response models
* Consistent error handling
* A modular driver architecture
* Easy extensibility for new providers

Your application only needs to know about Fusion—not the underlying provider.

```go
response, err := client.Generate(ctx, request)
```

Whether the request is handled by OpenAI, Gemini, Anthropic, or another provider is managed internally by Fusion.

---

## Features

* Unified API across providers
* Provider-agnostic request and response models
* Standardized error handling
* Context-aware request execution
* Modular driver architecture
* Easy provider registration
* Extensible by design

### Planned

* Streaming responses
* Tool / Function Calling
* Embeddings
* Middleware
* Retry support
* Logging
* Additional AI providers

---

## Installation

```bash
go get github.com/<your-username>/fusion
```

---

## Quick Start

```go
package main

import (
    "context"
    "fmt"

    "github.com/<your-username>/fusion/client"
    "github.com/<your-username>/fusion/models"
)

func main() {
    c := client.New(...)

    resp, err := c.Generate(context.Background(), models.GenerateRequest{
        Model: "gpt-4.1-mini",
        Messages: []models.Message{
            {
                Role: "user",
                Content: "Explain Go interfaces in simple terms.",
            },
        },
    })
    if err != nil {
        panic(err)
    }

    fmt.Println(resp.Text)
}
```

---

## Architecture

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

For a complete explanation of the architecture and design decisions, see **docs/ARCHITECTURE.md**.

---

## Project Structure

```text
fusion/
├── client/
├── drivers/
│   ├── openai/
│   ├── gemini/
│   └── ...
├── models/
├── registry/
├── fusionerrors/
├── docs/
└── ...
```

---

## Supported Providers

| Provider  | Status |
| --------- | ------ |
| OpenAI    | ✅      |
| Gemini    | ✅      |
| Anthropic | 🚧     |
| Groq      | 🚧     |
| Ollama    | 📋     |

> More providers will be added over time without changing the public API.

---

## Error Handling

Fusion converts provider-specific errors into common SDK errors.

Instead of parsing provider error messages:

```go
if strings.Contains(err.Error(), "Incorrect API Key") {
    ...
}
```

Use standard Go error handling:

```go
if errors.Is(err, fusionerrors.ErrUnauthorized) {
    ...
}
```

This keeps applications portable across providers.

---

## Design Philosophy

Fusion is built around a few simple principles:

* Keep the public API simple.
* Hide provider-specific complexity.
* Separate responsibilities into focused packages.
* Make adding providers straightforward.
* Prefer composition over duplication.
* Follow idiomatic Go.

---

## Documentation

* `docs/ARCHITECTURE.md` — Internal architecture and design
* API documentation *(coming soon)*
* Examples *(coming soon)*

---

## Contributing

Contributions are welcome.

If you're interested in adding a new provider or improving the SDK, please read the architecture documentation before getting started.

The goal is to keep Fusion modular, maintainable, and provider-agnostic.

---

## Roadmap

* Streaming
* Tool Calling
* Embeddings
* Middleware
* Logging
* Retry support
* More providers
* Additional examples and documentation

---

## License

Fusion is released under the MIT License.
