package workers

import "fmt"

type JobPanic struct {
	Value any
}

func (j *JobPanic) Error() string {
	return fmt.Sprintf("worker recovered: %v", j.Value)
}

type JobError struct {
	Value error
}

func (j *JobError) Error() string {
	return fmt.Sprintf("job error: %v", j.Value)
}
