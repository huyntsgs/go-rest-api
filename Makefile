BUILD_INFO_FLAGS := -X main.BuildTime=$(shell date -u '+%Y-%m-%d_%H:%M:%S') -X main.CommitHash=$(shell git rev-parse HEAD)
BINARY_NAME := gorestapi
IMAGE_NAME := goapi

.PHONY: build 
build:
	go build -ldflags "$(BUILD_INFO_FLAGS)" -a -o $(BINARY_NAME) .

.PHONY: build-linux
build-linux:
	env GOOS=linux GOARCH=amd64 GO11MODULE=ON go build -ldflags "$(BUILD_INFO_FLAGS)" -a -o $(BINARY_NAME) .
	
.PHONY: build-docker
build-docker:
	env GOOS=linux GOARCH=amd64 GO11MODULE=ON go build -ldflags "$(BUILD_INFO_FLAGS)" -a -o $(BINARY_NAME) .
	docker build -t 	$(IMAGE_NAME) .
.PHONY: clean
clean:
	rm -f $(BINARY_NAME)	
