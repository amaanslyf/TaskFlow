import uuid
import pytest


# ============================================================================
# B. Task Retrieval Tests
# ============================================================================

def test_get_all_tasks(client, task_factory):
    """
    Verify GET /api/v1/tasks returns a 200 OK and a JSON array of tasks.
    """
    # Create at least one task to guarantee presence
    task = task_factory(title=f"Retrieval Test Task {uuid.uuid4().hex[:6]}")

    response = client.get("/api/v1/tasks")
    assert response.status_code == 200, f"Expected 200, got {response.status_code}: {response.text}"
    tasks = response.json()
    assert isinstance(tasks, list), f"Expected list of tasks, got {type(tasks)}"

    # Confirm our created task is present in the list
    matching = [t for t in tasks if t["id"] == task["id"]]
    assert len(matching) == 1, f"Created task {task['id']} not found in task list"
    assert matching[0]["title"] == task["title"]


def test_get_task_by_id_exists(client, task_factory):
    """
    Verify GET /api/v1/tasks/{id} returns the specific task details.
    """
    unique_title = f"Individual Task {uuid.uuid4().hex[:6]}"
    unique_desc = "Specific task description for ID query test"
    created = task_factory(title=unique_title, description=unique_desc, status="todo")

    response = client.get(f"/api/v1/tasks/{created['id']}")
    assert response.status_code == 200, f"Expected 200, got {response.status_code}: {response.text}"
    fetched = response.json()

    assert fetched["id"] == created["id"]
    assert fetched["title"] == unique_title
    assert fetched["description"] == unique_desc
    assert fetched["status"] == "todo"
    assert "created_at" in fetched
    assert "updated_at" in fetched


def test_get_task_nonexistent_uuid(client):
    """
    Verify GET /api/v1/tasks/{id} returns 404 Not Found when a valid UUID
    does not correspond to any database record.
    """
    nonexistent_id = str(uuid.uuid4())
    response = client.get(f"/api/v1/tasks/{nonexistent_id}")
    assert response.status_code == 404, f"Expected 404 Not Found, got {response.status_code}: {response.text}"
    data = response.json()
    assert data.get("error") == "task not found"


# ============================================================================
# C. Task Creation Tests
# ============================================================================

def test_create_valid_task(client, task_factory):
    """
    Verify POST /api/v1/tasks creates a task and returns 201 Created
    with the expected fields, and that the task can be subsequently retrieved.
    """
    unique_title = f"New Task {uuid.uuid4().hex[:6]}"
    desc = "Created via automated test"
    status = "todo"

    payload = {
        "title": unique_title,
        "description": desc,
        "status": status,
    }
    response = client.post("/api/v1/tasks", json=payload)
    assert response.status_code == 201, f"Expected 201 Created, got {response.status_code}: {response.text}"

    task = response.json()
    assert "id" in task, "Response missing task ID"
    assert task["title"] == unique_title
    assert task["description"] == desc
    assert task["status"] == status
    assert "created_at" in task
    assert "updated_at" in task

    # Register for cleanup
    task_id = task["id"]
    try:
        # Verify subsequent retrieval
        get_res = client.get(f"/api/v1/tasks/{task_id}")
        assert get_res.status_code == 200
        retrieved = get_res.json()
        assert retrieved["id"] == task_id
        assert retrieved["title"] == unique_title
    finally:
        client.delete(f"/api/v1/tasks/{task_id}")


def test_create_task_default_status(client):
    """
    Verify creating a task without status defaults to 'todo'.
    """
    title = f"Default Status Task {uuid.uuid4().hex[:6]}"
    response = client.post("/api/v1/tasks", json={"title": title})
    assert response.status_code == 201
    task = response.json()
    try:
        assert task["status"] == "todo", f"Expected default status 'todo', got {task['status']}"
    finally:
        client.delete(f"/api/v1/tasks/{task['id']}")


# ============================================================================
# D. Task Updating Tests
# ============================================================================

def test_update_task_title_and_status(client, task_factory):
    """
    Verify PUT /api/v1/tasks/{id} updates title and status (to in_progress and completed),
    and that changes persist in the database.
    """
    task = task_factory(title=f"Initial Title {uuid.uuid4().hex[:6]}", status="todo")
    task_id = task["id"]

    # 1. Update title
    updated_title = f"Updated Title {uuid.uuid4().hex[:6]}"
    res_title = client.put(f"/api/v1/tasks/{task_id}", json={"title": updated_title})
    assert res_title.status_code == 200, f"Expected 200, got {res_title.status_code}: {res_title.text}"
    assert res_title.json()["title"] == updated_title

    # 2. Update status to in_progress
    res_status = client.put(f"/api/v1/tasks/{task_id}", json={"status": "in_progress"})
    assert res_status.status_code == 200
    assert res_status.json()["status"] == "in_progress"

    # 3. Verify changes persist via GET
    get_res = client.get(f"/api/v1/tasks/{task_id}")
    assert get_res.status_code == 200
    persisted = get_res.json()
    assert persisted["title"] == updated_title
    assert persisted["status"] == "in_progress"

    # 4. Update status to completed
    res_completed = client.put(f"/api/v1/tasks/{task_id}", json={"status": "completed"})
    assert res_completed.status_code == 200
    assert res_completed.json()["status"] == "completed"


def test_update_nonexistent_task(client):
    """
    Verify PUT /api/v1/tasks/{id} returns 404 for a valid UUID not in database.
    """
    nonexistent_id = str(uuid.uuid4())
    response = client.put(f"/api/v1/tasks/{nonexistent_id}", json={"title": "New Title"})
    assert response.status_code == 404
    assert response.json().get("error") == "task not found"


# ============================================================================
# E. Task Deletion Tests
# ============================================================================

def test_delete_task(client, task_factory):
    """
    Verify DELETE /api/v1/tasks/{id} deletes the record and subsequent GET returns 404.
    """
    task = task_factory(title=f"Task to Delete {uuid.uuid4().hex[:6]}")
    task_id = task["id"]

    # Delete task
    delete_res = client.delete(f"/api/v1/tasks/{task_id}")
    assert delete_res.status_code == 200, f"Expected 200, got {delete_res.status_code}: {delete_res.text}"
    assert delete_res.json().get("message") == "task deleted"

    # Verify subsequent GET returns 404
    get_res = client.get(f"/api/v1/tasks/{task_id}")
    assert get_res.status_code == 404, f"Expected 404 after deletion, got {get_res.status_code}"
    assert get_res.json().get("error") == "task not found"


def test_delete_nonexistent_task(client):
    """
    Verify DELETE /api/v1/tasks/{id} returns 404 for a valid UUID that does not exist.
    """
    nonexistent_id = str(uuid.uuid4())
    response = client.delete(f"/api/v1/tasks/{nonexistent_id}")
    assert response.status_code == 404
    assert response.json().get("error") == "task not found"


# ============================================================================
# F. Input Validation and Error Handling Tests
# ============================================================================

def test_create_task_empty_title_rejected(client):
    """
    Verify POST /api/v1/tasks rejects empty or whitespace-only titles with 400 Bad Request.
    """
    # Empty string
    res1 = client.post("/api/v1/tasks", json={"title": ""})
    assert res1.status_code == 400
    assert "title is required" in res1.json().get("error", "")

    # Whitespace only
    res2 = client.post("/api/v1/tasks", json={"title": "   "})
    assert res2.status_code == 400
    assert "title is required" in res2.json().get("error", "")


def test_create_task_invalid_status_rejected(client):
    """
    Verify POST /api/v1/tasks rejects statuses outside (todo, in_progress, completed).
    """
    payload = {
        "title": "Invalid Status Task",
        "status": "pending_approval",
    }
    response = client.post("/api/v1/tasks", json=payload)
    assert response.status_code == 400
    assert "status must be one of: todo, in_progress, completed" in response.json().get("error", "")


def test_update_task_empty_title_rejected(client, task_factory):
    """
    Verify PUT /api/v1/tasks/{id} rejects updating title to empty string.
    """
    task = task_factory()
    response = client.put(f"/api/v1/tasks/{task['id']}", json={"title": "  "})
    assert response.status_code == 400
    assert "title cannot be empty" in response.json().get("error", "")


def test_update_task_invalid_status_rejected(client, task_factory):
    """
    Verify PUT /api/v1/tasks/{id} rejects invalid status updates.
    """
    task = task_factory()
    response = client.put(f"/api/v1/tasks/{task['id']}", json={"status": "unknown_status"})
    assert response.status_code == 400
    assert "status must be one of: todo, in_progress, completed" in response.json().get("error", "")


def test_malformed_json_body_rejected(client, task_factory):
    """
    Verify sending unparseable JSON returns 400 Bad Request.
    """
    headers = {"Content-Type": "application/json"}
    malformed_body = "{invalid: json, not_quoted}"

    # POST
    res_post = client.post("/api/v1/tasks", data=malformed_body, headers=headers)
    assert res_post.status_code == 400
    assert "invalid JSON body" in res_post.json().get("error", "")

    # PUT
    task = task_factory()
    res_put = client.put(f"/api/v1/tasks/{task['id']}", data=malformed_body, headers=headers)
    assert res_put.status_code == 400
    assert "invalid JSON body" in res_put.json().get("error", "")


def test_invalid_uuid_rejected(client):
    """
    Verify invalid UUID path parameters return 400 Bad Request.
    """
    invalid_id = "not-a-valid-uuid-12345"

    # GET
    res_get = client.get(f"/api/v1/tasks/{invalid_id}")
    assert res_get.status_code == 400
    assert "invalid or missing task ID" in res_get.json().get("error", "")

    # PUT
    res_put = client.put(f"/api/v1/tasks/{invalid_id}", json={"title": "Valid Title"})
    assert res_put.status_code == 400
    assert "invalid or missing task ID" in res_put.json().get("error", "")

    # DELETE
    res_del = client.delete(f"/api/v1/tasks/{invalid_id}")
    assert res_del.status_code == 400
    assert "invalid or missing task ID" in res_del.json().get("error", "")

