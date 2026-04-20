package postgres

import (
    "context"
    "errors"
    "test-task/internal/domain/task"

    "github.com/jackc/pgx/v5"
    "github.com/jackc/pgx/v5/pgxpool"
)

type TaskRepository struct {
    pool *pgxpool.Pool
}

func NewTaskRepository(pool *pgxpool.Pool) *TaskRepository {
    return &TaskRepository{pool: pool}
}

func (r *TaskRepository) Create(ctx context.Context, t *task.Task) (*task.Task, error) {
    query := `
        INSERT INTO tasks (
            title, description, status, scheduled_at,
            recurrence_type, recurrence_interval, recurrence_day_of_month,
            recurrence_parity, recurrence_specific_dates, recurrence_end_date,
            parent_task_id
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
        RETURNING id, created_at, updated_at
    `

    err := r.pool.QueryRow(ctx, query,
        t.Title, t.Description, t.Status, t.ScheduledAt,
        t.RecurrenceType, t.RecurrenceInterval, t.RecurrenceDayOfMonth,
        t.RecurrenceParity, t.RecurrenceSpecificDates, t.RecurrenceEndDate,
        t.ParentTaskID,
    ).Scan(&t.ID, &t.CreatedAt, &t.UpdatedAt)

    if err != nil {
        return nil, err
    }

    return t, nil
}

func (r *TaskRepository) GetByID(ctx context.Context, id int) (*task.Task, error) {
    query := `
        SELECT id, title, description, status, scheduled_at, created_at, updated_at,
               recurrence_type, recurrence_interval, recurrence_day_of_month,
               recurrence_parity, recurrence_specific_dates, recurrence_end_date,
               parent_task_id
        FROM tasks
        WHERE id = $1
    `

    t := &task.Task{}
    err := r.pool.QueryRow(ctx, query, id).Scan(
        &t.ID, &t.Title, &t.Description, &t.Status, &t.ScheduledAt,
        &t.CreatedAt, &t.UpdatedAt,
        &t.RecurrenceType, &t.RecurrenceInterval, &t.RecurrenceDayOfMonth,
        &t.RecurrenceParity, &t.RecurrenceSpecificDates, &t.RecurrenceEndDate,
        &t.ParentTaskID,
    )

    if err != nil {
        if errors.Is(err, pgx.ErrNoRows) {
            return nil, task.ErrTaskNotFound
        }
        return nil, err
    }

    return t, nil
}

func (r *TaskRepository) List(ctx context.Context, limit, offset int) ([]*task.Task, error) {
    query := `
        SELECT id, title, description, status, scheduled_at, created_at, updated_at,
               recurrence_type, recurrence_interval, recurrence_day_of_month,
               recurrence_parity, recurrence_specific_dates, recurrence_end_date,
               parent_task_id
        FROM tasks
        ORDER BY created_at DESC
        LIMIT $1 OFFSET $2
    `

    rows, err := r.pool.Query(ctx, query, limit, offset)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var tasks []*task.Task
    for rows.Next() {
        t := &task.Task{}
        err := rows.Scan(
            &t.ID, &t.Title, &t.Description, &t.Status, &t.ScheduledAt,
            &t.CreatedAt, &t.UpdatedAt,
            &t.RecurrenceType, &t.RecurrenceInterval, &t.RecurrenceDayOfMonth,
            &t.RecurrenceParity, &t.RecurrenceSpecificDates, &t.RecurrenceEndDate,
            &t.ParentTaskID,
        )
        if err != nil {
            return nil, err
        }
        tasks = append(tasks, t)
    }

    return tasks, nil
}

func (r *TaskRepository) Update(ctx context.Context, t *task.Task) (*task.Task, error) {
    query := `
        UPDATE tasks
        SET title = $1, description = $2, status = $3, scheduled_at = $4,
            recurrence_type = $5, recurrence_interval = $6, recurrence_day_of_month = $7,
            recurrence_parity = $8, recurrence_specific_dates = $9, recurrence_end_date = $10,
            updated_at = CURRENT_TIMESTAMP
        WHERE id = $11
        RETURNING updated_at
    `

    err := r.pool.QueryRow(ctx, query,
        t.Title, t.Description, t.Status, t.ScheduledAt,
        t.RecurrenceType, t.RecurrenceInterval, t.RecurrenceDayOfMonth,
        t.RecurrenceParity, t.RecurrenceSpecificDates, t.RecurrenceEndDate,
        t.ID,
    ).Scan(&t.UpdatedAt)

    if err != nil {
        if errors.Is(err, pgx.ErrNoRows) {
            return nil, task.ErrTaskNotFound
        }
        return nil, err
    }

    return t, nil
}

func (r *TaskRepository) Delete(ctx context.Context, id int) error {
    query := `DELETE FROM tasks WHERE id = $1`

    result, err := r.pool.Exec(ctx, query, id)
    if err != nil {
        return err
    }

    if result.RowsAffected() == 0 {
        return task.ErrTaskNotFound
    }

    return nil
}

func (r *TaskRepository) GetRecurringTasks(ctx context.Context) ([]*task.Task, error) {
    query := `
        SELECT id, title, description, status, scheduled_at, created_at, updated_at,
               recurrence_type, recurrence_interval, recurrence_day_of_month,
               recurrence_parity, recurrence_specific_dates, recurrence_end_date,
               parent_task_id
        FROM tasks
        WHERE recurrence_type IS NOT NULL
          AND recurrence_type != ''
          AND parent_task_id IS NULL
    `

    rows, err := r.pool.Query(ctx, query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var tasks []*task.Task
    for rows.Next() {
        t := &task.Task{}
        err := rows.Scan(
            &t.ID, &t.Title, &t.Description, &t.Status, &t.ScheduledAt,
            &t.CreatedAt, &t.UpdatedAt,
            &t.RecurrenceType, &t.RecurrenceInterval, &t.RecurrenceDayOfMonth,
            &t.RecurrenceParity, &t.RecurrenceSpecificDates, &t.RecurrenceEndDate,
            &t.ParentTaskID,
        )
        if err != nil {
            return nil, err
        }
        tasks = append(tasks, t)
    }

    return tasks, nil
}