# Deployment Guide

## Prerequisites

- Docker 20.10+
- Docker Compose 2.0+
- 8GB RAM minimum
- 20GB disk space

## Quick Start

### 1. Clone and Setup

```bash
git clone <repository-url>
cd fraud-detection-system
chmod +x scripts/setup.sh
./scripts/setup.sh
```

### 2. Manual Setup

If you prefer manual setup:

```bash
# Copy environment file
cp .env.example .env

# Edit .env with your configuration
nano .env

# Build and start services
docker-compose build
docker-compose up -d

# Run migrations
make migrate-up

# Seed database
make seed
```

### 3. Verify Installation

```bash
# Check service health
curl http://localhost:8080/health
curl http://localhost:8000/health

# View logs
docker-compose logs -f

# Check running services
docker-compose ps
```

## Configuration

### Environment Variables

Key variables to configure in `.env`:

```bash
# Security (IMPORTANT!)
JWT_SECRET=your-super-secret-key-here

# Database
DATABASE_PASSWORD=secure-password-here

# Hugging Face (optional, for future ML models)
HUGGINGFACE_TOKEN=your-token-here
```

### Service Ports

- API Gateway: 8080
- Auth Service: 8081
- Verification Service: 8082
- ML Service: 8000
- PostgreSQL: 5432
- Redis: 6379
- Kafka: 9093
- Nginx: 80

## Database Management

### Migrations

```bash
# Run all migrations
make migrate-up

# Rollback last migration
make migrate-down

# Create new migration
make migrate-create NAME=add_new_table
```

### Backup and Restore

```bash
# Backup database
docker-compose exec postgres pg_dump -U frauddetection frauddetection_db > backup.sql

# Restore database
docker-compose exec -T postgres psql -U frauddetection frauddetection_db < backup.sql
```

## Scaling

### Horizontal Scaling

Scale specific services:

```bash
# Scale worker service
docker-compose up -d --scale worker=3

# Scale verification service
docker-compose up -d --scale verification-service=2
```

### Resource Limits

Edit `docker-compose.yml` to set resource limits:

```yaml
services:
  api-gateway:
    deploy:
      resources:
        limits:
          cpus: '1'
          memory: 512M
        reservations:
          cpus: '0.5'
          memory: 256M
```

## Monitoring

### Logs

```bash
# All services
docker-compose logs -f

# Specific service
docker-compose logs -f api-gateway

# Last 100 lines
docker-compose logs --tail=100 ml-service
```

### Health Checks

```bash
# API Gateway
curl http://localhost:8080/health
curl http://localhost:8080/ready

# ML Service
curl http://localhost:8000/health
curl http://localhost:8000/ready
```

### Resource Usage

```bash
# Container stats
docker stats

# Disk usage
docker system df
```

## Troubleshooting

### Service Won't Start

```bash
# Check logs
docker-compose logs <service-name>

# Restart service
docker-compose restart <service-name>

# Rebuild and restart
docker-compose up -d --build <service-name>
```

### Database Connection Issues

```bash
# Check PostgreSQL logs
docker-compose logs postgres

# Verify database is running
docker-compose exec postgres pg_isready -U frauddetection

# Connect to database
docker-compose exec postgres psql -U frauddetection -d frauddetection_db
```

### Kafka Issues

```bash
# Check Kafka logs
docker-compose logs kafka

# List topics
docker-compose exec kafka kafka-topics --list --bootstrap-server localhost:9092

# Check consumer groups
docker-compose exec kafka kafka-consumer-groups --list --bootstrap-server localhost:9092
```

### ML Service Issues

```bash
# Check ML service logs
docker-compose logs ml-service

# Verify models are loaded
curl http://localhost:8000/health

# Restart ML service
docker-compose restart ml-service
```

## Production Deployment

### Security Checklist

- [ ] Change all default passwords
- [ ] Set strong JWT_SECRET
- [ ] Enable HTTPS/TLS
- [ ] Configure firewall rules
- [ ] Enable database encryption
- [ ] Set up VPN for internal services
- [ ] Configure rate limiting
- [ ] Enable audit logging
- [ ] Set up backup automation

### Performance Optimization

1. **Database**
   - Enable connection pooling
   - Configure appropriate indexes
   - Set up read replicas
   - Enable query caching

2. **Redis**
   - Configure persistence
   - Set up Redis cluster
   - Optimize cache TTL

3. **Kafka**
   - Configure replication
   - Optimize partition count
   - Set up monitoring

4. **Application**
   - Enable production mode
   - Optimize worker count
   - Configure timeouts
   - Enable compression

### Monitoring Setup

1. **Prometheus** (Future)
```yaml
# Add to docker-compose.yml
prometheus:
  image: prom/prometheus
  volumes:
    - ./monitoring/prometheus:/etc/prometheus
  ports:
    - "9090:9090"
```

2. **Grafana** (Future)
```yaml
grafana:
  image: grafana/grafana
  ports:
    - "3000:3000"
  environment:
    - GF_SECURITY_ADMIN_PASSWORD=admin
```

## Backup Strategy

### Automated Backups

```bash
# Add to crontab
0 2 * * * /path/to/scripts/backup.sh
```

### Backup Script

```bash
#!/bin/bash
DATE=$(date +%Y%m%d_%H%M%S)
docker-compose exec -T postgres pg_dump -U frauddetection frauddetection_db | gzip > backup_$DATE.sql.gz
# Upload to S3 or backup server
```

## Update and Maintenance

### Update Services

```bash
# Pull latest images
docker-compose pull

# Rebuild services
docker-compose build

# Restart with new images
docker-compose up -d
```

### Database Maintenance

```bash
# Vacuum database
docker-compose exec postgres psql -U frauddetection -d frauddetection_db -c "VACUUM ANALYZE;"

# Check database size
docker-compose exec postgres psql -U frauddetection -d frauddetection_db -c "SELECT pg_size_pretty(pg_database_size('frauddetection_db'));"
```

## Cleanup

### Remove Old Data

```bash
# Remove old verifications (non-fraud, older than 30 days)
docker-compose exec postgres psql -U frauddetection -d frauddetection_db -c "DELETE FROM verifications WHERE is_fraud = false AND created_at < NOW() - INTERVAL '30 days';"
```

### Clean Docker

```bash
# Remove unused containers
docker container prune

# Remove unused images
docker image prune

# Remove unused volumes
docker volume prune

# Clean everything
docker system prune -a
```

## Shutdown

```bash
# Stop all services
docker-compose down

# Stop and remove volumes
docker-compose down -v
```

