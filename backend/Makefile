GENERATED_DIR := api/generated
OPENAPI_SPEC := api/openapi.yaml

run:
	go run cmd/server/main.go

lint:
	golangci-lint run

build:
	go build -o taskapp

docker-build:
	docker build -t taskapp .

docker-run:
	docker run -p 8080:8080 taskapp

clean-spec:
	rm -rf $(GENERATED_DIR)

generate:
	@echo "Generating server and types..."
	mkdir -p $(GENERATED_DIR)
	oapi-codegen -generate chi-server -o $(GENERATED_DIR)/server.gen.go -package generated $(OPENAPI_SPEC)
	oapi-codegen -generate types -o $(GENERATED_DIR)/types.gen.go -package generated $(OPENAPI_SPEC)

install-and-generate: install generate

clean-install-generate: clean-spec install-and-generate

up:
	docker-compose up -d

down:
	docker-compose down

logs:
	@echo "Displaying logs for Docker Compose services..."
	docker-compose logs -f

clean-container:
	@echo "Removing containers, and volumes..."
	docker-compose down -v

unit-test:
	go fmt ./...
	go test -vet all ./...

integration-test:
	go test -v -vet all -tags=integration ./...

race-test:
	go test -v -race ./...

test: unit-test integration-test race-test

get-mocks-generator:
	go install -mod=mod go.uber.org/mock/mockgen

generate-mocks: get-mocks-generator
	go generate -v ./...

clean-mocks:
	find . -type f -name '*_mocks.go' -exec rm -f {} \;

.PHONY: run lint build docker-build docker-run clean-spec install generate install-and-generate clean-install-generate up down logs clean-container unit-test integration-test race-test test get-mocks-generator generate-mocks clean-mocks
