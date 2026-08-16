

# ============================================================================
# Загрузка переменных окружения
# ============================================================================
ifneq (,$(wildcard .env))
    include .env
    export
endif



# ============================================================================
# Сборка строки подключения
# ============================================================================
CONNECTION_STRING = postgresql://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=disable
MIGRATIONS_PATH ?= ./db/migrations



# ============================================================================
# Значения по умолчанию, если .env не найден
# ============================================================================
POSTGRES_USER ?= waves2k
POSTGRES_PASSWORD ?= rootpassword
POSTGRES_DB ?= go_wallet
POSTGRES_HOST ?= localhost
POSTGRES_PORT ?= 5438



# ============================================================================
# Builds the application
# Usage: make build
# ============================================================================
build:
	@go build -o ./bin/app ./cmd/main.go



# ============================================================================
# Builds and runs the application
# Usage: make run
# ============================================================================
run: build
	@./bin/app



# ============================================================================
# Runs all the tests
# Usage: make test
# ============================================================================
test:
	@go test -v ./...	



# ============================================================================
# Stop services from docker-compose.yml file
# Usage: make docker-compose-down
# ============================================================================
docker-compose-down:
	@docker-compose down -v



# ============================================================================
# Run services from docker-compose.yml file
# Usage: make docker-compose-up
# ============================================================================
docker-compose-up:
	@docker-compose up -d



# ============================================================================
# Migration command
# Usage: make migrate-create NAME=create_users_table
# ============================================================================
migrate-create:
	@migrate create -ext sql -dir ${MIGRATIONS_PATH} -seq ${NAME}
	@echo "✅ Migration created successfully"



# ============================================================================
# Commits N migraitons
# Usage: make migrate-up STEPS=n
# ============================================================================
migrate-up:
	@echo "🔄 Applying migrations..."
	@migrate -path ${MIGRATIONS_PATH} -database "${CONNECTION_STRING}" up ${STEPS}
	@echo "✅ Migrations applied successfully"



# ============================================================================
# Rollbacks N migraitons
# Usage: make migrate-down STEPS=n
# ============================================================================
migrate-down:
	@echo "🔄 Applying migrations..."
	@migrate -path ${MIGRATIONS_PATH} -database "${CONNECTION_STRING}" down ${STEPS}
	@echo "✅ Migrations applied successfully"



# ============================================================================
# Gets migration current version
# Usage: make migrate-version
# ============================================================================
migrate-version:
	@migrate -path ${MIGRATIONS_PATH} -database "${CONNECTION_STRING}" version



# ============================================================================
# Force to certain migration version
# Usage: make migrate-force VERSION=N
# ============================================================================
migrate-force:
	@echo "⚠️  WARNING: force will change version without running migrations"
	@read -p "Are you sure you want to force version $(VERSION)? (y/N) " confirm; \
	if [ "$$confirm" != "y" ]; then exit 1; fi
	@migrate -path ${MIGRATIONS_PATH} -database "${CONNECTION_STRING}" force $(VERSION)
