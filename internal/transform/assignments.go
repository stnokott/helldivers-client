package transform

import (
	"errors"
	"fmt"

	"github.com/stnokott/helldivers-client/internal/api"
	"github.com/stnokott/helldivers-client/internal/db"
	"github.com/stnokott/helldivers-client/internal/db/gen"
)

func Assignments(data APIData) ([]db.EntityMerger, error) {
	if data.Assignments == nil {
		return nil, errors.New("got nil assignments slice")
	}

	src := *data.Assignments
	assignments := make([]db.EntityMerger, len(src))
	for i, assignment := range src {
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

		assignments[i] = &db.Assignment{
			Assignment: gen.Assignment{
				ID:           *assignment.Id,
				Title:        *assignment.Title,
				Briefing:     *assignment.Briefing,
				Description:  *assignment.Description,
				Expiration:   db.PGTimestamp(*assignment.Expiration),
				RewardType:   *reward.Type,
				RewardAmount: *reward.Amount,
			},
			Tasks: tasks,
		}
	}
	return assignments, nil
}

func parseAssignmentReward(in *api.Assignment2_Reward) (api.Reward2, error) {
	reward, err := in.AsReward2()
	if err != nil {
		return api.Reward2{}, fmt.Errorf("parse assignment reward: %w", err)
	}
	if reward.Amount == nil || reward.Type == nil {
		return api.Reward2{}, errFromNils(&reward)
	}
	return reward, nil
}

func convertAssignmentTasks(in *[]api.Task2) ([]gen.AssignmentTask, error) {
	tasks := make([]gen.AssignmentTask, len(*in))
	for i, task := range *in {
		if task.Type == nil || task.ValueTypes == nil || task.Values == nil {
			return nil, errFromNils(&task)
		}
		tasks[i] = gen.AssignmentTask{
			TaskType:   *task.Type,
			Values:     *task.Values,
			ValueTypes: *task.ValueTypes,
		}
	}
	return tasks, nil
}
