package database

import (
	"context"
	"fmt"
	"holiday_bot/internal/apperrors"
	"holiday_bot/internal/config"
	"time"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const timeOut = 5

func InitClient(ctx context.Context, conf *config.Config, log *logrus.Logger) (*mongo.Database, error) {
	var clientOpt *options.ClientOptions
	mongoDBURL := fmt.Sprintf("mongodb://%s:%s@%s:%s", conf.UserName, conf.Password, conf.MongoHost, conf.MongoPort)
	credential := options.Credential{
		Username: conf.UserName,
		Password: conf.Password,
	}

	ctx, cancel := context.WithTimeout(ctx, time.Second*time.Duration(timeOut))
	defer cancel()

	clientOpt = options.Client().ApplyURI(mongoDBURL).SetAuth(credential)
	client, err := mongo.Connect(ctx, clientOpt)
	if err != nil {
		appErr := apperrors.MongoInitFailedError.AppendMessage(fmt.Sprintf("error to connect mongoDB [%v]", err))
		log.Info(appErr)
		return nil, appErr
	}

	if err = client.Ping(ctx, nil); err != nil {
		appErr := apperrors.MongoInitFailedError.AppendMessage(fmt.Sprintf("error to ping mongoDB [%v]", err))
		return nil, appErr
	}

	log.Info("Ping success")

	mongoDB := client.Database(conf.DBName)
	return mongoDB, nil
}
