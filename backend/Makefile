GENERATED_DIR := api/generated
OPENAPI_SPEC := api/openapi.yaml

run:
	go run main.go

lint:
	golangci-lint run

build:
	go build -o myapp

docker-build:
	docker build -t myapp .

docker-run:
	docker run -p 8080:8080 myapp

clean-spec:
	@echo "Cleaning generated files..."
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

.PHONY: run lint build docker-build docker-run clean-spec install generate install-and-generate clean-install-generate up down logs clean-container
