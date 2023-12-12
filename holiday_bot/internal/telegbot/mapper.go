package telegbot

import (
	"fmt"
	"holiday_bot/internal/apperrors"
	"holiday_bot/internal/models"
)

func MapModelsPostToResponse(holiday *models.Fast) (string, error) {
	if holiday == nil {
		return "", apperrors.BotMapperEncodingErr.AppendMessage("data.Main == nil")
	}

	response := fmt.Sprintf("Дата сьогодні %v, дозволено їсти %v", holiday.Date.Time().Format("02.01.2006"), holiday.Title)
	return response, nil
}

func MapModelsHolidayToResponse(holiday *models.Holiday) (string, error) {
	if holiday == nil {
		return "", apperrors.BotMapperEncodingErr.AppendMessage("data.Main == nil")
	}

	response := fmt.Sprintf("Дата сьогодні %v, свято %v", holiday.Date.Time().Format("02.01.2006"), holiday.Title)
	return response, nil
}
