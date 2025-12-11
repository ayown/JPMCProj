import requests
import uuid

BASE_URL = "http://localhost:3000"
REGISTER_URL = f"{BASE_URL}/api/v1/auth/register"
LOGIN_URL = f"{BASE_URL}/api/v1/auth/login"
REPORTS_URL = f"{BASE_URL}/api/v1/reports"
TIMEOUT = 30


def test_report_creation_and_listing_with_authentication():
    # Register a new user for authentication
    username = f"user_{uuid.uuid4().hex[:8]}"
    password = "StrongPassw0rd!"
    register_payload = {
        "username": username,
        "password": password,
        "email": f"{username}@example.com"
    }
    try:
        register_resp = requests.post(REGISTER_URL, json=register_payload, timeout=TIMEOUT)
        # It might be possible user already exists if run multiple times rapidly; accept 201 Created or 400 Bad Request if already exists
        assert register_resp.status_code in (201, 400), f"Unexpected register status: {register_resp.status_code}"

        # Login to get JWT token
        login_payload = {
            "username": username,
            "password": password
        }
        login_resp = requests.post(LOGIN_URL, json=login_payload, timeout=TIMEOUT)
        assert login_resp.status_code == 200, f"Login failed with status: {login_resp.status_code}"
        login_json = login_resp.json()
        assert "access_token" in login_json, "No access_token in login response"
        token = login_json["access_token"]
        headers = {"Authorization": f"Bearer {token}"}

        # Create a new report with valid authentication
        report_payload = {
            "title": "Suspicious SMS Fraud",
            "message": "Urgent: Please verify your account details immediately!",
            "sender": "BANK123",
            "reported_at": "2025-12-10T10:00:00Z",
            "details": "Contains urgency manipulation and header spoofing characteristics."
        }
        create_resp = requests.post(REPORTS_URL, json=report_payload, headers=headers, timeout=TIMEOUT)
        assert create_resp.status_code == 201, f"Report creation failed with status: {create_resp.status_code}"
        created_report = create_resp.json()
        assert "id" in created_report, "Created report response missing 'id'"
        report_id = created_report["id"]

        # List reports with valid authentication
        list_resp = requests.get(REPORTS_URL, headers=headers, timeout=TIMEOUT)
        assert list_resp.status_code == 200, f"Listing reports failed with status: {list_resp.status_code}"
        reports_list = list_resp.json()
        assert isinstance(reports_list, list), "Reports list response is not a list"
        assert any(r.get("id") == report_id for r in reports_list), "Created report not found in list"

        # Attempt to create a report without authentication (should fail with 401)
        unauthorized_create_resp = requests.post(REPORTS_URL, json=report_payload, timeout=TIMEOUT)
        assert unauthorized_create_resp.status_code == 401, f"Unauthorized create expected 401 but got {unauthorized_create_resp.status_code}"

        # Attempt to list reports without authentication (should fail with 401)
        unauthorized_list_resp = requests.get(REPORTS_URL, timeout=TIMEOUT)
        assert unauthorized_list_resp.status_code == 401, f"Unauthorized list expected 401 but got {unauthorized_list_resp.status_code}"

    finally:
        # Cleanup: delete the created report if exists and authenticated
        if 'report_id' in locals() and 'headers' in locals():
            try:
                delete_resp = requests.delete(f"{REPORTS_URL}/{report_id}", headers=headers, timeout=TIMEOUT)
                # Accept 200 or 204 as successful deletion, or 404 if already deleted
                assert delete_resp.status_code in (200, 204, 404), f"Unexpected delete status: {delete_resp.status_code}"
            except Exception:
                pass


test_report_creation_and_listing_with_authentication()