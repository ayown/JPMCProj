import requests

BASE_URL = "http://localhost:3000"
LOGIN_ENDPOINT = f"{BASE_URL}/api/v1/auth/login"
TIMEOUT = 30

def test_user_login_endpoint_authentication():
    # Valid credentials for successful login test
    valid_credentials = {
        "username": "testuser",
        "password": "testpassword"
    }
    # Invalid credentials for unauthorized test
    invalid_credentials = {
        "username": "invaliduser",
        "password": "wrongpassword"
    }

    # Test successful login with valid credentials
    try:
        response = requests.post(LOGIN_ENDPOINT, json=valid_credentials, timeout=TIMEOUT)
    except requests.RequestException as e:
        assert False, f"Request failed: {e}"

    assert response.status_code == 200, f"Expected status 200, got {response.status_code}"
    try:
        json_data = response.json()
    except ValueError:
        assert False, "Response is not valid JSON"
    # Expecting a token or similar field in success response, check presence of some key like 'token'
    assert "token" in json_data or "access_token" in json_data, "JWT token not found in response"

    # Test unauthorized access with invalid credentials
    try:
        response_unauth = requests.post(LOGIN_ENDPOINT, json=invalid_credentials, timeout=TIMEOUT)
    except requests.RequestException as e:
        assert False, f"Request failed: {e}"

    assert response_unauth.status_code == 401, f"Expected status 401, got {response_unauth.status_code}"

test_user_login_endpoint_authentication()