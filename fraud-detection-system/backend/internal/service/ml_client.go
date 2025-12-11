package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/fraud-detection-system/backend/internal/config"
	"github.com/fraud-detection-system/backend/internal/models"
	"github.com/fraud-detection-system/backend/internal/utils"
)

type MLClient struct {
	baseURL    string
	httpClient *http.Client
	config     *config.Config
}

func NewMLClient(cfg *config.Config) *MLClient {
	return &MLClient{
		baseURL: cfg.ML.ServiceURL,
		httpClient: &http.Client{
			Timeout: cfg.ML.InferenceTimeout,
		},
		config: cfg,
	}
}

// Predict sends a prediction request to the ML service
func (c *MLClient) Predict(ctx context.Context, req *models.MLInferenceRequest) (*models.MLInferenceResponse, error) {
	url := fmt.Sprintf("%s/api/v1/predict", c.baseURL)

	// Marshal request
	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request
	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")

	// Send request
	startTime := time.Now()
	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		utils.GetLogger().WithError(err).Error("Failed to call ML service")
		return nil, fmt.Errorf("failed to call ML service: %w", err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		utils.GetLogger().WithField("status", resp.StatusCode).WithField("body", string(body)).Error("ML service returned error")
		return nil, fmt.Errorf("ML service returned status %d", resp.StatusCode)
	}

	// Parse response
	var mlResp models.MLInferenceResponse
	if err := json.NewDecoder(resp.Body).Decode(&mlResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Add inference time if not set
	if mlResp.InferenceTimeMs == 0 {
		mlResp.InferenceTimeMs = int(time.Since(startTime).Milliseconds())
	}

	utils.GetLogger().WithField("inference_time_ms", mlResp.InferenceTimeMs).Info("ML prediction completed")

	return &mlResp, nil
}

// HealthCheck checks if the ML service is healthy
func (c *MLClient) HealthCheck(ctx context.Context) error {
	url := fmt.Sprintf("%s/health", c.baseURL)

	httpReq, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("failed to call ML service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("ML service returned status %d", resp.StatusCode)
	}

	return nil
}

