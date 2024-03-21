GO_BUILD_PATH ?= bin
GO_BUILD_WAREHOUSES_PATH ?= $(GO_BUILD_PATH)/warehouses/

PATH := $(PATH):/home/linuxbrew/.linuxbrew/opt/go/libexec/bin
export PATH



.PHONY: build
build:
	CGO_ENABLED=$(CGO) GOOS=linux go build -o $(GO_BUILD_WAREHOUSES_PATH) ./cmd/warehouses/

.PHONY: up
up:
	docker compose up --build -d

.PHONY: down
down:
	docker compose down

.PHONY: test
make test:
	go test ./...

.PHONY: clean
clean:
	@rm -rf ./bin