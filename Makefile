POSTGRESQL_URL = $(shell echo %DB_CONN_URL%)

server:
	@go run ./cmd/http/.

exp:
	@go run ./cmd/exp/exp.go

db:
	@docker compose up -d	

migrate-up:
	@migrate -database ${POSTGRESQL_URL} -path db/migrations up

migrate-down:
	@migrate -database ${POSTGRESQL_URL} -path db/migrations down

.PHONY: server exp migrate-up migrate-down	