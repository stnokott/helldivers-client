package db

import (
	"context"
	"fmt"

	"github.com/stnokott/helldivers-client/internal/db/gen"
)

// compile-time implementation check
var _ EntityMerger = (*Assignment)(nil)

// Assignment implements EntityMerger
type Assignment struct {
	gen.Assignment
	Tasks []gen.AssignmentTask
}

func (a *Assignment) Merge(ctx context.Context, tx *gen.Queries, stats tableMergeStats) error {
	// Since we have static assignment IDs, but Identity task IDs, we cannot easily merge both together.
	// (Composite types also don't work properly yet, see https://github.com/sqlc-dev/sqlc/issues/2760)
	// This is why we apply the following procedure:
	//   1. If assignment exists, delete all connected tasks first
	//   2. Then, merge the assignment as usual, re-inserting the tasks along the way
	exists, err := tx.AssignmentExists(ctx, a.ID)
	if err != nil {
		return fmt.Errorf("failed to check if assignment (ID=%d) exists: %w", a.ID, err)
	}
	if exists {
		if err = tx.DeleteAssignmentTasks(ctx, a.TaskIds); err != nil {
			return fmt.Errorf("failed to delete assignment tasks: %w", err)
		}
	}

	taskIDs, err := insertAssignmentTasks(ctx, tx, a.Tasks, stats)
	if err != nil {
		return err
	}
	a.TaskIds = taskIDs

	if _, err = tx.MergeAssignment(ctx, gen.MergeAssignmentParams(a.Assignment)); err != nil {
		return fmt.Errorf("failed to insert assignment '%s': %v", a.Title, err)
	}
	stats.Incr("Assignments", false, 1)
	return nil
}

func insertAssignmentTasks(ctx context.Context, tx *gen.Queries, tasks []gen.AssignmentTask, stats tableMergeStats) ([]int64, error) {
	taskIDs := make([]int64, len(tasks))
	for i, task := range tasks {
		taskID, err := tx.InsertAssignmentTask(ctx, gen.InsertAssignmentTaskParams{
			TaskType:   task.TaskType,
			Values:     task.Values,
			ValueTypes: task.ValueTypes,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to insert assignment task (ID=%d): %v", task.ID, err)
		}
		stats.Incr("Assignment Tasks", false, 1)
		taskIDs[i] = taskID
	}
	return taskIDs, nil
}
