# 📚 Project Index - Quick Navigation

Welcome to the Banking Fraud Detection System! Use this index to quickly find what you need.

## 🚀 Getting Started

1. **First Time Setup**: Read [`QUICKSTART.md`](QUICKSTART.md)
2. **Run Setup Script**: `./scripts/setup.sh`
3. **Verify Installation**: `./scripts/verify-setup.sh`
4. **Test the System**: Follow [`TESTING.md`](TESTING.md)

## 📖 Documentation

### Essential Docs
- [`README.md`](README.md) - Project overview and introduction
- [`QUICKSTART.md`](QUICKSTART.md) - Get started in 5 minutes
- [`COMPLETE_PROJECT_SUMMARY.md`](COMPLETE_PROJECT_SUMMARY.md) - **Complete overview of everything**

### Technical Docs
- [`docs/API.md`](docs/API.md) - Complete API reference
- [`docs/ARCHITECTURE.md`](docs/ARCHITECTURE.md) - System architecture
- [`docs/DEPLOYMENT.md`](docs/DEPLOYMENT.md) - Deployment guide

### Testing & Development
- [`TESTING.md`](TESTING.md) - Testing guide with examples
- [`MISSING_COMPONENTS_ADDED.md`](MISSING_COMPONENTS_ADDED.md) - What was added to complete the project

### Component-Specific
- [`backend/README.md`](backend/) - Backend services documentation
- [`ml-service/README.md`](ml-service/) - ML service documentation
- [`frontend/README.md`](frontend/README.md) - Frontend documentation
- [`frontend/FRONTEND_COMPLETE.md`](frontend/FRONTEND_COMPLETE.md) - Frontend status

## 🏗️ Project Structure

### Backend (Go)
```
backend/
├── cmd/                    # Service entry points
│   ├── api-gateway/       # Main API gateway
│   ├── auth-service/      # Authentication service
│   ├── verification-service/  # Verification service
│   └── worker/            # Kafka worker
├── internal/              # Internal packages
│   ├── api/              # HTTP handlers & middleware
│   ├── models/           # Data models
│   ├── repository/       # Database layer
│   ├── service/          # Business logic
│   ├── queue/            # Kafka integration
│   └── database/         # DB & migrations
└── tests/                # Tests
```

### ML Service (Python)
```
ml-service/
├── app/                   # FastAPI application
│   ├── api/              # API routes
│   ├── models/           # ML models (4 models)
│   ├── inference/        # Prediction logic
│   └── utils/            # Utilities
├── training/             # Training pipeline
└── tests/                # Tests
```

### Frontend (React/TypeScript)
```
frontend/
├── src/
│   ├── components/       # React components
│   ├── pages/            # Page components
│   ├── services/         # API integration
│   ├── store/            # Redux state
│   ├── hooks/            # Custom hooks
│   └── types/            # TypeScript types
└── public/               # Static assets
```

## 🔧 Common Commands

### Setup & Start
```bash
./scripts/setup.sh          # Complete setup
make up                     # Start services
make down                   # Stop services
make restart                # Restart services
```

### Development
```bash
make logs                   # View all logs
make logs-api              # View API logs
make logs-ml               # View ML logs
make test                  # Run tests
```

### Database
```bash
make migrate-up            # Run migrations
make migrate-down          # Rollback migrations
make seed                  # Seed data
./scripts/backup.sh        # Backup database
```

### Deployment
```bash
./scripts/deploy.sh production    # Deploy to production
./scripts/verify-setup.sh         # Verify installation
```

## 🌐 Service URLs

| Service | URL | Port |
|---------|-----|------|
| Frontend | http://localhost:3000 | 3000 |
| API Gateway | http://localhost:8080 | 8080 |
| Auth Service | http://localhost:8081 | 8081 |
| Verification Service | http://localhost:8082 | 8082 |
| ML Service | http://localhost:8000 | 8000 |
| PostgreSQL | localhost:5432 | 5432 |
| Redis | localhost:6379 | 6379 |
| Kafka | localhost:9093 | 9093 |
| Nginx | http://localhost:80 | 80 |

## 📊 Key Features

### Fraud Detection
- ✅ Link detection
- ✅ Urgency analysis
- ✅ KYC fraud detection
- ✅ Sender verification
- ✅ RBI compliance checking
- ✅ Pattern analysis

### User Features
- ✅ User registration & login
- ✅ Message verification
- ✅ Verification history
- ✅ Report submission
- ✅ Statistics dashboard
- ✅ Educational resources

### Technical Features
- ✅ JWT authentication
- ✅ Rate limiting
- ✅ Caching
- ✅ Async processing
- ✅ Real-time detection
- ✅ Scalable architecture

## 🎯 Quick Actions

### Test Fraud Detection
```bash
curl -X POST http://localhost:8080/api/v1/verify \
  -H "Content-Type: application/json" \
  -d '{
    "content": "URGENT! Update KYC: http://fake.com",
    "sender_header": "FAKE-HDFC"
  }'
```

### Register User
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

### Check Health
```bash
curl http://localhost:8080/health
curl http://localhost:8000/health
```

## 🐛 Troubleshooting

### Service Won't Start
1. Check logs: `docker-compose logs <service>`
2. Restart: `docker-compose restart <service>`
3. Rebuild: `docker-compose up -d --build <service>`

### Database Issues
1. Check status: `docker-compose ps postgres`
2. View logs: `docker-compose logs postgres`
3. Connect: `docker-compose exec postgres psql -U frauddetection`

### Frontend Issues
1. Check logs: `docker-compose logs frontend`
2. Rebuild: `cd frontend && npm install && npm run build`
3. Check API URL in `.env`

## 📚 Learning Path

### Beginner
1. Read `README.md`
2. Run `./scripts/setup.sh`
3. Test with curl commands
4. Explore frontend at localhost:3000

### Intermediate
1. Read `docs/API.md`
2. Study `docs/ARCHITECTURE.md`
3. Explore backend code
4. Run tests

### Advanced
1. Train ML models
2. Customize features
3. Deploy to production
4. Scale infrastructure

## 🔗 Important Files

### Configuration
- `docker-compose.yml` - Service orchestration
- `.env.example` - Environment variables
- `Makefile` - Common commands

### Scripts
- `scripts/setup.sh` - Setup automation
- `scripts/verify-setup.sh` - Verification
- `scripts/deploy.sh` - Deployment
- `scripts/backup.sh` - Database backup

### Database
- `backend/internal/database/migrations/` - SQL migrations
- `database/seeds/` - Seed data

## 💡 Tips

1. **Always start with** `./scripts/setup.sh`
2. **Check health** before testing
3. **Read logs** when debugging
4. **Use Makefile** for common tasks
5. **Backup database** before major changes

## 📞 Need Help?

1. **Check documentation** in `docs/`
2. **View logs** with `make logs`
3. **Run verification** with `./scripts/verify-setup.sh`
4. **Read troubleshooting** in `docs/DEPLOYMENT.md`

## 🎉 Quick Wins

1. ✅ Run setup script
2. ✅ Access frontend
3. ✅ Test fraud detection
4. ✅ View statistics
5. ✅ Submit a report

---

**Ready to start?**

```bash
cd fraud-detection-system
./scripts/setup.sh
```

Then visit: **http://localhost:3000**

**Happy coding! 🚀**

