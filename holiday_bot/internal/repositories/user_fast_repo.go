package repositories

import (
	"context"
	"errors"
	"fmt"
	"holiday_bot/internal/apperrors"
	"holiday_bot/internal/config"
	"holiday_bot/internal/models"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2/bson"
)

const userPostCollection = "usersTelegramPost"

type UserFastRepo interface {
	SaveUserIfNotExist(ctx context.Context, user *models.User) error
	GetAllUsers(ctx context.Context) ([]*models.User, error)
	UpdateModification(ctx context.Context, user *models.User) error
	DropUser(ctx context.Context, chatID int64) error
}

type userFastRepo struct {
	log        *logrus.Logger
	collection *mongo.Collection
}

func NewUserPostRepo(config *config.Config, log *logrus.Logger, mongoDB *mongo.Database) UserFastRepo {
	return &userFastRepo{log: log, collection: mongoDB.Collection(userPostCollection)}
}

func (ufr *userFastRepo) DropUser(ctx context.Context, chatID int64) error {
	if chatID == 0 {
		appErr := apperrors.UserFastDropUserErr.AppendMessage("chatID = null")
		ufr.log.Error(appErr)
		return appErr
	}

	filter := bson.M{"chat_id": chatID}
	res, err := ufr.collection.DeleteOne(ctx, filter)
	if err != nil {
		appErr := apperrors.UserFastDropUserErr.AppendMessage(fmt.Sprintf("Cannot find user by chat_id. Err: %+v", err))
		ufr.log.Error(appErr)
		return appErr
	}

	if res.DeletedCount != 1 {
		appErr := apperrors.UserFastDropUserErr.AppendMessage(res)
		ufr.log.Error(appErr)
		return nil
	}

	ufr.log.Info("The user has been deleted, successfully.")
	return nil
}

func (ufr *userFastRepo) SaveUserIfNotExist(ctx context.Context, user *models.User) error {
	if user == nil {
		appErr := apperrors.UFRSaveUserIfNotExistErr.AppendMessage("insert Data If Not Exist user == nil")
		ufr.log.Error(appErr)
		return appErr
	}

	filter := bson.M{"chat_id": user.ChatID}
	res := ufr.collection.FindOne(ctx, filter)

	if res.Err() != nil && !errors.Is(res.Err(), mongo.ErrNoDocuments) {
		appErr := apperrors.UFRSaveUserIfNotExistErr.AppendMessage(fmt.Sprintf("Cannot find user by chat_id. Err: %+v", res.Err()))
		ufr.log.Error(appErr)
		return appErr
	}

	if res.Err() == nil {
		appErr := apperrors.UFRSaveUserIfNotExistErr.AppendMessage(fmt.Sprintf("The user already exists in the database with ChatID: %v", user.ChatID))
		ufr.log.Error(appErr)
		return nil
	}

	_, err := ufr.collection.InsertOne(ctx, user)
	if err != nil {
		appErr := apperrors.UFRSaveUserIfNotExistErr.AppendMessage(err)
		ufr.log.Error(appErr)
		return appErr
	}

	ufr.log.Info("The user has been added, successfully.")
	return nil
}

func (ufr *userFastRepo) GetAllUsers(ctx context.Context) ([]*models.User, error) {
	cursor, err := ufr.collection.Find(ctx, bson.M{})
	if err != nil {
		appErr := apperrors.UFRGetAllUsersErr.AppendMessage(err)
		ufr.log.Error(appErr)
		return nil, appErr
	}

	defer cursor.Close(ctx)

	return decodeUsers(ctx, cursor)
}

func (ufr *userFastRepo) UpdateModification(ctx context.Context, user *models.User) error {
	filter := bson.M{"chat_id": user.ChatID}
	update := bson.M{
		"$set": bson.M{
			"updated": user.Updated,
		},
	}

	res, err := ufr.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		appErr := apperrors.UFRUpdateModificationErr.AppendMessage(err)
		ufr.log.Error(appErr)
		return appErr
	}

	if res.ModifiedCount == 0 {
		appErr := apperrors.UFRUpdateModificationErr.AppendMessage("No documents were updated")
		ufr.log.Error(appErr)
		return appErr
	}

	return nil
}

func decodeUsers(ctx context.Context, cursor *mongo.Cursor) ([]*models.User, error) {
	defer cursor.Close(ctx)

	var users []*models.User
	for cursor.Next(ctx) {
		var user models.User
		err := cursor.Decode(&user)
		if err != nil {
			appErr := apperrors.UFRdecodeUsersErr.AppendMessage(err)
			return nil, appErr
		}

		users = append(users, &user)
	}

	if err := cursor.Err(); err != nil {
		appErr := apperrors.UFRdecodeUsersErr.AppendMessage(err)
		return nil, appErr
	}
	return users, nil
}
