import requests

def test_health_status_endpoint_verification():
    url = "http://localhost:3000/health"
    headers = {
        "Accept": "application/json"
    }
    try:
        response = requests.get(url, headers=headers, timeout=30)
        response.raise_for_status()
    except requests.RequestException as e:
        assert False, f"Request to /health failed: {e}"

    assert response.status_code == 200, f"Expected status code 200, got {response.status_code}"

test_health_status_endpoint_verification()