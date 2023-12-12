package apperrors

import "fmt"

type AppError struct {
	Message string
	Code    string
}

func NewAppError() *AppError {
	return &AppError{}
}

var (
	EnvConfigLoadError = AppError{
		Message: "Failed to load env file",
		Code:    EnvInit,
	}
	EnvConfigParseError = AppError{
		Message: "Failed to parse env file",
		Code:    EnvParse,
	}
	BotInitializationError = AppError{
		Message: "Failed to init new bot",
		Code:    BotInit,
	}
	BotSendMessageError = AppError{
		Message: "Failed to send the message",
		Code:    BotSendMsg,
	}
	WeatherClientError = AppError{
		Message: "Failed send request",
		Code:    HttpSendRequest,
	}
	BotMapperEncodingErr = AppError{
		Message: "Failed encoding mapper",
		Code:    MapperEncoding,
	}
	MongoDataExistsError = AppError{
		Message: "Failed send mongoDB",
		Code:    UserRepo,
	}
	MongoSaveUserFailedError = AppError{
		Message: "Failed save mongoDB",
		Code:    UserRepo,
	}
	MongoGetFailedError = AppError{
		Message: "Failed Get mongoDB",
		Code:    UserRepo,
	}
	MongoUpdateModFailedError = AppError{
		Message: "Failed Update mongoDB",
		Code:    UserRepo,
	}
	MongoInitFailedError = AppError{
		Message: "Failed Init mongoDB",
		Code:    InitMongo,
	}
	MongoDropUserByIDErr = AppError{
		Message: "Failed Delete user",
		Code:    UserRepo,
	}
	TxtRepoErr = AppError{
		Message: "Failed Delete user",
		Code:    TxtRepo,
	}
	ModelsCreateFastErr = AppError{
		Message: "Failed Create fast",
		Code:    Models,
	}
	ConvertDateTimeErr = AppError{
		Message: "Failed ConvertDateTime",
		Code:    Models,
	}
	ModelsCreateHolidayErr = AppError{
		Message: "Failed Create holiday",
		Code:    Models,
	}
	MongoSaveAllDatesError = AppError{
		Message: "Failed SaveAllDatesErr",
		Code:    FastRepo,
	}
	GetTodaysDate = AppError{
		Message: "Failed SaveAllDatesErr",
		Code:    FastRepo,
	}
	SaveDateIfNotExist = AppError{
		Message: "Failed SaveDateIfNotExistErr",
		Code:    FastRepo,
	}
	GetAllDatesErr = AppError{
		Message: "Failed GetAllDatesErr",
		Code:    FastRepo,
	}
	DecodeDatesErr = AppError{
		Message: "Failed decodeDatesErr",
		Code:    FastRepo,
	}
	HolidaySaveAllDatesErr = AppError{
		Message: "Failed SaveAllDatesErr",
		Code:    HolidayRepo,
	}
	HolidayGetTodaysDateErr = AppError{
		Message: "Failed GetTodaysDateErr",
		Code:    HolidayRepo,
	}
	HolidaySaveDateIfNotExistErr = AppError{
		Message: "Failed SaveDateIfNotExistErr",
		Code:    HolidayRepo,
	}
	HolidayGetAllDatesErr = AppError{
		Message: "Failed GetAllDatesErr",
		Code:    HolidayRepo,
	}
	UpdateModFailedError = AppError{
		Message: "Failed UpdateModFailedError",
		Code:    HolidayRepo,
	}
	DecodeHolidaysErr = AppError{
		Message: "Failed DecodeHolidaysErr",
		Code:    HolidayRepo,
	}
	UserFastDropUserErr = AppError{
		Message: "Failed DropUser",
		Code:    UserFastRepo,
	}
	UFRSaveUserIfNotExistErr = AppError{
		Message: "Failed SaveUserIfNotExist",
		Code:    UserFastRepo,
	}
	UFRGetAllUsersErr = AppError{
		Message: "Failed GetAllUsers",
		Code:    UserFastRepo,
	}
	UFRUpdateModificationErr = AppError{
		Message: "Failed UpdateModification",
		Code:    UserFastRepo,
	}
	UFRdecodeUsersErr = AppError{
		Message: "Failed decodeUsers",
		Code:    UserFastRepo,
	}
	UHRDropUserErr = AppError{
		Message: "Failed DropUser",
		Code:    UserHolidayRepo,
	}
	UHRSaveUserIfNotExistErr = AppError{
		Message: "Failed SaveUserIfNotExist",
		Code:    UserHolidayRepo,
	}
	UHRGetAllUsersErr = AppError{
		Message: "Failed GetAllUsers",
		Code:    UserHolidayRepo,
	}
	UHRUpdateModificationErr = AppError{
		Message: "Failed UpdateModification",
		Code:    UserHolidayRepo,
	}
	NewBotErr = AppError{
		Message: "Failed NewBot",
		Code:    Bot,
	}
	ReplyingOnMessagesErr = AppError{
		Message: "Failed ReplyingOnMessagesErr",
		Code:    Bot,
	}
	ReplyOnNewMessageErr = AppError{
		Message: "Failed ReplyingOnMessageErr",
		Code:    Bot,
	}
	PushScheduledFistErr = AppError{
		Message: "Failed PushScheduleFistErr",
		Code:    Bot,
	}
	PushScheduleHolidayErr = AppError{
		Message: "Failed PushScheduleHolidayErr",
		Code:    Bot,
	}
)

func (appError *AppError) Error() string {
	return appError.Code + ": " + appError.Message
}

func (appError *AppError) AppendMessage(anyErrs ...interface{}) *AppError {
	return &AppError{
		Message: fmt.Sprintf("%v : %v", appError.Message, anyErrs),
		Code:    appError.Code,
	}
}

func IsAppError(err1 error, err2 *AppError) bool {
	err, ok := err1.(*AppError)
	if !ok {
		return false
	}

	return err.Code == err2.Code
}
