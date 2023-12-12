package main

import (
	"context"
	"net/http"
	"weekend_bot/internal/config"
	"weekend_bot/internal/database"
	"weekend_bot/internal/log"
	"weekend_bot/internal/repositories"
	"weekend_bot/internal/telegbot"
)

var sendWeatherForecast = 6 //12 //hours
var timeOut = 55            //minuts

func main() {

	logger, err := log.NewLogAndSetLevel("info")
	if err != nil {
		logger.Error(err)
		return
	}

	conf := config.NewConfig()
	err = conf.ParseConfig(".env", logger)
	if err != nil {
		logger.Error(err)
		return
	}

	if err = log.SetLevel(logger, conf.LogLevel); err != nil {
		logger.Error(err)
		return
	}

	httpClient := &http.Client{}

	ctx := context.Background()
	mongoDB, err := database.InitClient(ctx, conf, logger)
	if err != nil {
		logger.Error(err)
		return
	}

	userPostRepo := repositories.NewUserPostRepo(conf, logger, mongoDB)
	postRepo := repositories.NewFastRepo(conf, logger, mongoDB)

	userHolidayRepo := repositories.NewUserHolidayRepo(conf, logger, mongoDB)
	holidayRepo := repositories.NewHolidayRepo(conf, logger, mongoDB)

	bot, err := telegbot.NewBot(conf, httpClient, postRepo, userPostRepo, logger, userHolidayRepo, holidayRepo)
	if err != nil {
		logger.Error(err)
		return
	}

	go func() {
		if err := bot.SendingNotificationsFoodInPost(ctx, sendWeatherForecast, timeOut); err != nil {
			logger.Error(err)
		}
	}()

	logger.Info("BOT is replying on messages.")
	if err = bot.ReplyingOnMessages(ctx); err != nil {
		logger.Error(err)
	}

}
