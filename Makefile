
.PHONY: run

# Start the server
run:
	@go run main.go


# Build the server
build:
	@go build -o bin/scheduler-extender main.go


clear:
	@rm -rf bin


