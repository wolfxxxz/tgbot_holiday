package telegbot

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"
	"weekend_bot/internal/apperrors"
	"weekend_bot/internal/config"
	"weekend_bot/internal/models"
	"weekend_bot/internal/repositories"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
)

type Bot struct {
	botClient         *tgbotapi.BotAPI
	log               *logrus.Logger
	postUsers         repositories.UserFastRepo
	postDate          repositories.FastRepo
	holidayUsers      repositories.UserHolidayRepo
	holidayDate       repositories.HolidayRepo
	timeoutMongoQuery int
}

func NewBot(config *config.Config, httpClient *http.Client, postRepo repositories.FastRepo,
	postUserRepo repositories.UserFastRepo, log *logrus.Logger, holidayUsers repositories.UserHolidayRepo,
	holidayDate repositories.HolidayRepo) (*Bot, error) {
	client, err := tgbotapi.NewBotAPIWithClient(config.Token, "https://api.telegram.org/bot%s/%s", httpClient)
	if err != nil {
		appErr := apperrors.NewBotErr.AppendMessage(err)
		log.Error(appErr)
		return nil, appErr
	}

	client.Token = config.Token
	log.Infof("Authorized on account %s", client.Self.UserName)
	timeoutMongoQuery, err := strconv.Atoi(config.TimeoutMongoQuery)
	if err != nil {
		appErr := apperrors.NewBotErr.AppendMessage(err)
		log.Error(appErr)
		return nil, appErr
	}
	return &Bot{botClient: client, log: log, postDate: postRepo, postUsers: postUserRepo,
		holidayUsers: holidayUsers, holidayDate: holidayDate, timeoutMongoQuery: timeoutMongoQuery}, nil
}

func (bot *Bot) ReplyingOnMessages(ctx context.Context) error {
	updateConfig := tgbotapi.NewUpdate(allAvailableUpdates)
	updateConfig.Timeout = expectAnswerSec
	updates := bot.botClient.GetUpdatesChan(updateConfig)

	for update := range updates {
		answerMessage := bot.replyOnNewMessage(ctx, &update)
		if answerMessage == nil {
			appErr := apperrors.ReplyingOnMessagesErr.AppendMessage("somebody sent geotranslation")
			bot.log.Error(appErr)
			continue
		}

		_, err := bot.botClient.Send(answerMessage)
		if err != nil {
			appErr := apperrors.ReplyingOnMessagesErr.AppendMessage(fmt.Sprintf("Error sending message: %v\n", err))
			bot.log.Error(appErr)
		}
	}

	return nil
}

func (bot *Bot) replyOnNewMessage(ctx context.Context, upd *tgbotapi.Update) *tgbotapi.MessageConfig {
	if upd.Message == nil {
		appErr := apperrors.ReplyOnNewMessageErr.AppendMessage("upd = nil")
		bot.log.Error(appErr)
		return nil
	}

	chatID := upd.Message.Chat.ID
	text := upd.Message.Text
	bot.log.Infof("Replying on message. Text: %s; ChatID %v\n", text, chatID)

	var reply string
	var keyboard tgbotapi.ReplyKeyboardMarkup
	switch text {
	case subscribeOnPist:
		reply = bot.subscribePist(ctx, chatID, upd)
	case subscribeOnNextHolliday:
		reply = bot.subscribeHoliday(ctx, chatID, upd)
	case unSubscribeOnPist:
		reply = bot.unSubscribePist(ctx, chatID, upd)
	case unSubscribeOnNextHolliday:
		reply = bot.unSubscribeHoliday(ctx, chatID, upd)
	default:
		reply = "Виберіть що саме вас цікавить"
	}

	keyboard = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton(subscribeOnPist), tgbotapi.NewKeyboardButton(subscribeOnNextHolliday)),
		tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton(unSubscribeOnPist), tgbotapi.NewKeyboardButton(unSubscribeOnNextHolliday)),
	)

	message := tgbotapi.NewMessage(chatID, reply)
	message.ReplyMarkup = keyboard
	return &message
}

func (bot *Bot) subscribePist(ctx context.Context, chatID int64, upd *tgbotapi.Update) string {
	user := models.Create(chatID)
	bot.log.Info(user)

	ctxWithTimeout, cancel := context.WithTimeout(ctx, time.Second*time.Duration(bot.timeoutMongoQuery))
	defer cancel()

	err := bot.postUsers.SaveUserIfNotExist(ctxWithTimeout, user)
	if err != nil {
		appErr := err.(*apperrors.AppError)
		bot.log.Error(appErr)
		reply := "Шось пішло не так винен провайдер"
		return reply
	}

	reply := "ви підписались на розсилку Піст"
	return reply
}

func (bot *Bot) unSubscribePist(ctx context.Context, chatID int64, upd *tgbotapi.Update) string {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, time.Second*time.Duration(bot.timeoutMongoQuery))
	defer cancel()

	err := bot.postUsers.DropUser(ctxWithTimeout, chatID)
	if err != nil {
		appErr := err.(*apperrors.AppError)
		bot.log.Error(appErr)
		reply := "Шось пішло не так винен провайдер"
		return reply
	}

	reply := "ви відписались від розсилки Піст"
	return reply
}

func (bot *Bot) subscribeHoliday(ctx context.Context, chatID int64, upd *tgbotapi.Update) string {
	user := models.Create(chatID)

	ctxWithTimeout, cancel := context.WithTimeout(ctx, time.Second*time.Duration(bot.timeoutMongoQuery))
	defer cancel()

	err := bot.holidayUsers.SaveUserIfNotExist(ctxWithTimeout, user)
	if err != nil {
		appErr := err.(*apperrors.AppError)
		bot.log.Error(appErr)
		reply := "Шось пішло не так винен провайдер"
		return reply
	}

	reply := "ви підписались на розсилку Святкові дні"
	return reply
}

func (bot *Bot) unSubscribeHoliday(ctx context.Context, chatID int64, upd *tgbotapi.Update) string {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, time.Second*time.Duration(bot.timeoutMongoQuery))
	defer cancel()

	err := bot.holidayUsers.DropUser(ctxWithTimeout, chatID)
	if err != nil {
		appErr := err.(*apperrors.AppError)
		bot.log.Error(appErr)
		reply := "Шось пішло не так винен провайдер"
		return reply
	}

	reply := "ви відписались від розсилки Святкові дні"
	return reply
}

func (bot *Bot) SendingNotificationsFoodInPost(ctx context.Context, hourRing, timeOutMinute int) error {
	ticker := time.NewTicker(time.Minute * time.Duration(timeOutMinute))
	defer ticker.Stop()

	for t := range ticker.C {
		if t.Hour() == hourRing {
			if err := bot.pushScheduledFist(ctx, t); err != nil {
				bot.log.Error(err)
				if err.Error() != "Forbidden: bot was blocked by the user" {
					return err
				}
			}

			if err := bot.pushScheduleHoliday(ctx, t); err != nil {
				bot.log.Error(err)
				if err.Error() != "Forbidden: bot was blocked by the user" {
					return err
				}
			}
		}
	}

	return nil
}

func (bot *Bot) pushScheduledFist(ctx context.Context, executedTime time.Time) error {
	users, err := bot.postUsers.GetAllUsers(ctx)
	if err != nil {
		appErr := err.(*apperrors.AppError)
		bot.log.Error(appErr)
		return appErr
	}

	for _, user := range users {
		time.Sleep(time.Second * 1)
		reply := ""
		weekend, err := bot.postDate.GetTodaysDate(ctx)
		if err != nil {
			appErr := err.(*apperrors.AppError)
			bot.log.Error(appErr)
			reply = "Something went wrong, provider is quilty"
			return appErr
		}

		reply, err = MapModelsPostToResponse(weekend)
		if err != nil {
			appErr := err.(*apperrors.AppError)
			bot.log.Error(appErr)
			reply = "Something went wrong, noone is to blame. Just a mistake"
			return appErr
		}

		message := tgbotapi.NewMessage(user.ChatID, reply)
		bot.log.Infof("Sending a reply to user. Reply: %+v; ChatID: %v;", reply, user.ChatID)
		_, err = bot.botClient.Send(message)
		if err != nil {
			appErr := apperrors.PushScheduledFistErr.AppendMessage(fmt.Sprintf("Error sending message: %v\n", err))
			bot.log.Error(appErr)
			return appErr
		}

		user.Update()
		ctxWithTimeout, cancel := context.WithTimeout(ctx, time.Second*time.Duration(bot.timeoutMongoQuery))
		defer cancel()
		if err = bot.postUsers.UpdateModification(ctxWithTimeout, user); err != nil {
			appErr := err.(*apperrors.AppError)
			bot.log.Error(appErr)
			return appErr
		}
	}

	bot.log.Infof("Execute a scheduled task at %v sent users %v", executedTime.Format("15:04"), len(users))
	return nil
}

func (bot *Bot) pushScheduleHoliday(ctx context.Context, executedTime time.Time) error {
	users, err := bot.holidayUsers.GetAllUsers(ctx)
	if err != nil {
		appErr := err.(*apperrors.AppError)
		bot.log.Error(appErr)
		return err
	}

	for _, user := range users {
		time.Sleep(time.Second * 1)
		reply := ""
		weekend, err := bot.holidayDate.GetTodaysDate(ctx)
		if err != nil {
			appErr := err.(*apperrors.AppError)
			bot.log.Error(appErr)
			reply = "Something went wrong, provider is quilty"
			return nil
		}

		bot.log.Infof("weekend: %+v", weekend)
		reply, err = MapModelsHolidayToResponse(weekend)
		if err != nil {
			appErr := err.(*apperrors.AppError)
			bot.log.Error(appErr)
			reply = "Something went wrong, noone is to blame. Just a mistake"
			return appErr
		}

		message := tgbotapi.NewMessage(user.ChatID, reply)
		bot.log.Infof("Sending a reply to user. Reply: %+v; ChatID: %v;", reply, user.ChatID)
		_, err = bot.botClient.Send(message)
		if err != nil {
			appErr := apperrors.PushScheduleHolidayErr.AppendMessage(fmt.Sprintf("Error sending message: %v\n", err))
			bot.log.Error(appErr)
			return appErr
		}

		user.Update()
		ctxWithTimeout, cancel := context.WithTimeout(ctx, time.Second*time.Duration(bot.timeoutMongoQuery))
		defer cancel()
		if err = bot.postUsers.UpdateModification(ctxWithTimeout, user); err != nil {
			appErr := err.(*apperrors.AppError)
			bot.log.Error(appErr)
			return appErr
		}
	}

	bot.log.Infof("Execute a scheduled task at %v sent users %v", executedTime.Format("15:04"), len(users))
	return nil
}
