local-db:
	docker-compose -f development/postgres_local/docker-compose.yaml up -d

local-init:
	go run ./cmd/bun db init

local-migrate:
	go run ./cmd/bun db migrate
