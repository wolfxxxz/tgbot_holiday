package repositories

import (
	"os"
	"strings"
	"weekend_bot/internal/apperrors"
	"weekend_bot/internal/models"

	"github.com/sirupsen/logrus"
)

type txtRepo struct {
	log *logrus.Logger
}

func NewTxtRepo(log *logrus.Logger) *txtRepo {
	return &txtRepo{log: log}
}

func (t *txtRepo) GetFasts(filetxt string) ([]*models.Fast, error) {
	data, err := os.ReadFile(filetxt)
	if err != nil {
		appErr := apperrors.TxtRepoErr.AppendMessage(err)
		t.log.Error(appErr)
		return nil, appErr
	}

	content := string(data)
	lines := strings.Split(content, "\n")
	var fasts []*models.Fast
	for _, line := range lines {
		if line == "" {
			continue
		}

		dataLines := strings.Split(line, "-")
		if len(dataLines) <= 1 {
			continue
		}

		for i, v := range dataLines {
			dataLines[i] = strings.TrimSpace(v)
		}

		onePist, err := models.CreateWeekend(dataLines[1], dataLines[0])
		if err != nil {
			appErr := apperrors.TxtRepoErr.AppendMessage(err)
			t.log.Error(appErr)
			return nil, appErr
		}

		fasts = append(fasts, onePist)
	}

	return fasts, nil
}

func (tr *txtRepo) GetHolidays(filetxt string) ([]*models.Holiday, error) {
	data, err := os.ReadFile(filetxt)
	if err != nil {
		appErr := apperrors.TxtRepoErr.AppendMessage(err)
		tr.log.Error(appErr)
		return nil, appErr
	}

	content := string(data)
	lines := strings.Split(content, "\n")
	var holidays []*models.Holiday
	for _, line := range lines {
		if line == "" {
			tr.log.Info("lines end")
			continue
		}

		dataLines := strings.Split(line, "-")
		if len(dataLines) <= 1 {
			continue
		}

		for i, v := range dataLines {
			dataLines[i] = strings.TrimSpace(v)
		}

		holiday, err := models.CreateHoliday(dataLines[1], dataLines[0])
		if err != nil {
			appErr := apperrors.TxtRepoErr.AppendMessage(err)
			tr.log.Error(appErr)
			return nil, appErr
		}

		holidays = append(holidays, holiday)
	}

	return holidays, nil
}
