down: clean
	@echo "clean and down serivices..."
	docker-compose down

up_build: build
	@echo "Stop docker images..."
	docker-compose down
	@echo "Building and Starting docker images ..."
	docker-compose up 

build:
	@echo "building ..."
	@GOOS=linux  CGO_ENABLED=0 go build -o cmd/api/dist/chat_service cmd/api/main.go


clean:
	@echo "clean..."
	@ rm cmd/api/dist/chat_service