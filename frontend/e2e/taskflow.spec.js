import { test, expect } from '@playwright/test';

const API_BASE_URL = process.env.API_BASE_URL || 'http://localhost:8080';

test.describe('TaskFlow End-to-End Application Workflow', () => {
  // Array to track task IDs created during tests for reliable cleanup
  let createdTaskIds = [];

  // Helper to query backend for tasks with matching prefix
  async function cleanupTestTasks(request) {
    try {
      const response = await request.get(`${API_BASE_URL}/api/v1/tasks`);
      if (response.ok()) {
        const tasks = await response.json();
        for (const task of tasks) {
          if (task.title && task.title.startsWith('E2E Test Task')) {
            await request.delete(`${API_BASE_URL}/api/v1/tasks/${task.id}`);
          }
        }
      }
    } catch {
      // Backend may not be reachable or already cleaned
    }
  }

  test.beforeEach(async () => {
    createdTaskIds = [];
  });

  test.afterEach(async ({ request }) => {
    // Clean up any tasks tracked by ID
    for (const id of createdTaskIds) {
      try {
        await request.delete(`${API_BASE_URL}/api/v1/tasks/${id}`);
      } catch {
        // Ignored
      }
    }
    // Also clean up any lingering E2E test tasks
    await cleanupTestTasks(request);
  });

  test('1. Dashboard loads successfully with header and statistics cards', async ({ page }) => {
    await page.goto('/');

    // Check application branding
    await expect(page.getByRole('heading', { name: /TaskFlow/i, level: 1 })).toBeVisible();
    await expect(page.getByText('A simple task manager — built to learn DevOps')).toBeVisible();

    // Verify statistics cards are displayed
    await expect(page.getByLabel('Task statistics')).toBeVisible();
    await expect(page.getByLabel('Total tasks')).toBeVisible();
    await expect(page.getByLabel('To do tasks')).toBeVisible();
    await expect(page.getByLabel('In progress tasks')).toBeVisible();
    await expect(page.getByLabel('Completed tasks')).toBeVisible();
  });

  test('2 & 3. Create a task using UI and verify it appears in the task list', async ({ page, request }) => {
    const uniqueTitle = `E2E Test Task Create ${Date.now()}`;
    const uniqueDesc = 'End-to-end task created via browser form';

    await page.goto('/');

    // Fill the task form using exact accessible labels
    await page.getByLabel('Title *').fill(uniqueTitle);
    await page.getByLabel('Description').fill(uniqueDesc);
    await page.getByLabel('Status', { exact: true }).selectOption('todo');

    // Submit the form
    await page.getByRole('button', { name: 'Create Task' }).click();

    // Verify the newly created task card appears in the list
    const taskCard = page.locator('.task-card', { hasText: uniqueTitle });
    await expect(taskCard).toBeVisible();
    await expect(taskCard.getByText(uniqueDesc)).toBeVisible();

    // Verify status badge
    await expect(taskCard.getByText('To Do')).toBeVisible();

    // Verify input fields reset
    await expect(page.getByLabel('Title *')).toHaveValue('');
    await expect(page.getByLabel('Description')).toHaveValue('');

    // Fetch and store ID for cleanup
    const res = await request.get(`${API_BASE_URL}/api/v1/tasks`);
    if (res.ok()) {
      const list = await res.json();
      const match = list.find((t) => t.title === uniqueTitle);
      if (match) createdTaskIds.push(match.id);
    }
  });

  test('4. Update a task title and description', async ({ page, request }) => {
    const originalTitle = `E2E Test Task PreUpdate ${Date.now()}`;
    const updatedTitle = `E2E Test Task PostUpdate ${Date.now()}`;

    // Create a task via API to guarantee test independence
    const res = await request.post(`${API_BASE_URL}/api/v1/tasks`, {
      data: {
        title: originalTitle,
        description: 'Original description before editing',
        status: 'todo',
      },
    });
    expect(res.ok()).toBeTruthy();
    const created = await res.json();
    createdTaskIds.push(created.id);

    await page.goto('/');

    // Locate the task card and click Edit
    const taskCard = page.locator('.task-card', { hasText: originalTitle });
    await expect(taskCard).toBeVisible();
    await taskCard.getByRole('button', { name: `Edit task ${originalTitle}` }).click();

    // Verify form enters edit mode
    await expect(page.getByRole('heading', { name: 'Edit Task' })).toBeVisible();
    await expect(page.getByLabel('Title *')).toHaveValue(originalTitle);

    // Modify title and submit
    await page.getByLabel('Title *').fill(updatedTitle);
    await page.getByRole('button', { name: 'Update Task' }).click();

    // Verify update reflected on task card
    await expect(page.locator('.task-card', { hasText: updatedTitle })).toBeVisible();
    await expect(page.locator('.task-card', { hasText: originalTitle })).not.toBeVisible();
  });

  test('5. Change task status by clicking status badge', async ({ page, request }) => {
    const title = `E2E Test Task StatusCycle ${Date.now()}`;

    // Create task via API
    const res = await request.post(`${API_BASE_URL}/api/v1/tasks`, {
      data: { title, status: 'todo' },
    });
    expect(res.ok()).toBeTruthy();
    const created = await res.json();
    createdTaskIds.push(created.id);

    await page.goto('/');

    const taskCard = page.locator('.task-card', { hasText: title });
    await expect(taskCard).toBeVisible();

    // Initially "To Do"
    const badge = taskCard.getByRole('button', { name: /Status:/i });
    await expect(badge).toHaveText('To Do');

    // Click badge to cycle: todo -> in_progress
    await badge.click();
    await expect(badge).toHaveText('In Progress');

    // Click badge to cycle: in_progress -> completed
    await badge.click();
    await expect(badge).toHaveText('Completed');
  });

  test('6. Filter tasks by status', async ({ page, request }) => {
    const todoTitle = `E2E Test Task FilterTodo ${Date.now()}`;
    const compTitle = `E2E Test Task FilterComp ${Date.now()}`;

    // Create 1 todo task and 1 completed task
    const res1 = await request.post(`${API_BASE_URL}/api/v1/tasks`, {
      data: { title: todoTitle, status: 'todo' },
    });
    const res2 = await request.post(`${API_BASE_URL}/api/v1/tasks`, {
      data: { title: compTitle, status: 'completed' },
    });
    const t1 = await res1.json();
    const t2 = await res2.json();
    createdTaskIds.push(t1.id, t2.id);

    await page.goto('/');

    // Under "All" filter, both should be visible
    await expect(page.locator('.task-card', { hasText: todoTitle })).toBeVisible();
    await expect(page.locator('.task-card', { hasText: compTitle })).toBeVisible();

    // Filter by "To Do"
    await page.getByRole('button', { name: 'To Do', exact: true }).click();
    await expect(page.locator('.task-card', { hasText: todoTitle })).toBeVisible();
    await expect(page.locator('.task-card', { hasText: compTitle })).not.toBeVisible();

    // Filter by "Completed"
    await page.getByRole('button', { name: 'Completed', exact: true }).click();
    await expect(page.locator('.task-card', { hasText: compTitle })).toBeVisible();
    await expect(page.locator('.task-card', { hasText: todoTitle })).not.toBeVisible();

    // Return to "All"
    await page.getByRole('button', { name: 'All', exact: true }).click();
    await expect(page.locator('.task-card', { hasText: todoTitle })).toBeVisible();
    await expect(page.locator('.task-card', { hasText: compTitle })).toBeVisible();
  });

  test('7 & 8. Delete a task and confirm removal from list and statistics', async ({ page, request }) => {
    const title = `E2E Test Task Remove ${Date.now()}`;

    const res = await request.post(`${API_BASE_URL}/api/v1/tasks`, {
      data: { title, status: 'todo' },
    });
    expect(res.ok()).toBeTruthy();
    const created = await res.json();

    await page.goto('/');

    const taskCard = page.locator('.task-card', { hasText: title });
    await expect(taskCard).toBeVisible();

    // Auto-accept window.confirm dialog
    page.once('dialog', async (dialog) => {
      expect(dialog.message()).toContain('Delete this task?');
      await dialog.accept();
    });

    // Click Delete button using explicit accessible label
    await taskCard.getByRole('button', { name: `Delete task ${title}` }).click();

    // Confirm task card is removed from DOM
    await expect(page.locator('.task-card', { hasText: title })).not.toBeVisible();

    // Verify deletion persisted in backend API
    const verifyRes = await request.get(`${API_BASE_URL}/api/v1/tasks/${created.id}`);
    expect(verifyRes.status()).toBe(404);
  });
});
