package models

import (
	"holiday_bot/internal/apperrors"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Holiday struct {
	ID      *uuid.UUID          `bson:"_id"`
	Title   string              `bson:"title"`
	Date    *primitive.DateTime `bson:"date"`
	Created *primitive.DateTime `bson:"_created"`
	Updated *primitive.DateTime `bson:"_updated"`
}

func CreateHoliday(title string, date string) (*Holiday, error) {
	id := uuid.New()
	now := primitive.NewDateTimeFromTime(time.Now())
	dateTime, err := convertDateTime(date)
	if err != nil {
		appErr := apperrors.ModelsCreateHolidayErr.AppendMessage(err)
		return nil, appErr
	}

	return &Holiday{ID: &id, Title: title, Date: dateTime, Created: &now}, nil

}
