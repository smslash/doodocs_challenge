# Docker
DOCKER_IMAGE_NAME := doodocs-app

build-docker:
	@docker build -t $(DOCKER_IMAGE_NAME) .

run-docker:
	@docker run -p 8080:8080 --env-file .env $(DOCKER_IMAGE_NAME)

build-and-run: build-docker run-docker

# Тест
test:
	go test ./tests/createArchive_test.go
	go test ./tests/getArchiveInfo_test.go

# CMD
run:
	go get github.com/joho/godotenv
	go get github.com/stretchr/testify/assert
	go run ./cmd/server/main.go