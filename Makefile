.DEFAULT_GOAL := help

DOCKER_TAG := latest

.PHONY: build
build: ## Build go module
	go build -o ./bin/go_skel .

.PHONY: build-image
build-image: ## Build image for deploy
	docker build -t kijimad/go_skel:${DOCKER_TAG} \
	--target release ./

.PHONY: build-local
build-local: ## Build image for local development
	docker-compose build --no-cache

.PHONY: up
up: ## Do docker compose up
	docker-compose up -d

.PHONY: down
down: ## Do docker compose down
	docker-compose down

.PHONY: logs
logs: ## Tail docker compose logs
	docker-compose logs -f

.PHONY: ps
ps: ## Check container status
	docker-compose ps

.PHONY: lint
lint: ## Run lint
	docker run --rm -v ${PWD}:/app -w /app golangci/golangci-lint:v1.51.2 golangci-lint run -v
	docker run --rm -v ${PWD}:/app -w /app golang:1.20 go vet ./...

.PHONY: run
run: ## Run
	go run .

.PHONY: test
test: ## Run test
	go test -race -shuffle=on -v ./...

.PHONY: help
help: ## Show help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
