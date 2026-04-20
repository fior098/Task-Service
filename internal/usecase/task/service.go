package task

import (
    "context"
    "test-task/internal/domain/task"
)

type Service struct {
    repo      Repository
    scheduler Scheduler
}

func NewService(repo Repository, scheduler Scheduler) *Service {
    return &Service{
        repo:      repo,
        scheduler: scheduler,
    }
}

func (s *Service) CreateTask(ctx context.Context, t *task.Task) (*task.Task, error) {
    if err := t.Validate(); err != nil {
        return nil, err
    }

    created, err := s.repo.Create(ctx, t)
    if err != nil {
        return nil, err
    }

    if created.IsRecurring() {
        go s.scheduler.ScheduleRecurringTasks(context.Background())
    }

    return created, nil
}

func (s *Service) GetTask(ctx context.Context, id int) (*task.Task, error) {
    return s.repo.GetByID(ctx, id)
}

func (s *Service) ListTasks(ctx context.Context, limit, offset int) ([]*task.Task, error) {
    if limit <= 0 {
        limit = 10
    }
    if offset < 0 {
        offset = 0
    }

    return s.repo.List(ctx, limit, offset)
}

func (s *Service) UpdateTask(ctx context.Context, t *task.Task) (*task.Task, error) {
    if err := t.Validate(); err != nil {
        return nil, err
    }

    existing, err := s.repo.GetByID(ctx, t.ID)
    if err != nil {
        return nil, err
    }
    if existing == nil {
        return nil, task.ErrTaskNotFound
    }

    updated, err := s.repo.Update(ctx, t)
    if err != nil {
        return nil, err
    }

    if updated.IsRecurring() {
        go s.scheduler.ScheduleRecurringTasks(context.Background())
    }

    return updated, nil
}

func (s *Service) DeleteTask(ctx context.Context, id int) error {
    existing, err := s.repo.GetByID(ctx, id)
    if err != nil {
        return err
    }
    if existing == nil {
        return task.ErrTaskNotFound
    }

    return s.repo.Delete(ctx, id)
}