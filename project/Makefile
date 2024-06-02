FRONTEND_BINARY=frontendApp
ANALYZER_BINARY=analyzerApp

## up: starts all containers in the background without forcing build
up:
	@echo "Starting Docker images..."
	docker-compose up -d
	@echo "Docker images started!"

## up_build: stops docker-compose (if running), builds all projects and starts docker compose
up_build: build_analyzer
	@echo "Stopping docker images (if running...)"
	docker-compose down
	@echo "Building (when required) and starting docker images..."
	docker-compose up --build -d
	@echo "Docker images built and started!"

## down: stop docker compose
down:
	@echo "Stopping docker compose..."
	docker-compose down
	@echo "Done!"

## build_analyzer: builds the analyzer binary as a linux executable
build_analyzer:
	@echo "Building analyzer binary..."
	cd ../webpage-analyzer-service && env GOOS=linux CGO_ENABLED=0 go build -o ${ANALYZER_BINARY} ./cmd/api
	@echo "Done!"

## build_frontend: builds the frontend binary
build_frontend:
	@echo "Building frontend binary..."
	cd ../frontend && env CGO_ENABLED=0 go build -o ${FRONTEND_BINARY} ./cmd/web
	@echo "Done!"

## start: starts the frontend
start: build_frontend
	@echo "Starting frontend"
	cd ../frontend && ./${FRONTEND_BINARY} &

## stop: stop the frontend
stop:
	@echo "Stopping frontend..."
	@-pkill -SIGTERM -f "./${FRONTEND_BINARY}"
	@echo "Stopped frontend!"