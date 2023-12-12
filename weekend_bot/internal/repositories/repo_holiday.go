package repositories

import (
	"context"
	"errors"
	"fmt"
	"time"
	"weekend_bot/internal/apperrors"
	"weekend_bot/internal/config"
	"weekend_bot/internal/models"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2/bson"
)

const dateCollectionHoliday = "fast_days"

type HolidayRepo interface {
	SaveDateIfNotExist(ctx context.Context, user *models.Holiday) error
	GetAllDates(ctx context.Context) ([]*models.Holiday, error)
	UpdateModification(ctx context.Context, user *models.Holiday) error
	GetTodaysDate(ctx context.Context) (*models.Holiday, error)
	SaveAllDates(ctx context.Context, weekends []*models.Holiday) error
}

type holidayRepo struct {
	log        *logrus.Logger
	collection *mongo.Collection
}

func NewHolidayRepo(config *config.Config, log *logrus.Logger, mongoDB *mongo.Database) HolidayRepo {
	return &holidayRepo{log: log, collection: mongoDB.Collection(dateCollectionHoliday)}
}

func (hr *holidayRepo) SaveAllDates(ctx context.Context, holiday []*models.Holiday) error {
	if holiday == nil {
		appErr := apperrors.HolidaySaveAllDatesErr.AppendMessage("holiday == nil")
		return appErr
	}

	for _, date := range holiday {
		err := hr.SaveDateIfNotExist(ctx, date)
		if err != nil {
			appErr := apperrors.HolidaySaveAllDatesErr.AppendMessage(err)
			hr.log.Error(appErr)
			return appErr
		}
	}

	hr.log.Info("All weekends saved in mongoDB")
	return nil
}

func (hr *holidayRepo) GetTodaysDate(ctx context.Context) (*models.Holiday, error) {
	today := time.Now()
	dateTime, err := convertDateTime(today.Format("02.01.2006"))
	if err != nil {
		appErr := apperrors.HolidayGetTodaysDateErr.AppendMessage(err)
		hr.log.Error(appErr)
		return nil, appErr
	}

	filter := bson.M{"date": dateTime}
	res := hr.collection.FindOne(ctx, filter)
	if res.Err() != nil {
		appErr := apperrors.HolidayGetTodaysDateErr.AppendMessage(fmt.Sprintf("Cannot find today's date. Err: %+v", res.Err()))
		hr.log.Error(appErr)
		return nil, appErr
	}

	dateToday := models.Holiday{}
	err = res.Decode(&dateToday)
	if err != nil {
		appErr := apperrors.HolidayGetTodaysDateErr.AppendMessage(err)
		hr.log.Error(appErr)
		return nil, appErr
	}

	return &dateToday, nil
}

func (fr *holidayRepo) SaveDateIfNotExist(ctx context.Context, date *models.Holiday) error {
	if date == nil {
		appErr := apperrors.HolidaySaveDateIfNotExistErr.AppendMessage("date == nil")
		fr.log.Error(appErr)
		return appErr
	}

	filter := bson.M{"chat_id": date.ID}
	res := fr.collection.FindOne(ctx, filter)

	if res.Err() != nil && !errors.Is(res.Err(), mongo.ErrNoDocuments) {
		appErr := apperrors.HolidaySaveDateIfNotExistErr.AppendMessage(fmt.Sprintf("res.Err: %+v", res.Err()))
		fr.log.Error(appErr)
		return appErr
	}

	if res.Err() == nil {
		appErr := apperrors.HolidaySaveDateIfNotExistErr.AppendMessage(fmt.Sprintf("The holiday already exists in the database with ChatID: %v", date.ID))
		fr.log.Error(appErr)
		return appErr
	}

	_, err := fr.collection.InsertOne(ctx, date)
	if err != nil {
		appErr := apperrors.HolidaySaveDateIfNotExistErr.AppendMessage(err)
		fr.log.Error(appErr)
		return appErr
	}

	fr.log.Info("The holiday has been added, successfully.")
	return nil
}

func (fr *holidayRepo) GetAllDates(ctx context.Context) ([]*models.Holiday, error) {
	cursor, err := fr.collection.Find(ctx, bson.M{})
	if err != nil {
		appErr := apperrors.HolidayGetAllDatesErr.AppendMessage(err)
		fr.log.Error(appErr)
		return nil, appErr
	}

	defer cursor.Close(ctx)

	return decodeHolidays(ctx, cursor)
}

func (fr *holidayRepo) UpdateModification(ctx context.Context, date *models.Holiday) error {
	filter := bson.M{"chat_id": date.ID}
	update := bson.M{
		"$set": bson.M{
			"updated": date.Updated,
		},
	}

	res, err := fr.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		appErr := apperrors.UpdateModFailedError.AppendMessage(err)
		fr.log.Error(appErr)
		return appErr
	}

	if res.ModifiedCount == 0 {
		appErr := apperrors.UpdateModFailedError.AppendMessage("No documents were updated")
		fr.log.Error(appErr)
		return appErr
	}

	return nil
}

func decodeHolidays(ctx context.Context, cursor *mongo.Cursor) ([]*models.Holiday, error) {
	defer cursor.Close(ctx)

	var dates []*models.Holiday
	for cursor.Next(ctx) {
		var date models.Holiday
		err := cursor.Decode(&date)
		if err != nil {
			appErr := apperrors.DecodeHolidaysErr.AppendMessage(err)
			return nil, appErr
		}

		dates = append(dates, &date)
	}

	if err := cursor.Err(); err != nil {
		appErr := apperrors.DecodeHolidaysErr.AppendMessage(err)
		return nil, appErr
	}

	return dates, nil
}
