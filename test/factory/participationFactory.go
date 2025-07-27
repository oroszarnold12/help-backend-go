package factory

import (
	"help/model"

	"github.com/google/uuid"
)

func NewTestParticipation(overrides ...func(*model.Participation)) model.Participation {
	participation := model.Participation{
		Id:              1,
		Uuid:            uuid.New(),
		ShowOnDashboard: true,
		User:            NewTestUser(),
		Course:          NewTestCourse(),
	}

	for _, override := range overrides {
		override(&participation)
	}

	return participation
}
