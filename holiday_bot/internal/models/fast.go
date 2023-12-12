package models

import (
	"holiday_bot/internal/apperrors"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Fast struct {
	ID      *uuid.UUID          `bson:"_id"`
	Title   string              `bson:"title"`
	Date    *primitive.DateTime `bson:"date"`
	Created *primitive.DateTime `bson:"_created"`
	Updated *primitive.DateTime `bson:"_updated"`
}

func CreateWeekend(title string, date string) (*Fast, error) {
	id := uuid.New()
	now := primitive.NewDateTimeFromTime(time.Now())
	dateTime, err := convertDateTime(date)
	if err != nil {
		appErr := apperrors.ModelsCreateFastErr.AppendMessage(err)
		return nil, appErr
	}

	return &Fast{ID: &id, Title: title, Date: dateTime, Created: &now}, nil

}

func convertDateTime(dateTime string) (*primitive.DateTime, error) {
	parsedTime, err := time.Parse("02.01.2006", dateTime)
	if err != nil {
		appErr := apperrors.ConvertDateTimeErr.AppendMessage(err)
		return nil, appErr
	}

	dateTimePrimitive := primitive.NewDateTimeFromTime(parsedTime)
	return &dateTimePrimitive, nil
}
