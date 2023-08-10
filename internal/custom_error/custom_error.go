package custom_error

import "errors"

var (
	ErrEmptyID               = errors.New("empty id param")
	ErrInvalidIDParameter    = errors.New("invalid id param")
	ErrTaskNotFound          = errors.New("task not found")
	ErrMessageTooLong        = errors.New("more than 200 char")
	ErrInvalidActiveAtFormat = errors.New("activeAt invalid format")
	ErrDuplicateTask         = errors.New("a task with the same title already exists")
)
