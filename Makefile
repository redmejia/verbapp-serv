down: clean
	@echo "Clean and down serivices..."
	docker compose down

start: clean build
	@echo "Starting..."
	docker compose up --build 

up_build: build
	@echo "Stop docker images..."
	docker compose down
	@echo "Building and Starting docker images ..."
	docker compose up --build 

build:
	@echo "Building..."
	@GOOS=linux  CGO_ENABLED=0 go build -o cmd/api/dist/chat_service cmd/api/main.go


clean:
	@echo "Clean..."
	@rm -f cmd/api/dist/chat_service