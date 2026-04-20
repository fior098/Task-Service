CREATE TABLE IF NOT EXISTS tasks (
                                     id SERIAL PRIMARY KEY,
                                     title VARCHAR(255) NOT NULL,
    description TEXT,
    status VARCHAR(50) NOT NULL DEFAULT 'new',
    scheduled_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    recurrence_type VARCHAR(50),
    recurrence_interval INTEGER,
    recurrence_day_of_month INTEGER,
    recurrence_parity VARCHAR(10),
    recurrence_specific_dates TEXT,
    recurrence_end_date TIMESTAMP,
    parent_task_id INTEGER REFERENCES tasks(id) ON DELETE SET NULL
    );

CREATE INDEX idx_tasks_scheduled_at ON tasks(scheduled_at);
CREATE INDEX idx_tasks_status ON tasks(status);
CREATE INDEX idx_tasks_parent_task_id ON tasks(parent_task_id);
CREATE INDEX idx_tasks_recurrence_type ON tasks(recurrence_type);