all: start
build:
	go build -o ./cmd/authority/authority ./cmd/authority
	go build -o ./cmd/frontend/frontend ./cmd/frontend
	go build -o ./cmd/treasury/treasury ./cmd/treasury
	docker-compose build
start: build
	docker-compose up
