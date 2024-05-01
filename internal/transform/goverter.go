package transform

import (
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stnokott/helldivers-client/internal/api"
	"github.com/stnokott/helldivers-client/internal/db"
	"github.com/stnokott/helldivers-client/internal/db/gen"
)

//go:generate goverter gen github.com/stnokott/helldivers-client/internal/transform

// Converter converts API structs into mergable DB structs.
//
// goverter:converter
// goverter:wrapErrors
//
// goverter:extend Must.*
//
// goverter:output:package github.com/stnokott/helldivers-client/internal/transform
// goverter:output:file ./generated.go
type Converter interface {
	ConvertAssignment(source api.Assignment2) (*db.Assignment, error)
	// goverter:map Id ID
	// goverter:ignore TaskIds
	// goverter:map Reward RewardType | parseAssignmentRewardType
	// goverter:map Reward RewardAmount | parseAssignmentRewardAmount
	ConvertSingleAssignment(source api.Assignment2) (gen.Assignment, error)
	// goverter:ignore ID
	// goverter:map Type TaskType
	ConvertAssignmentTask(source api.Task2) (gen.AssignmentTask, error)
	ConvertAssignmentTasks(source []api.Task2) ([]gen.AssignmentTask, error)
	// goverter:map Id ID
	// goverter:map Published CreateTime
	ConvertDispatch(source api.Dispatch) (*db.Dispatch, error)
}

func MustAssignment(c Converter, source api.Assignment2) (*db.Assignment, error) {
	assignment, err := c.ConvertSingleAssignment(source)
	if err != nil {
		return nil, err
	}
	if source.Tasks == nil {
		return nil, errors.New("Tasks is nil")
	}
	tasks, err := c.ConvertAssignmentTasks(*source.Tasks)
	if err != nil {
		return nil, err
	}
	return &db.Assignment{
		Assignment: assignment,
		Tasks:      tasks,
	}, nil
}

func MustInt32Ptr(ptr *int32) (int32, error) {
	return mustPtr(ptr)
}

func MustInt64Ptr(ptr *int64) (int64, error) {
	return mustPtr(ptr)
}

func MustString(ptr *string) (string, error) {
	return mustPtr(ptr)
}

func MustInt32Slice(ptr *[]int32) ([]int32, error) {
	return mustPtr(ptr)
}

func MustTimestamp(ptr *time.Time) (pgtype.Timestamp, error) {
	t, err := mustPtr(ptr)
	if err != nil {
		return pgtype.Timestamp{}, err
	}
	return pgtype.Timestamp{Time: t, Valid: true}, nil
}

// nolint: ireturn
func mustPtr[T any](ptr *T) (T, error) {
	if ptr == nil {
		var zero T
		return zero, fmt.Errorf("%T is nil", ptr)
	}
	return *ptr, nil
}
