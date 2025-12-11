# Quick Start Guide

## Prerequisites

- Docker & Docker Compose installed
- 8GB RAM minimum
- 20GB free disk space

## Setup (5 minutes)

### 1. Clone and Navigate

```bash
cd fraud-detection-system
```

### 2. Configure Environment

```bash
cp .env.example .env
# Edit .env and set JWT_SECRET to a secure random string
nano .env
```

### 3. Run Setup Script

```bash
chmod +x scripts/setup.sh
./scripts/setup.sh
```

The script will:
- Build Docker images
- Start all services
- Run database migrations
- Seed initial data
- Verify service health

## Test the System

### 1. Check Services

```bash
# API Gateway
curl http://localhost:8080/health

# ML Service
curl http://localhost:8000/health
```

### 2. Register a User

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "SecurePass123!",
    "full_name": "Test User",
    "phone_number": "+919876543210"
  }'
```

### 3. Test Fraud Detection

```bash
# Test a fraudulent message
curl -X POST http://localhost:8080/api/v1/verify \
  -H "Content-Type: application/json" \
  -d '{
    "content": "URGENT! Your account will be blocked. Update KYC now: http://fake-bank.com",
    "sender_header": "FAKE-HDFC"
  }'
```

Expected response:
```json
{
  "success": true,
  "data": {
    "is_fraud": true,
    "fraud_score": 0.85,
    "fraud_type": "kyc_fraud",
    "risk_level": "HIGH",
    "explanation": "⚠️ This message has been flagged as potentially fraudulent...",
    "recommendations": [
      "Do not click on any links in this message",
      "Do not share personal or financial information",
      ...
    ]
  }
}
```

## Common Commands

```bash
# View logs
make logs

# Stop services
make down

# Restart services
make restart

# Run tests
make test

# View database
docker-compose exec postgres psql -U frauddetection -d frauddetection_db
```

## Access Points

- **API Gateway**: http://localhost:8080
- **ML Service**: http://localhost:8000
- **PostgreSQL**: localhost:5432
- **Redis**: localhost:6379
- **Kafka**: localhost:9093

## Next Steps

1. Read [API Documentation](docs/API.md)
2. Explore [Architecture](docs/ARCHITECTURE.md)
3. Run [Integration Tests](TESTING.md)
4. Check [Deployment Guide](docs/DEPLOYMENT.md)

## Troubleshooting

### Services won't start

```bash
# Check logs
docker-compose logs

# Restart services
docker-compose restart
```

### Port already in use

Edit `docker-compose.yml` to change port mappings.

### Database connection failed

```bash
# Check PostgreSQL
docker-compose logs postgres

# Restart PostgreSQL
docker-compose restart postgres
```

## Support

For issues, check:
1. Service logs: `docker-compose logs <service-name>`
2. Service status: `docker-compose ps`
3. Resource usage: `docker stats`

## Clean Up

```bash
# Stop and remove everything
make clean

# Or manually
docker-compose down -v
```

