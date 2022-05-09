NAME=mailchimp-importer
VERSION=0.1.0

all: docker
build: compile
test: unit

run:
	go run cmd/main.go

compile:
	go build -o /dev/null -ldflags "-s -w" ./cmd
	# compiling without binary output (remove -o /dev/null if you want to generate a binary)

deps:
	go mod download
	@go mod tidy

unit:
	export GO_TEST_RUN_INTEGRATION_TESTS=false && \
	go test -coverprofile=cover.out -timeout 300s ./...

docker:
	docker build -f build/Dockerfile . -t $(NAME):$(VERSION)

docker-run:
	sudo docker run --name $(NAME) $(NAME):$(VERSION)

.PHONY: all build test run deps unit docker
