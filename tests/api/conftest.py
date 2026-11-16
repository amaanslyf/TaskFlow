import os
import uuid
import pytest
import requests

DEFAULT_BASE_URL = "http://localhost:8080"


@pytest.fixture(scope="session")
def api_base_url():
    """Returns the base URL for the TaskFlow Go backend API."""
    return os.getenv("API_BASE_URL", DEFAULT_BASE_URL).rstrip("/")


@pytest.fixture(scope="session", autouse=True)
def ensure_backend_healthy(api_base_url):
    """
    Verifies that the Go backend is up and running before executing tests.
    Fails fast with a helpful DevOps message if the server is unreachable.
    """
    health_url = f"{api_base_url}/healthz"
    try:
        response = requests.get(health_url, timeout=3)
        if response.status_code != 200:
            pytest.fail(
                f"Backend at {health_url} returned HTTP {response.status_code}. "
                "Ensure the Go backend is running and healthy."
            )
    except requests.exceptions.RequestException as exc:
        pytest.fail(
            f"Could not connect to Go backend at {health_url}: {exc}. "
            "Prerequisite: Start the Go backend with 'go run ./cmd/api' before running API tests."
        )


class APIClient:
    """Convenience wrapper around requests for TaskFlow API."""

    def __init__(self, base_url: str):
        self.base_url = base_url

    def url(self, path: str) -> str:
        return f"{self.base_url}/{path.lstrip('/')}"

    def get(self, path: str, **kwargs):
        return requests.get(self.url(path), timeout=5, **kwargs)

    def post(self, path: str, json=None, data=None, headers=None, **kwargs):
        return requests.post(self.url(path), json=json, data=data, headers=headers, timeout=5, **kwargs)

    def put(self, path: str, json=None, data=None, headers=None, **kwargs):
        return requests.put(self.url(path), json=json, data=data, headers=headers, timeout=5, **kwargs)

    def delete(self, path: str, **kwargs):
        return requests.delete(self.url(path), timeout=5, **kwargs)


@pytest.fixture
def client(api_base_url):
    """Provides an API client instance."""
    return APIClient(api_base_url)


@pytest.fixture
def task_factory(client):
    """
    Factory fixture for generating test tasks with automatic database cleanup.
    Tracks all created task IDs and deletes them after test completion,
    preventing test pollution and preserving any existing user tasks.
    """
    created_task_ids = []

    def _create_task(title=None, description="Automated API test task", status="todo"):
        if title is None:
            title = f"Test Task {uuid.uuid4().hex[:8]}"

        payload = {
            "title": title,
            "description": description,
            "status": status,
        }
        res = client.post("/api/v1/tasks", json=payload)
        assert res.status_code == 201, f"Failed to create task in factory: {res.text}"
        task = res.json()
        created_task_ids.append(task["id"])
        return task

    yield _create_task

    # Teardown: delete all tasks created by this test
    for task_id in created_task_ids:
        try:
            client.delete(f"/api/v1/tasks/{task_id}")
        except Exception:
            pass

