package db

import (
	"context"
	"fmt"
	"log"

	"github.com/stnokott/helldivers-client/internal/db/gen"
)

// compile-time implementation check
var _ EntityMerger = (*Assignment)(nil)

// Assignment implements EntityMerger
type Assignment struct {
	gen.Assignment
	Tasks []AssignmentTask
}

type AssignmentTask gen.AssignmentTask

func (a *Assignment) Merge(ctx context.Context, tx *gen.Queries, stats *MergeStats, logger *log.Logger) error {
	logger.Printf("** merging assignment '%s'", a.Title)
	taskIDs, err := mergeAssignmentTasks(ctx, tx, a.Tasks, logger)
	if err != nil {
		return err
	}
	a.TaskIds = taskIDs

	id, err := tx.GetAssignment(ctx, a.ID)
	exists, err := entityExistsByPK(id, err, a.ID)
	if err != nil {
		return fmt.Errorf("failed to check for existing assignment: %v", err)
	}
	if exists {
		// perform UPDATE
		if _, err = tx.UpdateAssignment(ctx, gen.UpdateAssignmentParams(a.Assignment)); err != nil {
			return fmt.Errorf("failed to update assignment ('%s'): %v", a.Title, err)
		}
		stats.Updates++
	} else {
		// perform INSERT
		if _, err = tx.InsertAssignment(ctx, gen.InsertAssignmentParams(a.Assignment)); err != nil {
			return fmt.Errorf("failed to insert assignment ('%s'): %v", a.Title, err)
		}
		stats.Inserts++
	}
	return nil
}

func mergeAssignmentTasks(ctx context.Context, tx *gen.Queries, tasks []AssignmentTask, logger *log.Logger) ([]int64, error) {
	taskIDs := make([]int64, len(tasks))
	for i, task := range tasks {
		id, err := tx.GetAssignmentTask(ctx, task.ID)
		exists, err := entityExistsByPK(id, err, task.ID)
		var taskID int64
		if exists {
			// perform UPDATE
			taskID, err = tx.UpdateAssignmentTask(ctx, gen.UpdateAssignmentTaskParams(task))
			if err != nil {
				return nil, fmt.Errorf("failed to update assignment task (ID=%d): %v", task.ID, err)
			}
		} else {
			// perform INSERT
			taskID, err = tx.InsertAssignmentTask(ctx, gen.InsertAssignmentTaskParams{
				Type:       task.Type,
				Values:     task.Values,
				ValueTypes: task.ValueTypes,
			})
			if err != nil {
				return nil, fmt.Errorf("failed to insert assignment task (ID=%d): %v", task.ID, err)
			}
		}
		taskIDs[i] = taskID
	}
	return taskIDs, nil
}
