#!/bin/bash

# Banking Fraud Detection System - Setup Script

set -e

echo "🚀 Setting up Banking Fraud Detection System..."

# Check if Docker is installed
if ! command -v docker &> /dev/null; then
    echo "❌ Docker is not installed. Please install Docker first."
    exit 1
fi

# Check if Docker Compose is installed
if ! command -v docker-compose &> /dev/null; then
    echo "❌ Docker Compose is not installed. Please install Docker Compose first."
    exit 1
fi

# Create .env file if it doesn't exist
if [ ! -f .env ]; then
    echo "📝 Creating .env file from .env.example..."
    cp .env.example .env
    echo "⚠️  Please edit .env file with your configuration before proceeding."
    echo "   Especially set JWT_SECRET to a secure random string."
    read -p "Press enter to continue after editing .env..."
fi

# Create necessary directories
echo "📁 Creating necessary directories..."
mkdir -p models/{distilbert,roberta,lstm,xgboost,ensemble}
mkdir -p ml-service/data/{raw,processed,external}
mkdir -p database/backups
mkdir -p logs

# Create .gitkeep files
touch ml-service/data/raw/.gitkeep
touch ml-service/data/processed/.gitkeep
touch ml-service/data/external/.gitkeep
touch database/backups/.gitkeep

# Build Docker images
echo "🐳 Building Docker images..."
docker-compose build

# Start services
echo "🚀 Starting services..."
docker-compose up -d

# Wait for services to be ready
echo "⏳ Waiting for services to be ready..."
sleep 10

# Check if PostgreSQL is ready
echo "🔍 Checking PostgreSQL..."
until docker-compose exec -T postgres pg_isready -U frauddetection; do
    echo "Waiting for PostgreSQL..."
    sleep 2
done

# Run database migrations
echo "📊 Running database migrations..."
docker-compose exec -T postgres psql -U frauddetection -d frauddetection_db << EOF
$(cat backend/internal/database/migrations/001_create_users.sql)
$(cat backend/internal/database/migrations/002_create_messages.sql)
$(cat backend/internal/database/migrations/003_create_verifications.sql)
$(cat backend/internal/database/migrations/004_create_reports.sql)
$(cat backend/internal/database/migrations/005_create_rbi_circulars.sql)
$(cat backend/internal/database/migrations/006_create_sender_registry.sql)
EOF

# Seed database
echo "🌱 Seeding database..."
docker-compose exec -T postgres psql -U frauddetection -d frauddetection_db < database/seeds/rbi_circulars.sql
docker-compose exec -T postgres psql -U frauddetection -d frauddetection_db < database/seeds/sender_registry.sql

# Check service health
echo "🏥 Checking service health..."
sleep 5

# Check API Gateway
if curl -f http://localhost:8080/health > /dev/null 2>&1; then
    echo "✅ API Gateway is healthy"
else
    echo "❌ API Gateway is not responding"
fi

# Check ML Service
if curl -f http://localhost:8000/health > /dev/null 2>&1; then
    echo "✅ ML Service is healthy"
else
    echo "❌ ML Service is not responding"
fi

echo ""
echo "✅ Setup complete!"
echo ""
echo "🌐 Services are running:"
echo "   - API Gateway: http://localhost:8080"
echo "   - ML Service: http://localhost:8000"
echo "   - PostgreSQL: localhost:5432"
echo "   - Redis: localhost:6379"
echo "   - Kafka: localhost:9093"
echo ""
echo "📚 API Documentation:"
echo "   - Health Check: curl http://localhost:8080/health"
echo "   - Register User: curl -X POST http://localhost:8080/api/v1/auth/register -H 'Content-Type: application/json' -d '{...}'"
echo "   - Verify Message: curl -X POST http://localhost:8080/api/v1/verify -H 'Content-Type: application/json' -d '{...}'"
echo ""
echo "📖 For more information, see docs/API.md"
echo ""
echo "🛑 To stop services: make down"
echo "📋 To view logs: make logs"

