# Name of the CLI binary
CLI_NAME=todo-cli

install:
	@echo "Installing dependencies..."
	go mod tidy

local-run:
	@echo "Running the project locally on 8080 port ..."
	go run main.go serve

build:
	@echo "Building the project and creating binary $(CLI_NAME)..."
	go build -o $(CLI_NAME)

run:
	@echo "Running the binary $(CLI_NAME) with arguments: $(ARGS)..."
	./$(CLI_NAME) $(ARGS)

build-and-run:
	@echo "Building and running the project..."
	./build-and-run.sh $(ARGS)
