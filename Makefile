build:
	@go build -o bin/help cmd/main.go

run: build
	@./bin/help

migration:
	@migrate create -ext sql -dir cmd/migrate/migrations $(filter-out $@,$(MAKECMDGOALS))

migrate-up:
	@go run cmd/migrate/main.go up

migrate-down:
	@go run cmd/migrate/main.go down 

.PHONY: test
test:
	@if [ -n "$(TEST)" ]; then \
		echo "Running single test: $(TEST)"; \
		go test ./... -run $(TEST); \
	else \
		echo "Running all tests..."; \
		go test ./...; \
	fi
