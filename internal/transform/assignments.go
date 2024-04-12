package transform

import (
	"errors"

	"github.com/stnokott/helldivers-client/internal/db"
	"github.com/stnokott/helldivers-client/internal/db/structs"
)

// Assignments implements worker.docTransformer
type Assignments struct{}

func (_ Assignments) Transform(data APIData) (*db.DocsProvider[structs.Assignment], error) {
	if data.Assignments == nil {
		return nil, errors.New("got nil assignments slice")
	}

	assignments := *data.Assignments
	assignmentDocs := make([]db.DocWrapper[structs.Assignment], len(assignments))

	for i, assignment := range assignments {
		if assignment.Id == nil ||
			assignment.Title == nil ||
			assignment.Briefing == nil ||
			assignment.Description == nil ||
			assignment.Expiration == nil ||
			assignment.Tasks == nil ||
			assignment.Reward == nil {
			return nil, errFromNils(&assignment)
		}

		reward, err := assignment.Reward.AsReward2()
		if err != nil {
			return nil, err
		}
		tasksRaw := *assignment.Tasks
		tasks := make([]structs.AssignmentTask, len(tasksRaw))
		for i, task := range tasksRaw {
			tasks[i] = structs.AssignmentTask{
				Type:       *task.Type,
				Values:     *task.Values,
				ValueTypes: *task.ValueTypes,
			}
		}

		assignmentDocs[i] = db.DocWrapper[structs.Assignment]{
			DocID: *assignment.Id,
			Document: structs.Assignment{
				ID:       *assignment.Id,
				Title:    *assignment.Title,
				Briefing: *assignment.Briefing,

				Description: *assignment.Description,
				Expiration:  db.PrimitiveTime(*assignment.Expiration),
				Progress:    *assignment.Progress,
				Reward: structs.AssignmentReward{
					Type:   *reward.Type,
					Amount: *reward.Amount,
				},
				Tasks: tasks,
			},
		}
	}
	return &db.DocsProvider[structs.Assignment]{
		CollectionName: db.CollAssignments,
		Docs:           assignmentDocs,
	}, nil
}
