#!/bin/bash

# Deployment script for Banking Fraud Detection System
# Usage: ./scripts/deploy.sh [environment]

set -e

ENVIRONMENT=${1:-production}

echo "🚀 Deploying Banking Fraud Detection System to $ENVIRONMENT..."

# Color codes
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Check environment
if [ "$ENVIRONMENT" != "production" ] && [ "$ENVIRONMENT" != "staging" ]; then
    echo -e "${RED}❌ Invalid environment: $ENVIRONMENT${NC}"
    echo "Usage: ./scripts/deploy.sh [production|staging]"
    exit 1
fi

echo -e "${YELLOW}⚠️  Deploying to $ENVIRONMENT environment${NC}"
read -p "Are you sure? (yes/no): " confirm
if [ "$confirm" != "yes" ]; then
    echo "Deployment cancelled."
    exit 0
fi

# Pre-deployment checks
echo ""
echo "📋 Running pre-deployment checks..."

# Check if .env exists
if [ ! -f ".env.$ENVIRONMENT" ]; then
    echo -e "${RED}❌ Environment file .env.$ENVIRONMENT not found${NC}"
    exit 1
fi
echo -e "${GREEN}✓${NC} Environment file found"

# Check if Docker is running
if ! docker info > /dev/null 2>&1; then
    echo -e "${RED}❌ Docker is not running${NC}"
    exit 1
fi
echo -e "${GREEN}✓${NC} Docker is running"

# Run tests
echo ""
echo "🧪 Running tests..."
# TODO: Uncomment when tests are implemented
# make test || {
#     echo -e "${RED}❌ Tests failed${NC}"
#     exit 1
# }
echo -e "${YELLOW}⚠${NC} Tests skipped (not implemented)"

# Backup database
echo ""
echo "💾 Backing up database..."
./scripts/backup.sh || {
    echo -e "${RED}❌ Backup failed${NC}"
    exit 1
}
echo -e "${GREEN}✓${NC} Database backed up"

# Build images
echo ""
echo "🔨 Building Docker images..."
docker-compose -f docker-compose.yml -f docker-compose.$ENVIRONMENT.yml build || {
    echo -e "${RED}❌ Build failed${NC}"
    exit 1
}
echo -e "${GREEN}✓${NC} Images built successfully"

# Stop old containers
echo ""
echo "🛑 Stopping old containers..."
docker-compose -f docker-compose.yml -f docker-compose.$ENVIRONMENT.yml down || true

# Start new containers
echo ""
echo "🚀 Starting new containers..."
docker-compose -f docker-compose.yml -f docker-compose.$ENVIRONMENT.yml up -d || {
    echo -e "${RED}❌ Failed to start containers${NC}"
    exit 1
}

# Wait for services to be ready
echo ""
echo "⏳ Waiting for services to be ready..."
sleep 15

# Run migrations
echo ""
echo "📊 Running database migrations..."
make migrate-up || {
    echo -e "${RED}❌ Migrations failed${NC}"
    echo "Rolling back..."
    docker-compose -f docker-compose.yml -f docker-compose.$ENVIRONMENT.yml down
    exit 1
}
echo -e "${GREEN}✓${NC} Migrations completed"

# Health checks
echo ""
echo "🏥 Running health checks..."

# Check API Gateway
if curl -f http://localhost:8080/health > /dev/null 2>&1; then
    echo -e "${GREEN}✓${NC} API Gateway is healthy"
else
    echo -e "${RED}❌ API Gateway health check failed${NC}"
    exit 1
fi

# Check ML Service
if curl -f http://localhost:8000/health > /dev/null 2>&1; then
    echo -e "${GREEN}✓${NC} ML Service is healthy"
else
    echo -e "${RED}❌ ML Service health check failed${NC}"
    exit 1
fi

# Smoke tests
echo ""
echo "🔥 Running smoke tests..."
# TODO: Add smoke tests
echo -e "${YELLOW}⚠${NC} Smoke tests skipped (not implemented)"

# Clean up old images
echo ""
echo "🧹 Cleaning up old images..."
docker image prune -f

echo ""
echo -e "${GREEN}✅ Deployment completed successfully!${NC}"
echo ""
echo "📊 Deployment Summary:"
echo "  Environment: $ENVIRONMENT"
echo "  API Gateway: http://localhost:8080"
echo "  ML Service: http://localhost:8000"
echo "  Deployment time: $(date)"
echo ""
echo "📋 Next steps:"
echo "  1. Monitor logs: make logs"
echo "  2. Check metrics: docker stats"
echo "  3. Verify functionality: run integration tests"
echo ""
echo "🔄 To rollback: ./scripts/rollback.sh"

