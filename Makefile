all: up

up:
	docker-compose up --build

lint:
	golangci-lint run
	golint ./...

test:
	docker build -f Dockerfile.test .
