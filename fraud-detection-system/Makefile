.PHONY: help up down restart logs build clean test migrate seed

# Default target
help:
	@echo "Banking Fraud Detection System - Available Commands:"
	@echo ""
	@echo "  make up              - Start all services"
	@echo "  make down            - Stop all services"
	@echo "  make restart         - Restart all services"
	@echo "  make logs            - View logs from all services"
	@echo "  make build           - Build all Docker images"
	@echo "  make clean           - Clean up containers and volumes"
	@echo ""
	@echo "  make migrate-up      - Run database migrations"
	@echo "  make migrate-down    - Rollback database migrations"
	@echo "  make seed            - Seed database with initial data"
	@echo ""
	@echo "  make test            - Run all tests"
	@echo "  make test-backend    - Run backend tests"
	@echo "  make test-ml         - Run ML service tests"
	@echo ""
	@echo "  make train-models    - Train all ML models"
	@echo "  make lint            - Run linters"
	@echo "  make format          - Format code"

# Docker Compose commands
up:
	docker-compose up -d
	@echo "Services started. API Gateway: http://localhost:8080, ML Service: http://localhost:8000"

down:
	docker-compose down

restart:
	docker-compose restart

logs:
	docker-compose logs -f

logs-api:
	docker-compose logs -f api-gateway

logs-ml:
	docker-compose logs -f ml-service

logs-worker:
	docker-compose logs -f worker

build:
	docker-compose build

clean:
	docker-compose down -v
	rm -rf backend/bin/
	find . -type d -name "__pycache__" -exec rm -rf {} +
	find . -type f -name "*.pyc" -delete

# Database commands
migrate-up:
	docker-compose exec api-gateway /app/migrate -path=/app/migrations -database="postgres://frauddetection:frauddetection_password@postgres:5432/frauddetection_db?sslmode=disable" up

migrate-down:
	docker-compose exec api-gateway /app/migrate -path=/app/migrations -database="postgres://frauddetection:frauddetection_password@postgres:5432/frauddetection_db?sslmode=disable" down 1

migrate-create:
	@read -p "Enter migration name: " name; \
	migrate create -ext sql -dir backend/internal/database/migrations -seq $$name

seed:
	docker-compose exec postgres psql -U frauddetection -d frauddetection_db -f /docker-entrypoint-initdb.d/seeds/rbi_circulars.sql
	docker-compose exec postgres psql -U frauddetection -d frauddetection_db -f /docker-entrypoint-initdb.d/seeds/sender_registry.sql

# Testing
test: test-backend test-ml

test-backend:
	cd backend && go test ./... -v -cover

test-ml:
	cd ml-service && python -m pytest tests/ -v

test-integration:
	cd tests && go test ./integration/... -v

# ML Model training
train-models:
	docker-compose exec ml-service python scripts/train_all_models.py

train-distilbert:
	docker-compose exec ml-service python training/train_distilbert.py

train-roberta:
	docker-compose exec ml-service python training/train_roberta.py

train-lstm:
	docker-compose exec ml-service python training/train_lstm.py

train-xgboost:
	docker-compose exec ml-service python training/train_xgboost.py

# Code quality
lint:
	cd backend && golangci-lint run
	cd ml-service && flake8 app/ training/

format:
	cd backend && gofmt -s -w .
	cd ml-service && black app/ training/

# Development
dev-backend:
	cd backend && go run cmd/api-gateway/main.go

dev-ml:
	cd ml-service && uvicorn app.main:app --reload --port 8000

dev-worker:
	cd backend && go run cmd/worker/main.go

# Install dependencies
install-backend:
	cd backend && go mod download

install-ml:
	cd ml-service && pip install -r requirements.txt

# Database backup
backup-db:
	./scripts/backup.sh

# Setup
setup:
	cp .env.example .env
	@echo "Please edit .env file with your configuration"
	@echo "Then run: make up && make migrate-up && make seed"

