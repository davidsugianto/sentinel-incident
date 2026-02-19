DB_URL=postgres://user:pass@localhost:5432/sentinel_incident?sslmode=disable
MIGRATE_CMD=docker run --rm -v $(PWD)/migrations:/migrations --network host migrate/migrate -path=/migrations -database "$(DB_URL)"

# Build Go app
build:
	@echo "🚀 Building app..."
	go build -o bin/app ./cmd/http

# Run with Air (hot reload)
dev:
	@echo "🔄 Starting development with Air..."
	docker compose -f docker-compose.yml up --build

# Run migrations
migrate-up:
	@echo "⬆️  Running all migrations..."
	$(MIGRATE_CMD) up

migrate-down:
	@echo "⬇️  Rolling back last migration..."
	$(MIGRATE_CMD) down 1

migrate-drop:
	@echo "💣 Dropping all migrations..."
	$(MIGRATE_CMD) drop -f

migrate-force:
	@echo "⚡ Forcing migration version..."
	$(MIGRATE_CMD) force 1

# Clean containers
clean:
	@echo "🧹 Cleaning containers..."
	docker compose -f docker-compose.dev.yml down -v