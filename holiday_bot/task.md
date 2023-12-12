/*
		if res != nil && res.Err() != nil {
			if errors.Is(res.Err(), mongo.ErrNoDocuments) {
				return apperrors.MongoDataExistsError.AppendMessage(res.Err())
			}

			ur.log.Errorf("Cannot find user by chat_id. Err: %+v", res.Err())
			return apperrors.MongoGetFailedError.AppendMessage(res.Err())
		}
	*/

	/*
		defer cursor.Close(ctx)

		var users []*models.User
		for cursor.Next(ctx) {
			var user models.User
			err := cursor.Decode(&user)
			if err != nil {
				return nil, apperrors.MongoGetFailedError.AppendMessage(err)
			}

			users = append(users, &user)
		}

		if err := cursor.Err(); err != nil {
			return nil, apperrors.MongoGetFailedError.AppendMessage(err)
		}
	*/

	/*
	test

func TestUserRepo_GetAllUsers(t *testing.T) {
	// Создаем макет для коллекции
	collection := new(MockCollection)

	// Создаем UserRepo, используя макет
	userRepo := &UserRepo{collection: collection}

	// Фейковый контекст
	ctx := context.Background()

	// Модели пользователей для симуляции результатов запроса
	user1 := &models.User{
		ChatID: 1,
		// Другие поля пользователя
	}
	user2 := &models.User{
		ChatID: 2,
		// Другие поля пользователя
	}

	// Ожидаемый результат запроса, симулирующий наличие пользователей
	cursor := new(MockCursor)
	cursor.On("Next", ctx).Return(true)
	cursor.On("Next", ctx).Return(true)
	cursor.On("Next", ctx).Return(false)
	cursor.On("Decode", mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		user := args.Get(0).(*models.User)
		*user = *user1 // Заполняем пользователями для симуляции декодирования
	}).Once()
	cursor.On("Decode", mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		user := args.Get(0).(*models.User)
		*user = *user2 // Заполняем пользователями для симуляции декодирования
	}).Once()
	collection.On("Find", ctx, mock.Anything, mock.Anything).Return(cursor, nil)

	// Вызываем метод GetAllUsers
	users, err := userRepo.GetAllUsers(ctx)

	// Проверяем, что ошибки нет
	assert.NoError(t, err)

	// Проверяем, что вернуты ожидаемые пользователи
	assert.Len(t, users, 2)
	assert.Equal(t, user1, users[0])
	assert.Equal(t, user2, users[1])
}

func TestUserRepo_UpdateModification(t *testing.T) {
	// Создаем макет для коллекции
	collection := new(MockCollection)

	// Создаем UserRepo, используя макет
	userRepo := &UserRepo{collection: collection}

	// Фейковый контекст и пользователь для обновления
	ctx := context.Background()
	user := &models.User{
		ChatID:  12345,
		Updated: time.Now(),
		// Другие поля пользователя
	}

	// Устанавливаем ожидание вызова метода UpdateOne для успешного обновления
	collection.On("UpdateOne", ctx, mock.Anything, mock.Anything, mock.Anything).Return(
		&mongo.UpdateResult{ModifiedCount: 1}, nil,
	)

	// Вызываем метод UpdateModification
	err := userRepo.UpdateModification(ctx, user)

	// Проверяем, что ошибки нет
	assert.NoError(t, err)

	// Устанавливаем ожидание вызова метода UpdateOne для неудачного обновления
	collection.On("UpdateOne", ctx, mock.Anything, mock.Anything, mock.Anything).Return(
		&mongo.UpdateResult{ModifiedCount: 0}, nil,
	)

	// Вызываем метод UpdateModification, чтобы симулировать неудачное обновление
	err = userRepo.UpdateModification(ctx, user)

	// Проверяем, что вернулась ошибка, связанная с неудачным обновлением
	assert.Error(t, err)
}
*/

/*if res.Err() != nil && !errors.Is(res.Err(), mongo.ErrNoDocuments) {
		ur.log.Errorf("Cannot find user by chat_id. Err: %+v", res.Err())
		return apperrors.MongoGetFailedError.AppendMessage(res.Err())
	}

	if errors.Is(res.Err(), mongo.ErrNoDocuments) {
		ur.log.Infof("The user already exists in the database with ChatID: %v", user.ChatID)
		return apperrors.MongoDataExistsError.AppendMessage(res.Err())
	}

	_, err := ur.collection.InsertOne(ctx, user)
	if err != nil {
		return apperrors.MongoSaveUserFailedError.AppendMessage(err)
	}

	ur.log.Info("The user has been added, successfully.")
	return nil
	*/

	/*
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	mockUserRepo := mock.NewMockUserRepoInterface(ctl)

	//mockUserRepo.EXPECT().GetAllUsers(gomock.Any()).Return([]*models.User{}, nil)
	//mockUserRepo.EXPECT().SaveUserIfNotExist(gomock.Any(), gomock.Any()).Return(nil)
	//mockUserRepo.EXPECT().UpdateModification(gomock.Any(), gomock.Any()).Return(nil)
	*/

	/*
	

/*
	type testTransport2 struct {
		responses []*http.Response
		index     int
	}

	func (t *testTransport2) RoundTrip(req *http.Request) (*http.Response, error) {
		if t.index >= len(t.responses) {
			return nil, errors.New("no more responses")
		}

		response := t.responses[t.index]
		t.index++
		return response, nil
	}

	func fakeHTTPBotClientWithMultipleResponses(responses []*http.Response) *http.Client {
		return &http.Client{
			Transport: &testTransport2{
				responses: responses,
				index:     0,
			},
		}
	}

	func fakeBotWithWeatherClientMultipleResponses(weatherClient *httpclient.WeatherClient, responses []*http.Response) *Bot {
		apiToken := os.Getenv("TOKEN")
		testConfig := &config.Config{
			Token: apiToken,
		}

		testTgClientHttp := fakeHTTPBotClientWithMultipleResponses(responses)
		db, err := database.InitClient(context.TODO(), testConfig, logger)
		if err != nil {
			logger.Debug(err)
		}

		userRepo := repositories.NewUserRepo(testConfig, logger, db)
		bot, err := NewBot(testConfig, testTgClientHttp, weatherClient, logger, userRepo)
		if err != nil {
			logger.Fatal(err)
			return nil
		}

		return bot
	}
*/

/*
	func TestGetMessageByUpdate(t *testing.T) {
		ttPass := []struct {
			name                 string
			messageChatId        int64
			givenWeatherResponse *httpclient.GetWeatherResponse
			givenMessage         *tgbotapi.Update
			botReply             *tgbotapi.MessageConfig
		}{
			{
				"existing location command",
				300,
				&httpclient.GetWeatherResponse{
					Name: "Barselona",
					Main: &httpclient.Main{
						Temp:      10,
						FeelsLike: 15,
					},
					Weather: []*httpclient.Weather{
						{Description: "holly crap"},
					},
				},
				&tgbotapi.Update{
					UpdateID: 0,
					Message: &tgbotapi.Message{
						Chat: &tgbotapi.Chat{
							ID: 300,
						},
						Location: &tgbotapi.Location{
							Longitude: 21.017532,
							Latitude:  52.237049,
						},
					},
				},
				&tgbotapi.MessageConfig{
					Text:      "gh",
					ParseMode: "HTML",
				},
			},
		}

		testConfig := config.NewConfig()
		err := testConfig.ParseConfig(".env", logger)
		if err != nil {
			logger.Fatal(err)
		}

		for _, tc := range ttPass {

			responseJSON, err := json.Marshal(tc.givenWeatherResponse)
			if err != nil {
				t.Fatal(err)
			}

			httpClient := fakeHTTPBotClient(200, string(responseJSON))
			weatherClient := httpclient.NewWeatherClient(testConfig, httpClient, logger)
			bot := fakeBotWithWeatherClient(weatherClient)
			bot.botClient.Debug = true
			msg := bot.replyOnNewMessage(context.TODO(), tc.givenMessage)
			if err != nil {
				t.Error(err)
			}

			expMessage, err := MapGetWeatherResponseToWeatherAnswer(tc.givenWeatherResponse)
			if err != nil {
				t.Error(err)
			}
			if msg != nil && msg.Text != expMessage {
				t.Errorf("bot reply should be %s, but got %s", tc.botReply.Text, msg.Text)
			}
		}
	}

	func fakeBotWithWeatherClient(weatherClient *httpclient.WeatherClient) *Bot {
		apiToken := os.Getenv("TOKEN")
		testConfig := &config.Config{
			Token: apiToken,
		}

		response, err := generateBotOkJsonApiResponse()
		if err != nil {
			logger.Fatal(err)
			return nil
		}

		testTgClientHttp := fakeHTTPBotClient(200, response)
		tgClient, err := tgbotapi.NewBotAPIWithClient(testConfig.Token, "https://api.telegram.org/bot%s/%s", testTgClientHttp)
		if err != nil {
			logger.Fatal(err)
		}

		db, err := database.InitClient(context.TODO(), testConfig, logger)
		if err != nil {
			logger.Debug(err)
		}

		usrRepo := repositories.NewUserRepo(testConfig, logger, db)
		bot, err := NewBot(testConfig, testTgClientHttp, weatherClient, logger, usrRepo)
		if err != nil {
			logger.Fatal(err)
			return nil
		}
		bot.botClient = tgClient
		return bot
	}
*/

/*
func TestGetMessageByUpdate(t *testing.T) {
	ttPass := []struct {
		name                 string
		messageChatId        int64
		givenWeatherResponse *httpclient.GetWeatherResponse
		givenMessage         *tgbotapi.Update
		botReply             *tgbotapi.MessageConfig
	}{
		{
			"existing location command",
			300,
			&httpclient.GetWeatherResponse{
				Name: "Barselona",
				Main: &httpclient.Main{
					Temp:      10,
					FeelsLike: 15,
				},
				Weather: []*httpclient.Weather{
					{Description: "holly crap"},
				},
			},
			&tgbotapi.Update{
				UpdateID: 0,
				Message: &tgbotapi.Message{
					Chat: &tgbotapi.Chat{
						ID: 300,
					},
					Location: &tgbotapi.Location{
						Longitude: 21.017532,
						Latitude:  52.237049,
					},
				},
			},
			&tgbotapi.MessageConfig{
				Text:      "gh",
				ParseMode: "HTML",
			},
		},
	}

	testConfig := config.NewConfig()
	err := testConfig.ParseConfig(".env", logger)
	if err != nil {
		logger.Fatal(err)
	}

	for _, tc := range ttPass {

		responseJSON, err := json.Marshal(tc.givenWeatherResponse)
		if err != nil {
			t.Fatal(err)
		}

		httpClient := fakeHTTPBotClient(200, string(responseJSON))
		weatherClient := httpclient.NewWeatherClient(testConfig, httpClient, logger)
		bot := fakeBotWithWeatherClient(weatherClient, t)
		bot.botClient.Debug = true
		msg := bot.replyOnNewMessage(context.TODO(), tc.givenMessage)
		if err != nil {
			t.Error(err)
		}

		expMessage, err := MapGetWeatherResponseToWeatherAnswer(tc.givenWeatherResponse)
		if err != nil {
			t.Error(err)
		}
		if msg != nil && msg.Text != expMessage {
			t.Errorf("bot reply should be %s, but got %s", tc.botReply.Text, msg.Text)
		}
	}
}

func fakeBotWithWeatherClient(weatherClient *httpclient.WeatherClient, t *testing.T) *Bot {
	apiToken := os.Getenv("TOKEN")
	testConfig := &config.Config{
		Token:             apiToken,
		TimeoutMongoQuery: "5",
	}

	response, err := generateBotOkJsonApiResponse()
	if err != nil {
		logger.Fatal(err)
		return nil
	}

	testTgClientHttp := fakeHTTPBotClient(200, response)
	tgClient, err := tgbotapi.NewBotAPIWithClient(testConfig.Token, "https://api.telegram.org/bot%s/%s", testTgClientHttp)
	if err != nil {
		logger.Fatal(err)
	}

	ctl := gomock.NewController(t)
	defer ctl.Finish()

	mockUserRepo := mock.NewMockUserRepoInterface(ctl)
	mockUserRepo.EXPECT().SaveUserIfNotExist(gomock.Any(), gomock.Any()).Return(nil)
	bot, err := NewBot(testConfig, testTgClientHttp, weatherClient, logger, mockUserRepo)
	if err != nil {
		logger.Fatal(err)
		return nil
	}
	bot.botClient = tgClient
	return bot
}
*/

ERRO[247505]/home/mvmir/Рабочий стол/Святковий/weekend_bot/cmd/main.go:55 main.main.func1() Forbidden: bot was kicked from the supergroup chat 
2023/12/06 10:33:00 Post "https://api.telegram.org/bot6593825661:AAE-uCjzVTPMsAdvH4W9Bf_9pTaxiB8JtQE/getUpdates": read tcp 192.168.3.111:59220->149.154.167.220:443: read: connection reset by peer
2023/12/06 10:33:00 Failed to get updates, retrying in 3 seconds...
