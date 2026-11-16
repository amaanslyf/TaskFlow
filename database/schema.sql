-- ============================================================
-- TaskFlow Database Schema
-- ============================================================
-- This file creates the database objects required by TaskFlow.
-- Run it ONCE against a freshly created "taskflow" database.
--
--   1. Create the database:
--        createdb taskflow
--
--   2. Apply this schema:
--        psql -d taskflow -f database/schema.sql
-- ============================================================

-- The uuid-ossp extension gives us uuid_generate_v4() so that
-- PostgreSQL can create unique IDs for every new task.
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- ── tasks table ─────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS tasks (
    id          UUID        PRIMARY KEY DEFAULT uuid_generate_v4(),
    title       TEXT        NOT NULL,
    description TEXT        NOT NULL DEFAULT '',
    status      TEXT        NOT NULL DEFAULT 'todo'
                            CHECK (status IN ('todo', 'in_progress', 'completed')),
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Index on status speeds up the filtered queries we will run
-- from the frontend (e.g. "show me only in-progress tasks").
CREATE INDEX IF NOT EXISTS idx_tasks_status ON tasks (status);

