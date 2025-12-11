import requests

BASE_URL = "http://localhost:3000"
REGISTER_URL = f"{BASE_URL}/api/v1/auth/register"
LOGIN_URL = f"{BASE_URL}/api/v1/auth/login"
PROFILE_URL = f"{BASE_URL}/api/v1/profile"
TIMEOUT = 30

def test_user_profile_access_with_authentication():
    # Dummy user data for registration and login
    user_data = {
        "username": "testuser_tc006",
        "email": "testuser_tc006@example.com",
        "password": "StrongPass!123"
    }
    # Register the user
    try:
        register_resp = requests.post(REGISTER_URL, json=user_data, timeout=TIMEOUT)
        # Registration might be 201 Created or 400 if user exists, accept both
        assert register_resp.status_code in (201, 400)

        # Login the user to get JWT token
        login_payload = {
            "username": user_data["username"],
            "password": user_data["password"]
        }
        login_resp = requests.post(LOGIN_URL, json=login_payload, timeout=TIMEOUT)
        assert login_resp.status_code == 200
        token = login_resp.json().get("access_token")
        assert token and isinstance(token, str)

        headers_auth = {"Authorization": f"Bearer {token}"}
        # Access profile endpoint with valid JWT token
        profile_resp = requests.get(PROFILE_URL, headers=headers_auth, timeout=TIMEOUT)
        assert profile_resp.status_code == 200
        profile_data = profile_resp.json()
        assert isinstance(profile_data, dict)
        # Basic sanity check that profile data contains username or email matching user_data
        username_match = profile_data.get("username") == user_data["username"] if "username" in profile_data else False
        email_match = profile_data.get("email") == user_data["email"] if "email" in profile_data else False
        assert username_match or email_match

        # Access profile endpoint without authentication headers
        profile_resp_unauth = requests.get(PROFILE_URL, timeout=TIMEOUT)
        assert profile_resp_unauth.status_code == 401

        # Access profile endpoint with invalid/expired token
        headers_invalid = {"Authorization": "Bearer invalid_token_example"}
        profile_resp_invalid = requests.get(PROFILE_URL, headers=headers_invalid, timeout=TIMEOUT)
        assert profile_resp_invalid.status_code == 401

    finally:
        # Cleanup: No delete user endpoint provided in PRD, so skip user removal
        pass

test_user_profile_access_with_authentication()
