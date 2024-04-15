package transform

import (
	"errors"
	"fmt"

	"github.com/stnokott/helldivers-client/internal/api"
	"github.com/stnokott/helldivers-client/internal/db"
	"github.com/stnokott/helldivers-client/internal/db/structs"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Assignments implements worker.DocTransformer
type Assignments struct{}

// Transform implements the worker.DocTransformer interface
func (Assignments) Transform(data APIData, errFunc func(error)) *db.DocsProvider[structs.Assignment] {
	provider := &db.DocsProvider[structs.Assignment]{
		CollectionName: db.CollAssignments,
		Docs:           []db.DocWrapper[structs.Assignment]{},
	}

	if data.Assignments == nil {
		errFunc(errors.New("got nil assignments slice"))
		return provider
	}

	assignments := *data.Assignments

	for _, assignment := range assignments {
		if assignment.Id == nil ||
			assignment.Title == nil ||
			assignment.Briefing == nil ||
			assignment.Description == nil ||
			assignment.Expiration == nil ||
			assignment.Tasks == nil ||
			assignment.Reward == nil {
			errFunc(errFromNils(&assignment))
			continue
		}

		reward, err := parseAssignmentReward(assignment.Reward)
		if err != nil {
			errFunc(err)
			continue
		}
		tasks, err := convertAssignmentTasks(assignment.Tasks)
		if err != nil {
			errFunc(err)
			continue
		}
		provider.Docs = append(provider.Docs, db.DocWrapper[structs.Assignment]{
			DocID: *assignment.Id,
			Document: structs.Assignment{
				ID:       *assignment.Id,
				Title:    *assignment.Title,
				Briefing: *assignment.Briefing,

				Description: *assignment.Description,
				Expiration:  primitive.NewDateTimeFromTime(*assignment.Expiration),
				Progress:    *assignment.Progress,
				Reward: structs.AssignmentReward{
					Type:   *reward.Type,
					Amount: *reward.Amount,
				},
				Tasks: tasks,
			},
		})
	}
	return provider
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
