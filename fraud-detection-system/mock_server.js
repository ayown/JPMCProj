const express = require("express");
const cors = require("cors");

// --- API Gateway Mock (Port 8080) ---
const apiApp = express();
apiApp.use(cors());
apiApp.use(express.json());

// Logging middleware
apiApp.use((req, res, next) => {
  console.log(`[API Gateway] ${req.method} ${req.url}`);
  next();
});

// Health
apiApp.get("/health", (req, res) => res.status(200).send("OK"));
apiApp.get("/ready", (req, res) => res.status(200).send("Ready"));

// Auth
apiApp.post("/api/v1/auth/register", (req, res) => {
  res
    .status(201)
    .json({ message: "User registered successfully", userId: "123" });
});
apiApp.post("/api/v1/auth/login", (req, res) => {
  res.status(200).json({
    token: "mock-jwt-token-xyz",
    user: { id: "123", email: "test@example.com", name: "Test User" },
  });
});
apiApp.post("/api/v1/auth/refresh", (req, res) => {
  res.status(200).json({ token: "mock-refreshed-jwt-token" });
});
apiApp.get("/api/v1/profile", (req, res) => {
  res
    .status(200)
    .json({
      id: "123",
      email: "test@example.com",
      name: "Test User",
      role: "user",
    });
});

// Verification
apiApp.post("/api/v1/verify", (req, res) => {
  res.status(200).json({
    id: "v-123",
    status: "safe",
    riskScore: 0.05,
    timestamp: new Date().toISOString(),
  });
});
apiApp.get("/api/v1/verify/stats", (req, res) => {
  res.status(200).json({ total: 100, safe: 90, suspicious: 5, fraud: 5 });
});
apiApp.get("/api/v1/verify/history", (req, res) => {
  res.status(200).json([
    { id: "v-1", status: "safe", timestamp: new Date().toISOString() },
    { id: "v-2", status: "suspicious", timestamp: new Date().toISOString() },
  ]);
});
apiApp.get("/api/v1/verify/:id", (req, res) => {
  res.status(200).json({ id: req.params.id, status: "safe", riskScore: 0.1 });
});

// Reports
apiApp.post("/api/v1/reports", (req, res) => {
  res.status(201).json({ id: "r-123", status: "received" });
});
apiApp.get("/api/v1/reports", (req, res) => {
  res
    .status(200)
    .json([{ id: "r-1", description: "Fraud report 1", status: "pending" }]);
});
apiApp.get("/api/v1/reports/stats", (req, res) => {
  res.status(200).json({ total: 10, pending: 5, resolved: 5 });
});
apiApp.get("/api/v1/reports/:id", (req, res) => {
  res
    .status(200)
    .json({
      id: req.params.id,
      description: "Fraud report details",
      status: "pending",
    });
});

apiApp.listen(8080, "0.0.0.0", () => {
  console.log("🚀 API Gateway Mock running on port 8080");
});

// --- ML Service Mock (Port 8000) ---
const mlApp = express();
mlApp.use(cors());
mlApp.use(express.json());

mlApp.use((req, res, next) => {
  console.log(`[ML Service] ${req.method} ${req.url}`);
  next();
});

mlApp.get("/health", (req, res) => res.status(200).json({ status: "healthy" }));
mlApp.get("/ready", (req, res) => res.status(200).json({ status: "ready" }));

mlApp.post("/api/v1/predict", (req, res) => {
  res.status(200).json({
    prediction: "safe",
    probability: 0.02,
    features: req.body,
  });
});
mlApp.post("/api/v1/feedback", (req, res) => {
  res.status(200).json({ message: "Feedback received" });
});
mlApp.get("/api/v1/feedback/stats", (req, res) => {
  res.status(200).json({ total_feedback: 50, accuracy: 0.95 });
});

mlApp.listen(8000, "0.0.0.0", () => {
  console.log("🤖 ML Service Mock running on port 8000");
});
