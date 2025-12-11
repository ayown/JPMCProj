import requests
import time

BASE_URL = "http://localhost:3000"
PREDICT_ENDPOINT = "/api/v1/predict"
TIMEOUT = 30  # seconds
MAX_LATENCY_MS = 500  # max allowed latency in milliseconds

def test_ml_prediction_endpoint_accuracy_and_performance():
    url = BASE_URL + PREDICT_ENDPOINT
    headers = {
        "Content-Type": "application/json"
    }
    # Example payload simulating a typical message for fraud prediction
    payload = {
        "message": "URGENT: Your account has been compromised. Please verify immediately.",
        "sender": "+1234567890"
    }
    try:
        start_time = time.time()
        response = requests.post(url, json=payload, headers=headers, timeout=TIMEOUT)
        elapsed_ms = (time.time() - start_time) * 1000
    except requests.RequestException as e:
        assert False, f"Request to ML predict endpoint failed: {e}"

    # Validate response status code
    assert response.status_code == 200, f"Expected status code 200, got {response.status_code}"

    try:
        result = response.json()
    except ValueError:
        assert False, "Response is not valid JSON"

    # Validate that response contains the expected keys for fraud detection
    # According to PRD: fraud probability scores and explanations expected
    assert "fraud_probability" in result, "Response missing 'fraud_probability'"
    assert isinstance(result["fraud_probability"], (float, int)), "'fraud_probability' must be a number"
    assert 0.0 <= result["fraud_probability"] <= 1.0, "'fraud_probability' must be between 0 and 1"

    assert "explanation" in result, "Response missing 'explanation'"
    assert isinstance(result["explanation"], dict), "'explanation' must be a dictionary with details"

    # Validate latency constraints (< 500 ms total response)
    assert elapsed_ms < MAX_LATENCY_MS, f"Response latency {elapsed_ms:.1f}ms exceeds max allowed {MAX_LATENCY_MS}ms"

    # Optionally, validate sensible fraud probability range (e.g., precision >95%, but since no ground truth, just check range)
    # Test case does not specify a ground truth label to verify accuracy, so only basic validation performed here

test_ml_prediction_endpoint_accuracy_and_performance()
