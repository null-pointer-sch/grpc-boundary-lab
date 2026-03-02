PROTO_DIR=proto
PROTO_FILE=$(PROTO_DIR)/coffee_shop.proto

SERVER_DIR=server
CLIENT_DIR=client

SERVER_BIN=bin/server
CLIENT_BIN=bin/client

.PHONY: proto server client run-server run-client clean

proto:
	protoc \
	  --go_out=. --go_opt=paths=source_relative \
	  --go-grpc_out=. --go-grpc_opt=paths=source_relative \
	  $(PROTO_FILE)

server: proto
	mkdir -p bin
	go build -o $(SERVER_BIN) ./$(SERVER_DIR)

client: proto
	mkdir -p bin
	go build -o $(CLIENT_BIN) ./$(CLIENT_DIR)

run-server: server
	./$(SERVER_BIN)

run-client: client
	./$(CLIENT_BIN)

clean:
	rm -rf bin

.PHONY: run

run: proto
	@echo "Starting server..."
	@go run ./server & \
	SERVER_PID=$$!; \
	sleep 1; \
	echo "Starting client..."; \
	go run ./client; \
	echo "Stopping server..."; \
	kill $$SERVER_PID