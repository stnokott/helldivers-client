package transform

import (
	"errors"
	"fmt"

	"github.com/stnokott/helldivers-client/internal/api"
	"github.com/stnokott/helldivers-client/internal/db"
)

// Assignments converts API data into mergable DB entities.
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
		mergers[i] = a
	}
	return mergers, nil
}

// MustAssignment implements a converter for a single assignment.
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
		Assignment: *assignment,
		Tasks:      tasks,
	}, nil
}

// MustAssignmentTitle returns the default locale representation of a localized assignment title.
func MustAssignmentTitle(source *api.Assignment2_Title) (string, error) {
	if source == nil {
		return "", errors.New("Assignment Title is nil")
	}
	return source.AsAssignment2Title0()
}

// MustAssignmentBriefing returns the default locale representation of a localized assignment briefing.
func MustAssignmentBriefing(source *api.Assignment2_Briefing) (string, error) {
	if source == nil {
		return "", errors.New("Assignment Briefing is nil")
	}
	return source.AsAssignment2Briefing0()
}

// MustAssignmentDescription returns the default locale representation of a localized assignment description.
func MustAssignmentDescription(source *api.Assignment2_Description) (string, error) {
	if source == nil {
		return "", errors.New("Assignment Description is nil")
	}
	return source.AsAssignment2Description0()
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
