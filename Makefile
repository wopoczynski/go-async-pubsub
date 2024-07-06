.PHONY: help
help:
	@grep -E '^[0-9a-z.A-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: create-local-env-file
create-local-env-file:
	cp .env.example .env

.PHONY: build
build:
	@mkdir -p build
	go build -o build/todo \
		-ldflags=" \
			-X 'app/version.Version=dev' \
			-X 'app/version.CommitHash=xxxx' \
			-X 'app/version.BuildTimestamp=$$(date -Iseconds)' \
		" \
		./cmd

.PHONY: build-docker-image
build-docker-image:
	docker build \
		--build-arg "APP_BUILD_TIME=$$(date -Iseconds)" \

		--tag todoapp \
		--file ./docker/Dockerfile .

.PHONY: install-format-tools
install-format-tools:
	go install mvdan.cc/gofumpt@latest
	go install github.com/daixiang0/gci@latest

.PHONY: format
format:
	git ls-files | grep -E ".go$$" | \
		xargs gci write --skip-generated \
			--skip-vendor -s standard -s default -s blank -s dot
	gofumpt -w ./..

.PHONY: test
test:
	go test -count=1 -v -race $$(go list ./... | grep -v /vendor/ | grep -v /test/)
	go test -count=1 -v -race $$(go list ./test/...)

.PHONY: install-swag
install-swag:
	go install github.com/swaggo/swag/cmd/swag@v1.16.3

.PHONY: generate-docs
generate-docs:
	${HOME}/go/bin/swag init --parseDependency --parseInternal --parseDepth 4 -g internal/server/server.go

.PHONY: local-dev-env-up
local-up:
	docker-compose -f ./env/local/compose.yaml up -d --build

.PHONY: local-dev-env-down
local-down:
	docker-compose -f ./env/local/compose.yaml down

.PHONY: run
run:
	go run \
	-ldflags=" \
		-X 'app/version.Version=dev' \
		-X 'app/version.CommitHash=xxxx' \
		-X 'app/version.BuildTimestamp=$$(date -Iseconds)' \
	" \
	./cmd
