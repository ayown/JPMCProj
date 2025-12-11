# Banking Fraud Detection System

AI-powered real-time fraud detection system for banking SMS/message verification.

## Features

- 🤖 Multi-model ML ensemble for fraud detection
- 🔍 Real-time RBI regulatory compliance verification
- 🛡️ Header authentication and sender verification
- ⚡ High-throughput distributed processing with Kafka
- 🚀 Microservices architecture with Go and Python

## Architecture

### Backend Services (Go)
- **API Gateway**: REST endpoints for all services
- **Auth Service**: JWT-based authentication
- **Verification Service**: Message fraud verification orchestrator
- **Worker Service**: Kafka consumer for async processing

### ML Service (Python)
- **Inference API**: FastAPI-based ML model serving
- **Multi-Model Ensemble**: DistilBERT, RoBERTa, LSTM, XGBoost
- **Training Pipeline**: Continuous model training and evaluation

### Infrastructure
- **PostgreSQL**: Primary database
- **Redis**: Caching layer
- **Kafka**: Message queue for async processing
- **Docker Compose**: Container orchestration

## Quick Start

### Prerequisites
- Docker & Docker Compose
- Go 1.21+
- Python 3.10+
- Make

### Setup

1. Clone the repository:
```bash
git clone <repository-url>
cd fraud-detection-system
```

2. Copy environment variables:
```bash
cp .env.example .env
# Edit .env with your configuration
```

3. Start all services:
```bash
make up
```

4. Run database migrations:
```bash
make migrate
```

5. Seed initial data:
```bash
make seed
```

6. Train ML models (optional - pre-trained models included):
```bash
make train-models
```

### Verify Installation

```bash
# Check service health
curl http://localhost:8080/health

# Check ML service
curl http://localhost:8000/health
```

## API Endpoints

### Authentication
- `POST /api/v1/auth/register` - Register new user
- `POST /api/v1/auth/login` - Login user
- `POST /api/v1/auth/refresh` - Refresh JWT token

### Verification
- `POST /api/v1/verify` - Verify message for fraud
- `GET /api/v1/verify/:id` - Get verification result
- `GET /api/v1/verify/history` - Get verification history

### Reports
- `POST /api/v1/reports` - Submit fraud report
- `GET /api/v1/reports` - List reports
- `GET /api/v1/reports/:id` - Get report details

### RBI Compliance
- `GET /api/v1/rbi/circulars` - List RBI circulars
- `GET /api/v1/rbi/verify` - Verify RBI compliance

## Development

### Running Services Individually

```bash
# Backend API Gateway
cd backend
go run cmd/api-gateway/main.go

# ML Service
cd ml-service
uvicorn app.main:app --reload --port 8000

# Worker Service
cd backend
go run cmd/worker/main.go
```

### Running Tests

```bash
# All tests
make test

# Backend tests only
make test-backend

# ML tests only
make test-ml
```

### Database Migrations

```bash
# Create new migration
make migrate-create NAME=migration_name

# Run migrations
make migrate-up

# Rollback migrations
make migrate-down
```

## ML Models

### Ensemble Architecture

1. **DistilBERT** - Fast semantic fraud classification
2. **RoBERTa** - Multi-class fraud type detection
3. **LSTM + Attention** - Sequence pattern detection
4. **XGBoost** - Metadata feature classification

### Training

```bash
# Train all models
make train-models

# Train specific model
cd ml-service
python training/train_distilbert.py
```

### Model Performance Targets
- Precision: > 95%
- Recall: > 90%
- F1 Score: > 92%
- Inference Time: < 300ms

## Configuration

### Environment Variables

See `.env.example` for all configuration options.

Key variables:
- `DATABASE_URL` - PostgreSQL connection string
- `REDIS_URL` - Redis connection string
- `KAFKA_BROKERS` - Kafka broker addresses
- `JWT_SECRET` - JWT signing secret
- `ML_SERVICE_URL` - ML inference service URL

## Monitoring

### Metrics
- Request latency and throughput
- Model prediction distribution
- Error rates
- Queue depth

### Logs
- Structured JSON logging
- Centralized log aggregation
- Audit trails for all predictions

## Security

- End-to-end encryption for message text
- PII detection and masking
- Rate limiting per user/IP
- JWT authentication
- API key management

## Performance

- **Throughput**: 10,000+ verifications/second
- **Latency**: < 500ms (p95)
- **Availability**: 99.9% uptime

## Project Structure

```
fraud-detection-system/
├── backend/          # Go microservices
├── ml-service/       # Python ML service
├── models/           # Trained ML models
├── database/         # DB scripts and migrations
├── nginx/            # Reverse proxy config
├── scripts/          # DevOps scripts
├── docs/             # Documentation
└── tests/            # Integration tests
```

## Contributing

See [CONTRIBUTING.md](docs/CONTRIBUTING.md) for development guidelines.

## License

MIT License - see LICENSE file for details.

## Support

For issues and questions, please open a GitHub issue.

