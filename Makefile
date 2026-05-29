.PHONY: all frontend build build-linux clean

BINARY_NAME=internal-navigation
DIST_DIR=dist
FRONTEND_DIR=frontend

all: build

frontend:
	@echo "Building frontend..."
	cd $(FRONTEND_DIR) && npm install && npm run build

build: frontend
	@echo "Building $(BINARY_NAME)..."
	go build -buildvcs=false -o $(BINARY_NAME) .

build-linux: frontend
	@echo "Cross compiling for Linux amd64..."
	GOOS=linux GOARCH=amd64 CGO_ENABLED=1 CC=x86_64-linux-musl-gcc CXX=x86_64-linux-musl-g++ go build -buildvcs=false -ldflags '-linkmode external -extldflags "-static"' -o $(BINARY_NAME)-linux-amd64 .
	@echo "Binary: $(BINARY_NAME)-linux-amd64"

build-linux-simple: frontend
	@echo "Cross compiling for Linux amd64 (CGO disabled)..."
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -buildvcs=false -o $(BINARY_NAME)-linux-amd64 .
	@echo "Binary: $(BINARY_NAME)-linux-amd64"

clean:
	rm -f $(BINARY_NAME) $(BINARY_NAME)-linux-amd64
	rm -rf $(DIST_DIR)
