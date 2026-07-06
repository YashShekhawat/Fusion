# Fusion Architecture

> This document explains how Fusion is structured, why it is designed this way, and how to contribute without breaking the architecture.

---

# What is Fusion?

Fusion is a **provider-agnostic AI SDK for Go**.

Instead of learning and integrating multiple AI SDKs (OpenAI, Gemini, Anthropic, etc.), developers only interact with Fusion.

Fusion handles the provider-specific implementation behind the scenes.

For example, your application only needs to write:

```go
response, err := client.Generate(ctx, request)
```

Whether that request is sent to OpenAI, Gemini, Anthropic, or another provider is handled internally by Fusion.

---

# Architecture Overview

Every request flows through the SDK like this:

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

Each layer has **one responsibility**.

Keeping responsibilities separate makes the SDK easier to maintain, extend, and test.

---

# Project Structure

```text
fusion/
├── client/
├── drivers/
│   ├── openai/
│   ├── gemini/
│   └── ...
├── models/
├── fusionerrors/
├── docs/
└── ...
```

---

# Package Responsibilities

## models/

### Purpose

Contains the common request and response structures shared across the SDK.

Example:

```go
type GenerateRequest struct {
    ...
}

type GenerateResponse struct {
    ...
}
```

### Why?

Every provider has different request and response formats.

Fusion exposes one common model while each driver converts it into the provider's native format.

---

## fusionerrors/

### Purpose

Contains standard SDK errors.

Examples:

```go
ErrUnauthorized
ErrTimeout
ErrRateLimit
ErrProviderUnavailable
ErrInvalidResponse
```

### Why?

Every provider returns different error messages.

Fusion converts them into common errors so applications don't need provider-specific error handling.

Instead of:

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

---

## client/

### Purpose

The client is the public entry point of the SDK.

Applications should interact with the client instead of individual providers.

Example:

```go
client.Generate(ctx, request)
```

### Responsibilities

* Accept requests from the application
* Select the correct provider
* Call the driver
* Return a common response

### The client should NOT

* Build provider requests
* Send HTTP requests
* Parse provider responses
* Contain provider-specific logic

Those responsibilities belong to the drivers.

---

## registry/

### Purpose

The registry stores all available providers.

Internally it behaves like:

```go
map[string]Driver
```

Example:

```text
openai    -> OpenAIDriver
gemini    -> GeminiDriver
anthropic -> AnthropicDriver
```

### Why?

Without a registry the client would need code like:

```go
if provider == "openai" {
    ...
}

if provider == "gemini" {
    ...
}
```

Every new provider would require modifying the client.

With the registry:

```text
Client

↓

Registry

↓

Driver
```

Adding a new provider only requires registering it.

The client never changes.

---

## drivers/

### Purpose

Each driver knows how to communicate with one AI provider.

Example:

```text
drivers/
    openai/
    gemini/
    anthropic/
```

Each driver is independent.

Adding one provider should never require modifying another provider.

---

# Driver Structure

Each provider follows the same layout.

```text
drivers/
    openai/
        driver.go
        request.go
        response.go
        transport.go
        errors.go
```

---

## driver.go

The driver coordinates the request.

Responsibilities:

* Validate input
* Convert Fusion request
* Send request
* Convert provider response

The driver should **not** contain HTTP implementation details.

---

## request.go

Responsible for converting Fusion models into provider-specific models.

```text
Fusion Request

↓

OpenAI Request
```

Keeping mapping here means request format changes only affect one file.

---

## response.go

Responsible for converting provider responses into Fusion responses.

```text
OpenAI Response

↓

Fusion Response
```

Applications never need to understand provider-specific JSON.

---

## transport.go

Responsible only for communication.

Responsibilities:

* Build HTTP requests
* Send requests
* Receive responses
* Decode JSON

It should not contain business logic.

---

## errors.go

Responsible for converting provider-specific errors into Fusion errors.

Example:

```text
401

↓

ErrUnauthorized
```

Keeping this logic separate keeps transport clean and focused.

---

# Design Principles

Fusion follows several software engineering principles.

## Single Responsibility Principle

Every package should have one responsibility.

Instead of one large file doing everything:

```text
Generate()

↓

Create JSON

↓

Send HTTP

↓

Decode Response

↓

Handle Errors

↓

Retry

↓

Logging
```

the work is divided into focused components.

---

## Separation of Concerns

Each layer solves one problem.

| Component | Responsibility            |
| --------- | ------------------------- |
| Client    | Public API                |
| Registry  | Find the correct provider |
| Driver    | Coordinate provider logic |
| Request   | Build provider request    |
| Transport | HTTP communication        |
| Response  | Parse provider response   |
| Errors    | Map provider errors       |

---

## Extensibility

Adding a new provider should require as little work as possible.

Ideally:

1. Create a new driver.
2. Register it.
3. Done.

Existing providers should never need modification.

---

# Contribution Guidelines

When contributing to Fusion, please follow these rules.

## ✅ Do

* Keep packages focused on one responsibility.
* Follow the existing driver structure.
* Reuse common models and errors where appropriate.
* Keep provider-specific logic inside the provider package.
* Write clear comments when adding new architecture.

---

## ❌ Don't

Do not make one package responsible for multiple concerns.

Examples:

* Don't add HTTP code inside `client`.
* Don't build requests inside `transport`.
* Don't parse responses inside `driver`.
* Don't access one provider from another provider.
* Don't bypass the registry.

---

# Adding a New Provider

Every provider should follow this structure.

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
2. Map Fusion requests to provider requests.
3. Send requests using transport.
4. Map provider responses.
5. Convert provider errors into Fusion errors.
6. Register the driver.

No other package should require modification.

---

# Future Architecture

Fusion is designed to grow without changing its core architecture.

Planned features include:

* Middleware
* Logging
* Retry support
* Streaming
* Tool calling
* Additional providers

Whenever possible, new features should be provider-agnostic.

---

# Guiding Principle

Before adding new code, ask yourself:

> **Does this belong to Fusion, or does it belong to a specific provider?**

If every provider could use it, it probably belongs in Fusion.

If it only applies to one provider, keep it inside that provider's package.

Following this principle helps keep Fusion modular, maintainable, and easy to extend.
