package task

import (
    "context"
    "encoding/json"
    "strings"
    "test-task/internal/domain/task"
    "time"
)

type TaskScheduler struct {
    repo Repository
}

func NewScheduler(repo Repository) *TaskScheduler {
    return &TaskScheduler{repo: repo}
}

func (s *TaskScheduler) ScheduleRecurringTasks(ctx context.Context) error {
    tasks, err := s.repo.GetRecurringTasks(ctx)
    if err != nil {
        return err
    }

    now := time.Now()

    for _, t := range tasks {
        if !t.IsRecurring() {
            continue
        }

        nextDates := s.calculateNextOccurrences(t, now, 90)

        for _, nextDate := range nextDates {
            if t.RecurrenceEndDate != nil && nextDate.After(*t.RecurrenceEndDate) {
                continue
            }

            newTask := &task.Task{
                Title:        t.Title,
                Description:  t.Description,
                Status:       "new",
                ScheduledAt:  &nextDate,
                ParentTaskID: &t.ID,
            }

            s.repo.Create(ctx, newTask)
        }
    }

    return nil
}

func (s *TaskScheduler) calculateNextOccurrences(t *task.Task, from time.Time, daysAhead int) []time.Time {
    var dates []time.Time

    switch *t.RecurrenceType {
    case task.RecurrenceDaily:
        interval := *t.RecurrenceInterval
        current := from
        for i := 0; i < daysAhead/interval; i++ {
            current = current.AddDate(0, 0, interval)
            dates = append(dates, current)
        }

    case task.RecurrenceMonthly:
        dayOfMonth := *t.RecurrenceDayOfMonth
        current := from
        for i := 0; i < 12; i++ {
            current = current.AddDate(0, 1, 0)
            nextDate := time.Date(current.Year(), current.Month(), dayOfMonth, 0, 0, 0, 0, current.Location())
            if nextDate.Day() == dayOfMonth {
                dates = append(dates, nextDate)
            }
        }

    case task.RecurrenceEvenDays:
        current := from
        for i := 0; i < daysAhead; i++ {
            current = current.AddDate(0, 0, 1)
            if current.Day()%2 == 0 {
                dates = append(dates, current)
            }
        }

    case task.RecurrenceOddDays:
        current := from
        for i := 0; i < daysAhead; i++ {
            current = current.AddDate(0, 0, 1)
            if current.Day()%2 == 1 {
                dates = append(dates, current)
            }
        }

    case task.RecurrenceSpecificDates:
        if t.RecurrenceSpecificDates != nil {
            var specificDates []string
            json.Unmarshal([]byte(*t.RecurrenceSpecificDates), &specificDates)

            for _, dateStr := range specificDates {
                dateStr = strings.TrimSpace(dateStr)
                if parsed, err := time.Parse("2006-01-02", dateStr); err == nil {
                    if parsed.After(from) {
                        dates = append(dates, parsed)
                    }
                }
            }
        }
    }

    return dates
}