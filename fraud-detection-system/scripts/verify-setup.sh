#!/bin/bash

# Verification script for Banking Fraud Detection System

echo "🔍 Verifying Banking Fraud Detection System Setup..."
echo ""

ERRORS=0
WARNINGS=0

# Color codes
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Check function
check() {
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}✓${NC} $1"
    else
        echo -e "${RED}✗${NC} $1"
        ((ERRORS++))
    fi
}

warn() {
    echo -e "${YELLOW}⚠${NC} $1"
    ((WARNINGS++))
}

# 1. Check Docker
echo "1. Checking Docker..."
docker --version > /dev/null 2>&1
check "Docker installed"

docker-compose --version > /dev/null 2>&1
check "Docker Compose installed"

# 2. Check if services are running
echo ""
echo "2. Checking services..."

docker-compose ps | grep -q "api-gateway.*Up"
check "API Gateway running"

docker-compose ps | grep -q "ml-service.*Up"
check "ML Service running"

docker-compose ps | grep -q "postgres.*Up"
check "PostgreSQL running"

docker-compose ps | grep -q "redis.*Up"
check "Redis running"

docker-compose ps | grep -q "kafka.*Up"
check "Kafka running"

docker-compose ps | grep -q "worker.*Up"
check "Worker Service running"

# 3. Check service health
echo ""
echo "3. Checking service health..."

curl -f http://localhost:8080/health > /dev/null 2>&1
check "API Gateway health"

curl -f http://localhost:8000/health > /dev/null 2>&1
check "ML Service health"

# 4. Check database
echo ""
echo "4. Checking database..."

docker-compose exec -T postgres pg_isready -U frauddetection > /dev/null 2>&1
check "PostgreSQL ready"

# Check if tables exist
TABLES=$(docker-compose exec -T postgres psql -U frauddetection -d frauddetection_db -t -c "SELECT COUNT(*) FROM information_schema.tables WHERE table_schema='public';")
if [ "$TABLES" -ge 6 ]; then
    echo -e "${GREEN}✓${NC} Database tables created ($TABLES tables)"
else
    echo -e "${RED}✗${NC} Database tables missing (found $TABLES, expected 6+)"
    ((ERRORS++))
fi

# Check if seed data exists
SENDERS=$(docker-compose exec -T postgres psql -U frauddetection -d frauddetection_db -t -c "SELECT COUNT(*) FROM sender_registry;")
if [ "$SENDERS" -gt 0 ]; then
    echo -e "${GREEN}✓${NC} Sender registry seeded ($SENDERS senders)"
else
    warn "Sender registry not seeded"
fi

CIRCULARS=$(docker-compose exec -T postgres psql -U frauddetection -d frauddetection_db -t -c "SELECT COUNT(*) FROM rbi_circulars;")
if [ "$CIRCULARS" -gt 0 ]; then
    echo -e "${GREEN}✓${NC} RBI circulars seeded ($CIRCULARS circulars)"
else
    warn "RBI circulars not seeded"
fi

# 5. Check Redis
echo ""
echo "5. Checking Redis..."

docker-compose exec -T redis redis-cli ping > /dev/null 2>&1
check "Redis responding"

# 6. Check Kafka
echo ""
echo "6. Checking Kafka..."

docker-compose exec -T kafka kafka-broker-api-versions --bootstrap-server localhost:9092 > /dev/null 2>&1
check "Kafka broker responding"

# 7. Check file structure
echo ""
echo "7. Checking file structure..."

[ -f ".env" ] && check ".env file exists" || warn ".env file missing (using defaults)"
[ -f "docker-compose.yml" ] && check "docker-compose.yml exists"
[ -f "Makefile" ] && check "Makefile exists"
[ -d "backend" ] && check "backend directory exists"
[ -d "ml-service" ] && check "ml-service directory exists"
[ -d "database" ] && check "database directory exists"

# 8. Test API endpoints
echo ""
echo "8. Testing API endpoints..."

# Test registration
REGISTER_RESPONSE=$(curl -s -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d "{\"email\":\"test$(date +%s)@example.com\",\"password\":\"SecurePass123!\",\"full_name\":\"Test User\",\"phone_number\":\"+919876543210\"}")

if echo "$REGISTER_RESPONSE" | grep -q "success"; then
    check "User registration works"
else
    echo -e "${RED}✗${NC} User registration failed"
    ((ERRORS++))
fi

# Test verification (without auth)
VERIFY_RESPONSE=$(curl -s -X POST http://localhost:8080/api/v1/verify \
  -H "Content-Type: application/json" \
  -d '{"content":"Test message","sender_header":"TEST"}')

if echo "$VERIFY_RESPONSE" | grep -q "fraud_score"; then
    check "Message verification works"
else
    echo -e "${RED}✗${NC} Message verification failed"
    ((ERRORS++))
fi

# 9. Check logs for errors
echo ""
echo "9. Checking for errors in logs..."

API_ERRORS=$(docker-compose logs api-gateway 2>&1 | grep -i "error\|fatal" | wc -l)
ML_ERRORS=$(docker-compose logs ml-service 2>&1 | grep -i "error\|fatal" | wc -l)

if [ "$API_ERRORS" -eq 0 ]; then
    check "No errors in API Gateway logs"
else
    warn "Found $API_ERRORS errors in API Gateway logs"
fi

if [ "$ML_ERRORS" -eq 0 ]; then
    check "No errors in ML Service logs"
else
    warn "Found $ML_ERRORS errors in ML Service logs"
fi

# 10. Check resource usage
echo ""
echo "10. Checking resource usage..."

# Get container stats
STATS=$(docker stats --no-stream --format "table {{.Name}}\t{{.CPUPerc}}\t{{.MemUsage}}")
echo "$STATS"

# Summary
echo ""
echo "═══════════════════════════════════════════════"
if [ $ERRORS -eq 0 ] && [ $WARNINGS -eq 0 ]; then
    echo -e "${GREEN}✓ All checks passed!${NC}"
    echo ""
    echo "System is ready to use:"
    echo "  API Gateway: http://localhost:8080"
    echo "  ML Service: http://localhost:8000"
    echo ""
    echo "Next steps:"
    echo "  1. Read QUICKSTART.md for usage examples"
    echo "  2. Read docs/API.md for API documentation"
    echo "  3. Run TESTING.md test cases"
    exit 0
elif [ $ERRORS -eq 0 ]; then
    echo -e "${YELLOW}⚠ Setup complete with $WARNINGS warning(s)${NC}"
    echo ""
    echo "System is functional but has some warnings."
    echo "Review the warnings above and fix if necessary."
    exit 0
else
    echo -e "${RED}✗ Setup incomplete: $ERRORS error(s), $WARNINGS warning(s)${NC}"
    echo ""
    echo "Please fix the errors above before proceeding."
    echo "Common fixes:"
    echo "  - Run: docker-compose up -d"
    echo "  - Run: make migrate-up"
    echo "  - Run: make seed"
    echo "  - Check logs: docker-compose logs"
    exit 1
fi

