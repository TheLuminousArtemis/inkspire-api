include .envrc
MIGRATIONS_PATH = ./migrate/migrations

.PHONY: migrate create
migration:
	@migrate create -seq -ext sql -dir $(MIGRATIONS_PATH) $(filter-out $@, $(MAKECMDGOALS))

.PHONY: migrate up
migrate-up:
	@migrate -path=$(MIGRATIONS_PATH) -database=$(DB_ADDR) up

.PHONY: migrate down
migrate-down:
	@migrate -path=$(MIGRATIONS_PATH) -database=$(DB_ADDR) down $(filter-out $@, $(MAKECMDGOALS))

.PHONY: seed
seed:
	@go run ./migrate/seed/main.go

.PHONY: gen-docs
gen-docs:
	# @swag init -g ./api/main.go -d ./cmd/api,./internal/db,./internal/store, && swag fmt
	@swag init -g ./main.go -d cmd/api,./internal/db,./internal/store && swag fmt

.PHONY: migrate-force 
migrate-force:   
	@migrate -path=$(MIGRATIONS_PATH) -database=$(DB_ADDR) force $(filter-out $@,$(MAKECMDGOALS))

.PHONY: test
test:
	@go test -v ./...