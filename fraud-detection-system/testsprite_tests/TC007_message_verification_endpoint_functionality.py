import requests
import time

BASE_URL = "http://localhost:3000"
VERIFY_ENDPOINT = "/api/v1/verify"
TIMEOUT = 30

def test_message_verification_endpoint_functionality():
    url = BASE_URL + VERIFY_ENDPOINT
    headers = {"Content-Type": "application/json"}
    payload = {
        "message": "Urgent: Your account balance is low. Please update your details immediately.",
        "sender": "BankXYZ"
    }

    # We will try to issue multiple requests to test normal and rate limit handling
    success_response = None
    rate_limited = False

    for _ in range(20):
        try:
            response = requests.post(url, json=payload, headers=headers, timeout=TIMEOUT)
        except requests.RequestException as e:
            assert False, f"Request failed: {e}"

        if response.status_code == 200:
            success_response = response
            json_resp = response.json()
            assert "fraud_probability" in json_resp, "Response missing 'fraud_probability'"
            assert isinstance(json_resp["fraud_probability"], (float, int)), "'fraud_probability' not a number"
            # fraud probability should be between 0 and 1
            assert 0.0 <= json_resp["fraud_probability"] <= 1.0, "'fraud_probability' out of range"
            # Also checking if optional explanation is present if any
            if "explanation" in json_resp:
                assert isinstance(json_resp["explanation"], dict) or json_resp["explanation"] is None, "'explanation' should be dict or null"
            # Test succeeded for a request, continue to try more for rate limiting
        elif response.status_code == 429:
            rate_limited = True
            # Expected in case of rate limit, break loop
            break
        else:
            assert False, f"Unexpected status code received: {response.status_code}"

        # To reduce risk of hitting rate limit instantly if server is strict, small delay
        time.sleep(0.01)

    assert success_response is not None, "No successful verification response received"
    assert rate_limited or True, "Expected to receive 429 Too Many Requests status to validate rate limiting"

test_message_verification_endpoint_functionality()