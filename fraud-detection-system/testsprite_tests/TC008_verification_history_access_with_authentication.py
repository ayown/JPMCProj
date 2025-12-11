import requests

BASE_URL = "http://localhost:3000"
VERIFY_HISTORY_ENDPOINT = "/api/v1/verify/history"
AUTH_LOGIN_ENDPOINT = "/api/v1/auth/login"

TIMEOUT = 30

def verification_history_access_with_authentication():
    # Sample valid user credentials for login - these should be valid in the test environment
    login_payload = {
        "username": "testuser",
        "password": "testpassword"
    }

    # Step 1: Obtain valid JWT token by logging in
    try:
        login_resp = requests.post(
            BASE_URL + AUTH_LOGIN_ENDPOINT,
            json=login_payload,
            timeout=TIMEOUT
        )
    except requests.RequestException as e:
        assert False, f"Login request failed: {e}"
    assert login_resp.status_code == 200, f"Login failed with status code {login_resp.status_code}"
    login_json = login_resp.json()
    assert "token" in login_json or "access_token" in login_json, "Login response missing JWT token"
    token = login_json.get("token") or login_json.get("access_token")
    assert isinstance(token, str) and token != "", "JWT token is empty or invalid"

    headers_auth = {
        "Authorization": f"Bearer {token}"
    }

    # Step 2: Access /api/v1/verify/history with valid JWT token
    try:
        resp_auth = requests.get(
            BASE_URL + VERIFY_HISTORY_ENDPOINT,
            headers=headers_auth,
            timeout=TIMEOUT
        )
    except requests.RequestException as e:
        assert False, f"Authorized history request failed: {e}"
    assert resp_auth.status_code == 200, f"Expected 200 OK with valid token but got {resp_auth.status_code}"
    # Optionally check response json structure (should be list or dict)
    try:
        data = resp_auth.json()
        assert isinstance(data, (dict, list)), "Response JSON is not dict or list"
    except ValueError:
        assert False, "Response is not valid JSON"

    # Step 3: Access /api/v1/verify/history without JWT token to verify 401 Unauthorized
    try:
        resp_unauth = requests.get(
            BASE_URL + VERIFY_HISTORY_ENDPOINT,
            timeout=TIMEOUT
        )
    except requests.RequestException as e:
        assert False, f"Unauthorized history request failed: {e}"
    assert resp_unauth.status_code == 401, f"Expected 401 Unauthorized without token but got {resp_unauth.status_code}"

verification_history_access_with_authentication()