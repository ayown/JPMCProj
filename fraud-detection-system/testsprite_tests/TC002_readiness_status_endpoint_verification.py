import requests

ENDPOINT = "http://localhost:3000"

def test_readiness_status_endpoint_verification():
    url = f"{ENDPOINT}/ready"
    try:
        response = requests.get(url, timeout=30)
        response.raise_for_status()
        assert response.status_code == 200, f"Expected status code 200, got {response.status_code}"
    except requests.exceptions.RequestException as e:
        assert False, f"Request to {url} failed: {e}"

test_readiness_status_endpoint_verification()