.PHONY: build run clean test

APP_NAME=hotel-booking
BUILD_DIR=./build

build:
	@echo "Build app"
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(APP_NAME) ./cmd/app

run:
	@echo "Run app"
	@go run ./cmd/app

clean:
	@echo "Clean"
	@rm -rf $(BUILD_DIR)

test:
	@echo "Test"
	@go test -v ./...
