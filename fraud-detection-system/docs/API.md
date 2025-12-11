# Banking Fraud Detection System - API Documentation

## Base URL

```
http://localhost:8080/api/v1
```

## Authentication

Most endpoints require JWT authentication. Include the token in the Authorization header:

```
Authorization: Bearer <your_jwt_token>
```

## Endpoints

### Authentication

#### Register User

```http
POST /auth/register
```

**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "SecurePass123!",
  "full_name": "John Doe",
  "phone_number": "+919876543210"
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "email": "user@example.com",
    "full_name": "John Doe",
    "phone_number": "+919876543210",
    "is_active": true,
    "is_verified": false,
    "created_at": "2024-01-01T00:00:00Z"
  }
}
```

#### Login

```http
POST /auth/login
```

**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "SecurePass123!"
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "access_token": "eyJhbGc...",
    "refresh_token": "eyJhbGc...",
    "expires_in": 86400,
    "token_type": "Bearer"
  }
}
```

#### Refresh Token

```http
POST /auth/refresh
```

**Request Body:**
```json
{
  "refresh_token": "eyJhbGc..."
}
```

### Verification

#### Verify Message

```http
POST /verify
```

**Request Body:**
```json
{
  "content": "Dear customer, your account will be blocked. Update KYC immediately: http://suspicious-link.com",
  "sender_header": "AX-HDFC",
  "message_type": "SMS",
  "received_at": "2024-01-01T10:00:00Z"
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "message_id": "uuid",
    "is_fraud": true,
    "fraud_score": 0.85,
    "fraud_type": "kyc_fraud",
    "confidence": 0.92,
    "risk_level": "HIGH",
    "header_verified": true,
    "rbi_compliant": false,
    "explanation": "⚠️ This message has been flagged as potentially fraudulent. ML Analysis: High fraud probability detected. Sender Verification: Sender ID 'AX-HDFC' is verified but message content is suspicious. RBI Compliance: Message uses regulatory language but does not match any official RBI circulars",
    "recommendations": [
      "Do not click on any links in this message",
      "Do not share personal or financial information",
      "Do not call any phone numbers provided in the message",
      "The sender is not verified - this is likely a spoofed message",
      "This message references fake or expired regulatory requirements",
      "Report this message to your bank and regulatory authorities",
      "Delete this message immediately"
    ],
    "model_predictions": {
      "rule_based_model": {
        "fraud_score": 0.85,
        "is_fraud": true,
        "indicators": [
          "Contains suspicious links",
          "Uses urgent language (2 urgent words)",
          "References KYC/regulatory requirements"
        ]
      }
    },
    "processing_time_ms": 245,
    "verified_at": "2024-01-01T10:00:01Z"
  }
}
```

#### Get Verification by ID

```http
GET /verify/:id
```

#### Get Verification History

```http
GET /verify/history?limit=20&offset=0
```

**Requires Authentication**

#### Get Verification Statistics

```http
GET /verify/stats
```

**Response:**
```json
{
  "success": true,
  "data": {
    "total_verifications": 1000,
    "fraud_detected": 350,
    "fraud_rate": 0.35,
    "avg_fraud_score": 0.42,
    "avg_processing_time": 250.5,
    "last_24_hours": 150,
    "last_7_days": 800
  }
}
```

### Reports

#### Submit Fraud Report

```http
POST /reports
```

**Request Body:**
```json
{
  "report_type": "FRAUD",
  "content": "Suspicious message content",
  "sender_header": "FAKE-BANK",
  "description": "This message is clearly a phishing attempt",
  "verification_id": "uuid"
}
```

#### Get Report by ID

```http
GET /reports/:id
```

#### Get User Reports

```http
GET /reports?limit=20&offset=0
```

**Requires Authentication**

#### Get Report Statistics

```http
GET /reports/stats
```

### Health Check

#### API Health

```http
GET /health
```

**Response:**
```json
{
  "status": "healthy",
  "service": "fraud-detection-api"
}
```

#### Readiness Check

```http
GET /ready
```

**Response:**
```json
{
  "status": "ready",
  "checks": {
    "database": "healthy",
    "cache": "healthy",
    "ml": "healthy"
  }
}
```

## Error Responses

All errors follow this format:

```json
{
  "error": "error_type",
  "message": "Human-readable error message",
  "code": 400
}
```

### Common Error Codes

- `400` - Bad Request (invalid input)
- `401` - Unauthorized (missing or invalid token)
- `403` - Forbidden (insufficient permissions)
- `404` - Not Found
- `429` - Too Many Requests (rate limit exceeded)
- `500` - Internal Server Error

## Rate Limiting

- **Authentication endpoints**: 50 requests per minute
- **Verification endpoints**: 100 requests per minute
- **Other endpoints**: 100 requests per minute

Rate limit headers are included in responses:
```
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 95
```

## Examples

### cURL Examples

#### Register and Login

```bash
# Register
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "SecurePass123!",
    "full_name": "Test User",
    "phone_number": "+919876543210"
  }'

# Login
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "SecurePass123!"
  }'
```

#### Verify Message

```bash
curl -X POST http://localhost:8080/api/v1/verify \
  -H "Content-Type: application/json" \
  -d '{
    "content": "URGENT: Your account will be suspended. Update KYC now: http://fake-bank.com",
    "sender_header": "FAKE-BANK"
  }'
```

#### Verify Message (Authenticated)

```bash
TOKEN="your_jwt_token_here"

curl -X POST http://localhost:8080/api/v1/verify \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "content": "Your transaction of Rs. 5000 is successful. Ref: ABC123",
    "sender_header": "AX-HDFC"
  }'
```

## ML Service API

The ML service has its own API (internal use):

### Predict Fraud

```http
POST http://localhost:8000/api/v1/predict
```

**Request Body:**
```json
{
  "content": "message content",
  "sender_header": "sender-id",
  "features": {
    "content": "message content",
    "sender_header": "sender-id",
    "message_length": 150,
    "has_links": true,
    "link_count": 1,
    "extracted_urls": ["http://example.com"],
    "has_phone_number": false,
    "phone_number_count": 0,
    "has_urgent_words": true,
    "urgent_word_count": 2,
    "special_char_ratio": 0.05,
    "capital_ratio": 0.15,
    "number_ratio": 0.10,
    "has_kyc_keywords": true,
    "has_bank_names": true
  }
}
```

## WebSocket Support (Future)

Real-time fraud alerts will be available via WebSocket in future versions:

```javascript
const ws = new WebSocket('ws://localhost:8080/ws/alerts');
ws.onmessage = (event) => {
  const alert = JSON.parse(event.data);
  console.log('Fraud alert:', alert);
};
```

