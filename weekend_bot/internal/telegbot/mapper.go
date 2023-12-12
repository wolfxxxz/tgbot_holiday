package telegbot

import (
	"fmt"
	"weekend_bot/internal/apperrors"
	"weekend_bot/internal/models"
)

func MapModelsPostToResponse(weekend *models.Fast) (string, error) {
	if weekend == nil {
		return "", apperrors.BotMapperEncodingErr.AppendMessage("data.Main == nil")
	}

	response := fmt.Sprintf("Дата сьогодні %v, дозволено їсти %v", weekend.Date.Time().Format("02.01.2006"), weekend.Title)
	return response, nil
}

func MapModelsHolidayToResponse(weekend *models.Holiday) (string, error) {
	if weekend == nil {
		return "", apperrors.BotMapperEncodingErr.AppendMessage("data.Main == nil")
	}

	response := fmt.Sprintf("Дата сьогодні %v, свято %v", weekend.Date.Time().Format("02.01.2006"), weekend.Title)
	return response, nil
}
