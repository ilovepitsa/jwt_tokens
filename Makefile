.PHONY: up
up:
	docker compose up --build -d

.PHONY: down
down:
	docker compose down --remove-orphans

.PHONY: restart
restart:
	docker compose down --remove-orphans
	docker compose up --build -d

.PHONY: test
test:
	CGO_ENABLE=1 go test -v -race ./...

.PHONY: cover-html
cover-html: 
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o cover.html
	open cover.html
	rm coverage.out

.PHONY: run
run:
	go run ./cmd