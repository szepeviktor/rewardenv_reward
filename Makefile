.DEFAULT_GOAL 	= help

SHELL         	= bash
project       	= reward
GIT_AUTHOR    	= janosmiko

help: ## Outputs this help screen
	@grep -E '(^[\/a-zA-Z0-9_-]+:.*?##.*$$)|(^##)' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}{printf "\033[32m%-30s\033[0m %s\n", $$1, $$2}' | sed -e 's/\[32m##/[33m/'

# If the first argument is "gen"...
ifeq (gen,$(firstword $(MAKECMDGOALS)))
  # use the rest as arguments for "run"
  GEN_ARGS := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))
  # ...and turn them into do-nothing targets
  $(eval $(GEN_ARGS):;@:)
endif

## —— Commands —————————————————————————————————————————————————————————
build: ## Build the command to ./dist
	go mod download
	go generate ./...
	go install github.com/go-bindata/go-bindata/...@latest
	go-bindata -pkg internal -o internal/bindata.go templates/... VERSION.txt
	CGO_ENABLED=0 go build -ldflags="-s -w" -o dist/reward ./main.go

build-all: ## Build the binaries using goreleaser (without releasing it)
	goreleaser release --rm-dist --auto-snapshot --skip-publish

build-all: ## Build the binaries using goreleaser (without releasing it)
	goreleaser release --rm-dist --auto-snapshot --skip-publish

## —— Go Commands —————————————————————————————————————————————————————————
gomod: ## Update Go Dependencies
	go mod tidy

lint: ## Lint Go Code
	golangci-lint run ./...

test: ## Run Go tests
	go test -race ./internal/crypto ./internal/shell ./internal/docker ./internal/dockercompose ./internal/util -v
