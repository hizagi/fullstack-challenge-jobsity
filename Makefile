up:
	@echo "Starting the frontend, backend, and MongoDB services..."
	docker-compose up --build

down:
	@echo "Stopping all services..."
	docker-compose down