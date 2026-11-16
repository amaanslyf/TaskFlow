import pytest


def test_healthz_endpoint(client):
    """
    Test GET /healthz.
    Liveness probe: Backend process is alive and responding to HTTP requests.
    """
    response = client.get("/healthz")
    assert response.status_code == 200, f"Expected 200 OK, got {response.status_code}: {response.text}"
    data = response.json()
    assert data.get("status") == "ok", f"Expected {{'status': 'ok'}}, got {data}"


def test_healthz_method_not_allowed(client):
    """
    Test POST /healthz returns 405 Method Not Allowed.
    """
    response = client.post("/healthz", json={})
    assert response.status_code == 405, f"Expected 405, got {response.status_code}"


def test_readyz_endpoint(client):
    """
    Test GET /readyz.
    Readiness probe: Backend can successfully communicate with PostgreSQL database.
    """
    response = client.get("/readyz")
    assert response.status_code == 200, f"Expected 200 OK, got {response.status_code}: {response.text}"
    data = response.json()
    assert data.get("status") == "ready", f"Expected {{'status': 'ready'}}, got {data}"


def test_readyz_method_not_allowed(client):
    """
    Test POST /readyz returns 405 Method Not Allowed.
    """
    response = client.post("/readyz", json={})
    assert response.status_code == 405, f"Expected 405, got {response.status_code}"

