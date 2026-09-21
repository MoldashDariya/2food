# 2food API Documentation

This directory contains the auto-generated API documentation for the 2food backend service, powered by Swaggo and Scalar.

⚠️ **WARNING: Do not edit `swagger.json` or `swagger.yaml` manually!**
These files are automatically generated based on the Go source code annotations. Any manual changes will be lost upon the next generation.

## 🛠 Prerequisites

To generate the documentation, you need the `swag` command-line tool. If you don't have it installed, run:

```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

*(Note: Ensure your `$(go env GOPATH)/bin` directory is added to your system's `$PATH`).*

## 🔄 Generating the Docs

Whenever you add a new endpoint, change a router path, or modify request/response structs, you must regenerate the documentation.

**Run this command from the root of the project:**

```bash
swag init -g cmd/2food/main.go
```

## 🚀 Viewing the API Reference

1. Start the backend server:
  go run cmd/2food/main.go
2. Open your web browser and navigate to:
  [http://localhost:8080/docs](http://localhost:8080/docs) (or your actual server domain/port)
