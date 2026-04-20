package task

import (
    "context"
    "test-task/internal/domain/task"
)

type Repository interface {
    Create(ctx context.Context, task *task.Task) (*task.Task, error)
    GetByID(ctx context.Context, id int) (*task.Task, error)
    List(ctx context.Context, limit, offset int) ([]*task.Task, error)
    Update(ctx context.Context, task *task.Task) (*task.Task, error)
    Delete(ctx context.Context, id int) error
    GetRecurringTasks(ctx context.Context) ([]*task.Task, error)
}

type Scheduler interface {
    ScheduleRecurringTasks(ctx context.Context) error
}