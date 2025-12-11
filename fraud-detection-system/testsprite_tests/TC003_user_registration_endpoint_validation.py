import requests

BASE_URL = "http://localhost:3000"
REGISTER_ENDPOINT = "/api/v1/auth/register"
TIMEOUT = 30

def test_user_registration_endpoint_validation():
    url = BASE_URL + REGISTER_ENDPOINT
    headers = {"Content-Type": "application/json"}

    # Valid registration data
    valid_payload = {
        "username": "testuser_tc003",
        "email": "testuser_tc003@example.com",
        "password": "StrongPassw0rd!"
    }

    # Invalid registration payloads to test 400 Bad Request
    invalid_payloads = [
        {},  # empty body
        {"username": ""},  # missing required fields
        {"email": "not-an-email", "password": "123456"},  # invalid email format
        {"username": "user", "email": "user@example.com"},  # missing password
        {"username": "user", "password": "1234"},  # missing email
        {"username": "user"*1000, "email": "user@example.com", "password": "password"},  # excessively long username
    ]

    try:
        # Test successful registration
        resp = requests.post(url, json=valid_payload, headers=headers, timeout=TIMEOUT)
        assert resp.status_code == 201, f"Expected 201 Created, got {resp.status_code}"

        # Test invalid inputs result in 400 Bad Request
        for payload in invalid_payloads:
            resp = requests.post(url, json=payload, headers=headers, timeout=TIMEOUT)
            assert resp.status_code == 400, f"Expected 400 Bad Request for payload {payload}, got {resp.status_code}"

    except requests.RequestException as e:
        assert False, f"Request failed: {e}"

test_user_registration_endpoint_validation()
