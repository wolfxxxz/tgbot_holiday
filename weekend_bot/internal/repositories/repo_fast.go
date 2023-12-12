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
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2/bson"
)

const dateCollection = "post_days"

type FastRepo interface {
	SaveDateIfNotExist(ctx context.Context, user *models.Fast) error
	GetAllDates(ctx context.Context) ([]*models.Fast, error)
	UpdateModification(ctx context.Context, user *models.Fast) error
	GetTodaysDate(ctx context.Context) (*models.Fast, error)
	SaveAllDates(ctx context.Context, weekends []*models.Fast) error
}

type fastRepo struct {
	log        *logrus.Logger
	collection *mongo.Collection
}

func NewFastRepo(config *config.Config, log *logrus.Logger, mongoDB *mongo.Database) FastRepo {
	return &fastRepo{log: log, collection: mongoDB.Collection(dateCollection)}
}

func (fr *fastRepo) SaveAllDates(ctx context.Context, fasts []*models.Fast) error {
	if fasts == nil {
		appErr := apperrors.MongoSaveAllDatesError.AppendMessage("insert Data If Not Exist user == nil")
		fr.log.Error(appErr)
		return appErr
	}

	for _, weekend := range fasts {
		err := fr.SaveDateIfNotExist(ctx, weekend)
		if err != nil {
			appErr := apperrors.MongoSaveAllDatesError.AppendMessage(err)
			fr.log.Error(appErr)
			return appErr
		}
	}

	fr.log.Info("All weekends saved in mongoDB")
	return nil
}

func (fr *fastRepo) GetTodaysDate(ctx context.Context) (*models.Fast, error) {
	today := time.Now()
	dateTime, err := convertDateTime(today.Format("02.01.2006"))
	if err != nil {
		appErr := err.(*apperrors.AppError)
		fr.log.Error(appErr)
		return nil, appErr
	}

	filter := bson.M{"date": dateTime}
	res := fr.collection.FindOne(ctx, filter)
	if res.Err() != nil {
		appErr := apperrors.GetTodaysDate.AppendMessage(fmt.Sprintf("Cannot find today's date. Err: %+v", res.Err()))
		fr.log.Error(appErr)
		return nil, appErr
	}

	dateToday := models.Fast{}
	err = res.Decode(&dateToday)
	if err != nil {
		appErr := apperrors.GetTodaysDate.AppendMessage(err)
		fr.log.Error(appErr)
		return nil, appErr
	}

	return &dateToday, nil
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

func (fr *fastRepo) SaveDateIfNotExist(ctx context.Context, date *models.Fast) error {
	if date == nil {
		appErr := apperrors.SaveDateIfNotExist.AppendMessage("date == nil")
		fr.log.Error(appErr)
		return appErr
	}

	filter := bson.M{"chat_id": date.ID}
	res := fr.collection.FindOne(ctx, filter)

	if res.Err() != nil && !errors.Is(res.Err(), mongo.ErrNoDocuments) {
		appErr := apperrors.SaveDateIfNotExist.AppendMessage(fmt.Sprintf("Cannot find user by chat_id. Err: %+v", res.Err()))
		fr.log.Error(appErr)
		return appErr
	}

	if res.Err() == nil {
		appErr := apperrors.SaveDateIfNotExist.AppendMessage(fmt.Sprintf("The user already exists in the database with ChatID: %v", date.ID))
		fr.log.Info(appErr)
		return nil
	}

	_, err := fr.collection.InsertOne(ctx, date)
	if err != nil {
		appErr := apperrors.SaveDateIfNotExist.AppendMessage(err)
		fr.log.Info(appErr)
		return appErr
	}

	fr.log.Info("The user has been added, successfully.")
	return nil
}

func (fr *fastRepo) GetAllDates(ctx context.Context) ([]*models.Fast, error) {
	cursor, err := fr.collection.Find(ctx, bson.M{})
	if err != nil {
		appErr := apperrors.GetAllDatesErr.AppendMessage(err)
		fr.log.Error(appErr)
		return nil, appErr
	}

	defer cursor.Close(ctx)

	return decodeDates(ctx, cursor)
}

func (fr *fastRepo) UpdateModification(ctx context.Context, date *models.Fast) error {
	filter := bson.M{"chat_id": date.ID}
	update := bson.M{
		"$set": bson.M{
			"updated": date.Updated,
		},
	}

	res, err := fr.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return apperrors.MongoUpdateModFailedError.AppendMessage(err)
	}

	if res.ModifiedCount == 0 {
		return apperrors.MongoUpdateModFailedError.AppendMessage("No documents were updated")
	}

	return nil
}

func decodeDates(ctx context.Context, cursor *mongo.Cursor) ([]*models.Fast, error) {
	defer cursor.Close(ctx)

	var dates []*models.Fast
	for cursor.Next(ctx) {
		var date models.Fast
		err := cursor.Decode(&date)
		if err != nil {
			appErr := apperrors.DecodeDatesErr.AppendMessage(err)
			return nil, appErr
		}

		dates = append(dates, &date)
	}

	if err := cursor.Err(); err != nil {
		appErr := apperrors.DecodeDatesErr.AppendMessage(err)
		return nil, appErr
	}
	return dates, nil
}
