# Run tests
make test:
    go test ./...

# Run linter
make lint:
    golangci-lint run

# Generate mock interfaces
make mock:
    mockgen -source=internal/repository/interfaces.go -destination=internal/repository/mocks/repository_mocks.go

# Run migrations
make migrate:
    migrate -path migrations -database "postgresql://postgres:postgres@localhost:5432/ecommerce?sslmode=disable" up