APP_NAME=go-test-case-api

.PHONY: run
run:
	GOFLAGS=-buildvcs=false docker compose up --build


.PHONY: new-migration
new-migration:
ifndef name
	$(error "you need to input the migration name. ex: make new-migration name=create_accounts_table")
endif
	migrate create -ext sql -dir migrations -seq $(name)

# =========================
# TESTS
# =========================

.PHONY: test
test:
	- CGO_ENABLED=0 go test -v -short ./...


.PHONY: test-functional
test-functional:
	@echo "Starting test Postgres container..."
	docker-compose -f docker-compose.test.yml up -d
	@sleep 5
	@echo "Setting environment variables..."
	export DB_HOST=localhost
	export DB_PORT=5433
	export DB_USER=test
	export DB_PASSWORD=test
	export DB_NAME=testdb
	@echo "Running functional tests..."
	- CGO_ENABLED=0 go test -v -tags=functional ./internal/functional_tests || true
	@echo "Stopping test Postgres container..."
	docker-compose -f docker-compose.test.yml down

# =========================
# COVERAGE
# =========================

.PHONY: coverage
coverage:
	go test ./internal/handlers ./internal/service -coverprofile=coverage.out
	go tool cover -func=coverage.out
	grep -v "/mocks/" coverage.out > coverage.tmp
	go tool cover -html=coverage.out


# =========================
# UTIL
# =========================

.PHONY: tidy
tidy:
	go mod tidy


.PHONY: clean
clean:
	rm -f coverage*.out