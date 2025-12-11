# Testing the Banking Fraud Detection System

## Quick Test Guide

### 1. Start the System

```bash
cd fraud-detection-system
./scripts/setup.sh
```

Wait for all services to start (about 30 seconds).

### 2. Test Health Checks

```bash
# API Gateway
curl http://localhost:8080/health

# ML Service
curl http://localhost:8000/health
```

Expected response: `{"status":"healthy",...}`

### 3. Register a User

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

### 4. Login

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "SecurePass123!"
  }'
```

Save the `access_token` from the response.

### 5. Test Fraud Detection

#### Test Case 1: Obvious Fraud Message

```bash
curl -X POST http://localhost:8080/api/v1/verify \
  -H "Content-Type: application/json" \
  -d '{
    "content": "URGENT! Your HDFC account will be blocked. Update KYC immediately by clicking: http://fake-hdfc-bank.com/kyc or call 9999999999",
    "sender_header": "FAKE-HDFC"
  }'
```

Expected: High fraud score (> 0.7), multiple fraud indicators

#### Test Case 2: Legitimate Message

```bash
curl -X POST http://localhost:8080/api/v1/verify \
  -H "Content-Type: application/json" \
  -d '{
    "content": "Your transaction of Rs. 5000 at Amazon is successful. Available balance: Rs. 45000. Thank you for banking with us.",
    "sender_header": "AX-HDFC"
  }'
```

Expected: Low fraud score (< 0.3), legitimate classification

#### Test Case 3: Suspicious KYC Request

```bash
curl -X POST http://localhost:8080/api/v1/verify \
  -H "Content-Type: application/json" \
  -d '{
    "content": "Dear customer, RBI has made KYC update mandatory. Your account will be suspended if not updated within 24 hours. Click here to update: http://rbi-kyc-update.com",
    "sender_header": "RBI-BANK"
  }'
```

Expected: High fraud score, KYC fraud type, RBI compliance failure

#### Test Case 4: Phishing with Urgency

```bash
curl -X POST http://localhost:8080/api/v1/verify \
  -H "Content-Type: application/json" \
  -d '{
    "content": "ALERT! Suspicious activity detected on your account. Verify immediately to prevent blocking: http://verify-account-now.com",
    "sender_header": "ALERT-SBI"
  }'
```

Expected: High fraud score, phishing type

### 6. Test with Authentication

```bash
# Use the token from login
TOKEN="your_access_token_here"

curl -X POST http://localhost:8080/api/v1/verify \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "content": "Test message",
    "sender_header": "TEST"
  }'
```

### 7. Get Verification History

```bash
curl -X GET "http://localhost:8080/api/v1/verify/history?limit=10" \
  -H "Authorization: Bearer $TOKEN"
```

### 8. Get Statistics

```bash
curl http://localhost:8080/api/v1/verify/stats
```

### 9. Submit a Fraud Report

```bash
curl -X POST http://localhost:8080/api/v1/reports \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "report_type": "FRAUD",
    "content": "Suspicious message I received",
    "sender_header": "FAKE-SENDER",
    "description": "This is clearly a phishing attempt"
  }'
```

## Expected Fraud Indicators

The system detects fraud based on these indicators:

1. **Suspicious Links** (+0.25 score)
   - Presence of URLs in the message

2. **Urgent Language** (+0.15 score)
   - Words like: urgent, immediately, expire, blocked, suspended

3. **KYC/Regulatory Keywords** (+0.20 score)
   - References to KYC, RBI, mandatory updates

4. **Excessive Special Characters** (+0.10 score)
   - Ratio > 15%

5. **Excessive Capitals** (+0.10 score)
   - Ratio > 50%

6. **Phone Number + Urgency** (+0.15 score)
   - Callback scam indicator

7. **Bank + Urgency + Links** (+0.20 score)
   - High-risk combination

## Testing Sender Verification

### Verified Senders (Should have lower fraud scores)

```bash
# Test with verified HDFC sender
curl -X POST http://localhost:8080/api/v1/verify \
  -H "Content-Type: application/json" \
  -d '{
    "content": "Your OTP for transaction is 123456. Valid for 10 minutes.",
    "sender_header": "AX-HDFC"
  }'
```

### Unverified Senders (Should have higher fraud scores)

```bash
# Test with unverified sender
curl -X POST http://localhost:8080/api/v1/verify \
  -H "Content-Type: application/json" \
  -d '{
    "content": "Your OTP for transaction is 123456. Valid for 10 minutes.",
    "sender_header": "UNKNOWN-BANK"
  }'
```

## Performance Testing

### Load Test with Apache Bench

```bash
# Install Apache Bench
sudo apt-get install apache2-utils

# Test 1000 requests with 10 concurrent connections
ab -n 1000 -c 10 -p test-payload.json -T application/json \
  http://localhost:8080/api/v1/verify
```

### Expected Performance

- **Latency**: < 500ms (p95)
- **Throughput**: > 100 requests/second (single instance)
- **ML Inference**: < 300ms

## Database Verification

### Check Data

```bash
# Connect to database
docker-compose exec postgres psql -U frauddetection -d frauddetection_db

# Check users
SELECT * FROM users;

# Check verifications
SELECT id, is_fraud, fraud_score, fraud_type FROM verifications ORDER BY created_at DESC LIMIT 10;

# Check sender registry
SELECT sender_id, bank_name, reputation_score FROM sender_registry WHERE is_verified = true;

# Exit
\q
```

## Troubleshooting Tests

### Test Fails with Connection Error

```bash
# Check if services are running
docker-compose ps

# Check logs
docker-compose logs api-gateway
docker-compose logs ml-service
```

### Test Fails with 401 Unauthorized

- Token might be expired (24h expiry)
- Re-login to get a new token

### Test Fails with 429 Rate Limit

- Wait 1 minute
- Rate limit: 100 requests per minute

### ML Service Not Responding

```bash
# Restart ML service
docker-compose restart ml-service

# Check health
curl http://localhost:8000/health
```

## Integration Test Script

Save this as `test.sh`:

```bash
#!/bin/bash

API_URL="http://localhost:8080/api/v1"

echo "Testing Banking Fraud Detection System..."

# Test 1: Health Check
echo "1. Health Check..."
curl -s $API_URL/../health | jq .

# Test 2: Register
echo "2. Registering user..."
REGISTER_RESPONSE=$(curl -s -X POST $API_URL/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test'$(date +%s)'@example.com",
    "password": "SecurePass123!",
    "full_name": "Test User",
    "phone_number": "+919876543210"
  }')
echo $REGISTER_RESPONSE | jq .

# Test 3: Login
echo "3. Logging in..."
LOGIN_RESPONSE=$(curl -s -X POST $API_URL/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "SecurePass123!"
  }')
TOKEN=$(echo $LOGIN_RESPONSE | jq -r '.data.access_token')
echo "Token: $TOKEN"

# Test 4: Verify Fraud Message
echo "4. Testing fraud detection..."
curl -s -X POST $API_URL/verify \
  -H "Content-Type: application/json" \
  -d '{
    "content": "URGENT! Update KYC now: http://fake-bank.com",
    "sender_header": "FAKE"
  }' | jq .

echo "All tests completed!"
```

Run with: `chmod +x test.sh && ./test.sh`

## Expected Test Results

All tests should pass with:
- ✅ Services healthy
- ✅ User registration successful
- ✅ Login returns valid token
- ✅ Fraud messages detected (score > 0.5)
- ✅ Legitimate messages pass (score < 0.5)
- ✅ Response time < 500ms

