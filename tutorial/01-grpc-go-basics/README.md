# gRPC Go Basics – Coffee Shop Example

This tutorial is based on:

"2024 gRPC Golang Tutorial – The tutorial I wish I had when I was learning"  
by Code With Cypert  
https://www.youtube.com/watch?v=mPESsBfUKkc

The goal of this tutorial is to understand:

- gRPC service definitions with Protocol Buffers
- Unary RPC
- Server streaming RPC
- Code generation with protoc
- Client-server communication in Go
- Makefile-based automation

---

## Project Structure

```
01-grpc-go-basics/
├── proto/
│   └── coffee_shop.proto
├── server/
│   └── server.go
├── client/
│   └── client.go
├── Makefile
```

---

## Generate Protobuf Code

```bash
make proto
```

## Run Server

```bash
make run-server
```

## Run Client

```bash
make run-client
```

## Run Both

```bash
make run
```

---

## Learning Focus

This lab explores:

- Streaming responses (`GetMenu`)
- Unary RPC (`PlaceOrder`)
- Proto-to-Go type mapping
- Basic concurrency with goroutines
- Project structure for service-based systems

---

Part of my distributed systems & gRPC learning series.