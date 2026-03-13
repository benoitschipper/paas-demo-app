BINARY_NAME := paas-demo-app
IMAGE_NAME  ?= paas-demo-app
IMAGE_TAG   ?= latest
REGISTRY    ?= quay.io/your-org

.PHONY: build run test docker-build docker-push clean

build:
	CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w -X main.version=$(IMAGE_TAG)" -o bin/$(BINARY_NAME) ./cmd/server

run:
	go run ./cmd/server

test:
	go test ./... -v -race -count=1

docker-build:
	docker build -f Containerfile -t $(REGISTRY)/$(IMAGE_NAME):$(IMAGE_TAG) .

docker-push:
	docker push $(REGISTRY)/$(IMAGE_NAME):$(IMAGE_TAG)

clean:
	rm -rf bin/
