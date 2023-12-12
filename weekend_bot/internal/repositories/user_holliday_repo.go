package repositories

import (
	"context"
	"errors"
	"fmt"
	"weekend_bot/internal/apperrors"
	"weekend_bot/internal/config"
	"weekend_bot/internal/models"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2/bson"
)

const userHolidayCollection = "usersTelegramHoliday"

type UserHolidayRepo interface {
	SaveUserIfNotExist(ctx context.Context, user *models.User) error
	GetAllUsers(ctx context.Context) ([]*models.User, error)
	UpdateModification(ctx context.Context, user *models.User) error
	DropUser(ctx context.Context, chatID int64) error
}

type userHolydayRepo struct {
	log        *logrus.Logger
	collection *mongo.Collection
}

func NewUserHolidayRepo(config *config.Config, log *logrus.Logger, mongoDB *mongo.Database) UserHolidayRepo {
	return &userHolydayRepo{log: log, collection: mongoDB.Collection(userHolidayCollection)}
}

func (uhr *userHolydayRepo) DropUser(ctx context.Context, chatID int64) error {
	if chatID == 0 {
		appErr := apperrors.UHRDropUserErr.AppendMessage("chatID = null")
		uhr.log.Error(appErr)
		return appErr
	}

	filter := bson.M{"chat_id": chatID}
	res, err := uhr.collection.DeleteOne(ctx, filter)

	if err != nil {
		appErr := apperrors.UHRDropUserErr.AppendMessage(fmt.Sprintf("Cannot find user by chat_id. Err: %+v", err))
		uhr.log.Error(appErr)
		return appErr
	}

	if res.DeletedCount != 1 {
		appErr := apperrors.UHRDropUserErr.AppendMessage(res.DeletedCount)
		uhr.log.Error(appErr)
		return nil
	}

	uhr.log.Info("The user has been deleted, successfully.")
	return nil
}

func (uhr *userHolydayRepo) SaveUserIfNotExist(ctx context.Context, user *models.User) error {
	if user == nil {
		appErr := apperrors.UHRSaveUserIfNotExistErr.AppendMessage("insert Data If Not Exist user == nil")
		uhr.log.Error(appErr)
		return appErr
	}

	filter := bson.M{"chat_id": user.ChatID}
	res := uhr.collection.FindOne(ctx, filter)

	if res.Err() != nil && !errors.Is(res.Err(), mongo.ErrNoDocuments) {
		appErr := apperrors.UHRSaveUserIfNotExistErr.AppendMessage(fmt.Sprintf("Cannot find user by chat_id. Err: %+v", res.Err()))
		uhr.log.Error(appErr)
		return appErr
	}

	if res.Err() == nil {
		appErr := apperrors.UHRSaveUserIfNotExistErr.AppendMessage(fmt.Sprintf("The user already exists in the database with ChatID: %v", user.ChatID))
		uhr.log.Error(appErr)
		return nil
	}

	_, err := uhr.collection.InsertOne(ctx, user)
	if err != nil {
		appErr := apperrors.UHRSaveUserIfNotExistErr.AppendMessage(err)
		uhr.log.Error(appErr)
		return appErr
	}

	uhr.log.Info("The user has been added, successfully.")
	return nil
}

func (uhr *userHolydayRepo) GetAllUsers(ctx context.Context) ([]*models.User, error) {
	cursor, err := uhr.collection.Find(ctx, bson.M{})
	if err != nil {
		appErr := apperrors.UHRGetAllUsersErr.AppendMessage(err)
		uhr.log.Error(appErr)
		return nil, appErr
	}

	defer cursor.Close(ctx)

	return decodeUsers(ctx, cursor)
}

func (uhr *userHolydayRepo) UpdateModification(ctx context.Context, user *models.User) error {
	filter := bson.M{"chat_id": user.ChatID}
	update := bson.M{
		"$set": bson.M{
			"updated": user.Updated,
		},
	}

	res, err := uhr.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		appErr := apperrors.UHRUpdateModificationErr.AppendMessage(err)
		uhr.log.Error(appErr)
		return appErr
	}

	if res.ModifiedCount == 0 {
		appErr := apperrors.UHRUpdateModificationErr.AppendMessage("No documents were updated")
		uhr.log.Error(appErr)
		return appErr
	}

	return nil
}
