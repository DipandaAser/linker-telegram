package app

import (
	"github.com/DipandaAser/linker"
	"os"
)

var Config *linker.ProjectSettings
var TelegramBotToken = os.Getenv("BOT_TOKEN")

func Init() {
	Config = &linker.ProjectSettings{}
	Config.ServiceName = "telegram"
	Config.ProjectName = "Linker Telegram"
	Config.AuthKey = os.Getenv("APIKEY")
	Config.DBName = os.Getenv("DB_NAME")
	Config.MongodbURI = os.Getenv("MONGO_URI")
	Config.HTTPPort = os.Getenv("PORT")
	Config.WebUrl = os.Getenv("WEB_URL")
}
