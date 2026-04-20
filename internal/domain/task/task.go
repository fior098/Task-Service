package task

import "time"

type RecurrenceType string

const (
    RecurrenceNone          RecurrenceType = ""
    RecurrenceDaily         RecurrenceType = "daily"
    RecurrenceMonthly       RecurrenceType = "monthly"
    RecurrenceSpecificDates RecurrenceType = "specific_dates"
    RecurrenceEvenDays      RecurrenceType = "even_days"
    RecurrenceOddDays       RecurrenceType = "odd_days"
)

type Task struct {
    ID                     int             `json:"id"`
    Title                  string          `json:"title"`
    Description            string          `json:"description"`
    Status                 string          `json:"status"`
    ScheduledAt            *time.Time      `json:"scheduled_at,omitempty"`
    CreatedAt              time.Time       `json:"created_at"`
    UpdatedAt              time.Time       `json:"updated_at"`
    RecurrenceType         *RecurrenceType `json:"recurrence_type,omitempty"`
    RecurrenceInterval     *int            `json:"recurrence_interval,omitempty"`
    RecurrenceDayOfMonth   *int            `json:"recurrence_day_of_month,omitempty"`
    RecurrenceParity       *string         `json:"recurrence_parity,omitempty"`
    RecurrenceSpecificDates *string        `json:"recurrence_specific_dates,omitempty"`
    RecurrenceEndDate      *time.Time      `json:"recurrence_end_date,omitempty"`
    ParentTaskID           *int            `json:"parent_task_id,omitempty"`
}

func (t *Task) IsRecurring() bool {
    return t.RecurrenceType != nil && *t.RecurrenceType != RecurrenceNone
}

func (t *Task) Validate() error {
    if t.Title == "" {
        return ErrInvalidTitle
    }

    if t.Status == "" {
        return ErrInvalidStatus
    }

    if t.RecurrenceType != nil {
        switch *t.RecurrenceType {
        case RecurrenceDaily:
            if t.RecurrenceInterval == nil || *t.RecurrenceInterval < 1 {
                return ErrInvalidRecurrenceInterval
            }
        case RecurrenceMonthly:
            if t.RecurrenceDayOfMonth == nil || *t.RecurrenceDayOfMonth < 1 || *t.RecurrenceDayOfMonth > 30 {
                return ErrInvalidRecurrenceDayOfMonth
            }
        case RecurrenceSpecificDates:
            if t.RecurrenceSpecificDates == nil || *t.RecurrenceSpecificDates == "" {
                return ErrInvalidRecurrenceSpecificDates
            }
        case RecurrenceEvenDays, RecurrenceOddDays:
        case RecurrenceNone:
        default:
            return ErrInvalidRecurrenceType
        }
    }

    return nil
}