import requests
from requests.exceptions import RequestException

BASE_URL = "http://localhost:3000"

def test_token_refresh_endpoint_functionality():
    headers = {"Content-Type": "application/json"}
    timeout = 30

    # Step 1: Register a new user to get credentials for login
    register_payload = {
        "username": "testuser_refresh",
        "password": "TestPassword123!"
    }
    try:
        register_resp = requests.post(
            f"{BASE_URL}/api/v1/auth/register",
            json=register_payload,
            headers=headers,
            timeout=timeout,
        )
        assert register_resp.status_code == 201, f"Registration failed: {register_resp.text}"

        # Step 2: Login to obtain access and refresh tokens
        login_payload = {
            "username": register_payload["username"],
            "password": register_payload["password"]
        }
        login_resp = requests.post(
            f"{BASE_URL}/api/v1/auth/login",
            json=login_payload,
            headers=headers,
            timeout=timeout,
        )
        assert login_resp.status_code == 200, f"Login failed: {login_resp.text}"
        login_data = login_resp.json()
        assert "access_token" in login_data, "No access_token in login response"
        assert "refresh_token" in login_data, "No refresh_token in login response"
        refresh_token = login_data["refresh_token"]

        # Step 3: Refresh token with valid refresh token
        refresh_headers = {"Authorization": f"Bearer {refresh_token}"}
        refresh_resp = requests.post(
            f"{BASE_URL}/api/v1/auth/refresh",
            headers=refresh_headers,
            timeout=timeout,
        )
        assert refresh_resp.status_code == 200, f"Token refresh failed: {refresh_resp.text}"
        refresh_data = refresh_resp.json()
        assert "access_token" in refresh_data, "No access_token in refresh response"

        # Step 4: Refresh token with invalid/expired token should return 401
        invalid_headers = {"Authorization": "Bearer invalid_or_expired_token"}
        invalid_resp = requests.post(
            f"{BASE_URL}/api/v1/auth/refresh",
            headers=invalid_headers,
            timeout=timeout,
        )
        assert invalid_resp.status_code == 401, f"Unauthorized refresh did not return 401: {invalid_resp.text}"

    finally:
        # Cleanup: No explicit delete user endpoint given in PRD, so no cleanup possible here.
        pass

test_token_refresh_endpoint_functionality()