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