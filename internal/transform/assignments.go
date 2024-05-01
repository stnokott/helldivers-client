package transform

import (
	"errors"
	"fmt"

	"github.com/stnokott/helldivers-client/internal/api"
	"github.com/stnokott/helldivers-client/internal/db"
)

func Assignments(c Converter, data APIData) ([]db.EntityMerger, error) {
	if data.Assignments == nil {
		return nil, errors.New("got nil assignments slice")
	}

	src := *data.Assignments
	mergers := make([]db.EntityMerger, len(src))
	for i, assignment := range src {
		a, err := c.ConvertAssignment(assignment)
		if err != nil {
			return nil, err
		}
		mergers[i] = db.EntityMerger(a)
	}
	return mergers, nil
}

func parseAssignmentRewardType(source *api.Assignment2_Reward) (int32, error) {
	parsed, err := parseAssignmentReward(source)
	if err != nil {
		return -1, err
	}
	if parsed.Type == nil {
		return -1, errors.New("reward type is nil")
	}
	return *parsed.Type, nil
}

func parseAssignmentRewardAmount(source *api.Assignment2_Reward) (int32, error) {
	parsed, err := parseAssignmentReward(source)
	if err != nil {
		return -1, err
	}
	if parsed.Amount == nil {
		return -1, errors.New("reward amount is nil")
	}
	return *parsed.Amount, nil
}

func parseAssignmentReward(in *api.Assignment2_Reward) (api.Reward2, error) {
	reward, err := in.AsReward2()
	if err != nil {
		return api.Reward2{}, fmt.Errorf("parse assignment reward: %w", err)
	}
	return reward, nil
}
