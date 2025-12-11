# System Architecture

## Overview

The Banking Fraud Detection System is built using a microservices architecture with the following key components:

## Architecture Diagram

```
┌─────────────┐
│   Client    │
└──────┬──────┘
       │
       ▼
┌─────────────┐
│    Nginx    │ (Reverse Proxy)
└──────┬──────┘
       │
       ▼
┌─────────────────────────────────────────────┐
│           API Gateway (Go)                  │
│  - Request routing                          │
│  - Authentication                           │
│  - Rate limiting                            │
└──────┬──────────────────────────────────────┘
       │
       ├──────────────┬──────────────┬─────────────┐
       ▼              ▼              ▼             ▼
┌─────────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐
│Auth Service │ │Verification│ │  Worker  │ │ML Service│
│    (Go)     │ │Service(Go) │ │  (Go)    │ │ (Python) │
└──────┬──────┘ └─────┬──────┘ └────┬─────┘ └────┬─────┘
       │              │              │            │
       │              └──────┬───────┘            │
       │                     │                    │
       ▼                     ▼                    ▼
┌─────────────────────────────────────────────────┐
│              PostgreSQL Database                │
└─────────────────────────────────────────────────┘

       │                     │                    │
       ▼                     ▼                    ▼
┌─────────────┐       ┌──────────┐       ┌──────────┐
│   Redis     │       │  Kafka   │       │ Models   │
│  (Cache)    │       │ (Queue)  │       │ Storage  │
└─────────────┘       └──────────┘       └──────────┘
```

## Components

### 1. API Gateway (Go)

**Responsibilities:**
- Central entry point for all client requests
- Request routing to appropriate microservices
- JWT authentication and authorization
- Rate limiting and throttling
- Request/response logging
- CORS handling

**Technology:** Go, Gin framework

### 2. Auth Service (Go)

**Responsibilities:**
- User registration and authentication
- JWT token generation and validation
- Password hashing and verification
- User profile management
- Session management

**Technology:** Go, bcrypt, JWT

### 3. Verification Service (Go)

**Responsibilities:**
- Message fraud verification orchestration
- Feature extraction from messages
- ML service integration
- RBI compliance checking
- Header verification
- Result aggregation and scoring

**Technology:** Go, Gin framework

### 4. Worker Service (Go)

**Responsibilities:**
- Asynchronous message processing
- Kafka message consumption
- Background verification tasks
- Batch processing
- Retry logic for failed verifications

**Technology:** Go, Kafka consumer

### 5. ML Service (Python)

**Responsibilities:**
- Fraud detection inference
- Multi-model ensemble predictions
- Feature engineering
- Model serving
- Result caching
- Model versioning

**Technology:** Python, FastAPI, scikit-learn, transformers

**Models:**
- Rule-based detection (MVP)
- Future: DistilBERT, RoBERTa, LSTM, XGBoost

### 6. PostgreSQL Database

**Schema:**
- `users` - User accounts
- `messages` - Verified messages
- `verifications` - Verification results
- `reports` - User-submitted fraud reports
- `rbi_circulars` - RBI regulatory data
- `sender_registry` - Verified sender IDs

### 7. Redis Cache

**Usage:**
- ML prediction caching
- Rate limiting counters
- Session storage
- Temporary data storage

### 8. Kafka Message Queue

**Topics:**
- `verification-requests` - Async verification requests
- `fraud-reports` - Fraud report events
- `fraud-alerts` - High-priority fraud alerts

### 9. Nginx Reverse Proxy

**Responsibilities:**
- Load balancing
- SSL termination (production)
- Static file serving
- Request routing

## Data Flow

### Synchronous Verification Flow

1. Client sends verification request to API Gateway
2. API Gateway validates JWT (if provided) and applies rate limiting
3. Request forwarded to Verification Service
4. Verification Service:
   - Extracts message features
   - Calls ML Service for prediction
   - Calls RBI Compliance Service
   - Calls Header Verification Service
   - Aggregates results
5. Verification result saved to PostgreSQL
6. Response returned to client

### Asynchronous Verification Flow

1. Client sends verification request to API Gateway
2. API Gateway publishes message to Kafka
3. Worker Service consumes message from Kafka
4. Worker performs verification (same as synchronous)
5. Result saved to database
6. Optional: Alert sent via Kafka if high-risk fraud detected

## Security

### Authentication & Authorization
- JWT-based authentication
- Token expiration and refresh
- Role-based access control (future)

### Data Protection
- Password hashing with bcrypt
- PII masking in logs
- Encrypted database connections
- HTTPS in production

### Rate Limiting
- Per-user rate limits
- Per-IP rate limits
- Configurable limits per endpoint

## Scalability

### Horizontal Scaling
- Stateless services can be scaled horizontally
- Load balancing via Nginx
- Database connection pooling

### Caching Strategy
- ML predictions cached in Redis
- Cache invalidation on model updates
- TTL-based cache expiration

### Queue-based Processing
- Kafka for async processing
- Multiple worker instances
- Consumer groups for load distribution

## Monitoring & Observability

### Logging
- Structured JSON logging
- Log levels: DEBUG, INFO, WARN, ERROR
- Request/response logging
- Error tracking

### Metrics (Future)
- Request latency
- Throughput
- Error rates
- Model prediction distribution
- Queue depth

### Health Checks
- `/health` - Basic health check
- `/ready` - Readiness probe
- Service dependency checks

## Deployment

### Docker Compose (Development)
- All services in containers
- Shared network
- Volume mounts for development
- Hot reload enabled

### Kubernetes (Production - Future)
- Helm charts
- Auto-scaling
- Rolling updates
- Health checks and probes

## Performance Targets

- **API Latency**: < 500ms (p95)
- **ML Inference**: < 300ms
- **Database Queries**: < 50ms
- **Throughput**: 10,000+ verifications/second
- **Availability**: 99.9%

## Future Enhancements

1. **Advanced ML Models**
   - Train and deploy transformer models
   - Ensemble learning
   - Continuous learning pipeline

2. **Real-time Features**
   - WebSocket support for live alerts
   - Streaming analytics

3. **Advanced Analytics**
   - Fraud trend analysis
   - Predictive analytics
   - Anomaly detection

4. **Multi-tenancy**
   - Bank-specific models
   - Custom fraud rules
   - White-labeling

5. **Mobile SDKs**
   - Native iOS/Android integration
   - On-device ML inference

