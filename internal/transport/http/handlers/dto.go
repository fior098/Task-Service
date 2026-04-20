package handlers

import "time"

type CreateTaskRequest struct {
    Title                  string     `json:"title"`
    Description            string     `json:"description"`
    Status                 string     `json:"status"`
    ScheduledAt            *time.Time `json:"scheduled_at,omitempty"`
    RecurrenceType         *string    `json:"recurrence_type,omitempty"`
    RecurrenceInterval     *int       `json:"recurrence_interval,omitempty"`
    RecurrenceDayOfMonth   *int       `json:"recurrence_day_of_month,omitempty"`
    RecurrenceParity       *string    `json:"recurrence_parity,omitempty"`
    RecurrenceSpecificDates *string   `json:"recurrence_specific_dates,omitempty"`
    RecurrenceEndDate      *time.Time `json:"recurrence_end_date,omitempty"`
}

type UpdateTaskRequest struct {
    Title                  string     `json:"title"`
    Description            string     `json:"description"`
    Status                 string     `json:"status"`
    ScheduledAt            *time.Time `json:"scheduled_at,omitempty"`
    RecurrenceType         *string    `json:"recurrence_type,omitempty"`
    RecurrenceInterval     *int       `json:"recurrence_interval,omitempty"`
    RecurrenceDayOfMonth   *int       `json:"recurrence_day_of_month,omitempty"`
    RecurrenceParity       *string    `json:"recurrence_parity,omitempty"`
    RecurrenceSpecificDates *string   `json:"recurrence_specific_dates,omitempty"`
    RecurrenceEndDate      *time.Time `json:"recurrence_end_date,omitempty"`
}

type ErrorResponse struct {
    Error string `json:"error"`
}