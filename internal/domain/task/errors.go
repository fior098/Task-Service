package task

import "errors"

var (
    ErrTaskNotFound                   = errors.New("task not found")
    ErrInvalidTitle                   = errors.New("invalid title")
    ErrInvalidStatus                  = errors.New("invalid status")
    ErrInvalidRecurrenceType          = errors.New("invalid recurrence type")
    ErrInvalidRecurrenceInterval      = errors.New("recurrence interval must be >= 1")
    ErrInvalidRecurrenceDayOfMonth    = errors.New("recurrence day of month must be between 1 and 30")
    ErrInvalidRecurrenceSpecificDates = errors.New("recurrence specific dates cannot be empty")
)