package transform

import (
	"errors"
	"fmt"

	"github.com/stnokott/helldivers-client/internal/api"
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

		reward, err := parseAssignmentReward(assignment.Reward)
		if err != nil {
			return nil, err
		}
		tasks, err := convertAssignmentTasks(assignment.Tasks)
		if err != nil {
			return nil, err
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

func parseAssignmentReward(in *api.Assignment2_Reward) (api.Reward2, error) {
	reward, err := in.AsReward2()
	if err != nil {
		return api.Reward2{}, fmt.Errorf("cannot parse assignment reward: %w", err)
	}
	if reward.Amount == nil || reward.Type == nil {
		return api.Reward2{}, errFromNils(&reward)
	}
	return reward, nil
}

func convertAssignmentTasks(in *[]api.Task2) ([]structs.AssignmentTask, error) {
	tasks := make([]structs.AssignmentTask, len(*in))
	for i, task := range *in {
		if task.Type == nil || task.ValueTypes == nil || task.Values == nil {
			return nil, errFromNils(&task)
		}
		tasks[i] = structs.AssignmentTask{
			Type:       *task.Type,
			Values:     *task.Values,
			ValueTypes: *task.ValueTypes,
		}
	}
	return tasks, nil
}
