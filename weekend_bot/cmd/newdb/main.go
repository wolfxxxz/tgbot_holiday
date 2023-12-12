package main

import (
	"context"
	"weekend_bot/internal/config"
	"weekend_bot/internal/database"
	"weekend_bot/internal/log"
	"weekend_bot/internal/repositories"
)

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

	ctx := context.Background()
	mongoDB, err := database.InitClient(ctx, conf, logger)
	if err != nil {
		logger.Error(err)
		return
	}

	fastRepo := repositories.NewFastRepo(conf, logger, mongoDB)
	logger.Info("START DECODE fasts")
	txtRepo := repositories.NewTxtRepo(logger)
	pists, err := txtRepo.GetFasts("fasts.txt")
	if err != nil {
		logger.Error(err)
		return
	}

	err = fastRepo.SaveAllDates(ctx, pists)
	if err != nil {
		logger.Error(err)
		return
	}

	holidayRepo := repositories.NewHolidayRepo(conf, logger, mongoDB)
	logger.Info("START DECODE holiday")

	holidays, err := txtRepo.GetHolidays("holidays.txt")
	if err != nil {
		logger.Error(err)
		return
	}

	err = holidayRepo.SaveAllDates(ctx, holidays)
	if err != nil {
		logger.Error(err)
		return
	}

	logger.Info("holidays and fasts saved successfully")
}
