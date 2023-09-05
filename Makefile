
.PHONY: run build docker push clear

# Start the server
run:
	@go run main.go


# Build the server
build:
	@go build -o bin/scheduler-extender main.go


# Build the docker image
docker:
	@docker build -t zhengyscn/scheduler-extender:v1.0.6 .


# Push the docker image
push:
	@docker push zhengyscn/scheduler-extender:v1.0.6

# Clear the bin directory
clear:
	@rm -rf bin