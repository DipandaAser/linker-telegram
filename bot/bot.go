package bot

import (
	"fmt"
	"github.com/DipandaAser/linker"
	"github.com/DipandaAser/linker-telegram/app"
	"github.com/DipandaAser/linker-telegram/bot/groups"
	tb "github.com/DipandaAser/telegrambot"
	"github.com/gin-gonic/gin"
)

var myBot *tb.Bot
var webHook *tb.Webhook
var botCommands []tb.Command

func InitBot() error {

	webHook = tb.NewWebhook(tb.Webhook{
		Listen:   ":" + app.Config.HTTPPort,
		Endpoint: &tb.WebhookEndpoint{PublicURL: app.Config.WebUrl},
	})

	b, err := tb.NewBot(tb.Settings{
		Token:   app.TelegramBotToken,
		Poller:  webHook,
		Verbose: false,
	})

	if err != nil {
		return err
	}
	myBot = b

	botCmds := []tb.Command{}
	for _, command := range linker.GetCommands() {
		botCmds = append(botCmds, tb.Command{
			Text:        fmt.Sprintf("/%s", command.Text),
			Description: command.Description,
		})
	}
	botCommands = botCmds

	// Setting all bot commands
	_ = myBot.SetCommands(botCommands)
	return nil
}

func Start() {
	SetBotHandlers()
	myBot.Start()
}

func GetNoBotRouter() *gin.RouterGroup {
	return webHook.NoBotRouter
}

func GetNoBotEndpointPath() string {
	return webHook.GetNoBotEndpoint()
}

func SetBotHandlers() {

	//======COMMANDS HANDLERS====================================
	myBot.Handle("/help", botHelpCommandHandler)
	myBot.Handle("/config", botConfigCommandHandler)
	myBot.Handle("/link", botLinkCommandHandler)
	myBot.Handle("/diffuse", botDiffuseCommandHandler)
	myBot.Handle("/list", botListCommandHandler)

	//======MESSAGES HANDLERS====================================
	myBot.Handle(tb.OnText, botTextMessageHandler)
	myBot.Handle(tb.OnAudio, botAudioMessageHandler)
	myBot.Handle(tb.OnPhoto, botOnPhotoHandler)
	myBot.Handle(tb.OnVideo, botOnVideoHAndler)
	myBot.Handle(tb.OnDocument, botOnDocumentHandler)

	myBot.Handle(tb.OnAddedToGroup, func(m *tb.Message) {
		_, getErr := groups.VerifyGroupExistenceAndCreateIfNot(fmt.Sprintf("%d", m.Chat.ID))
		if getErr != nil {
			_, _ = myBot.Send(m.Chat, "It's like somethings wrong, try re-adding this bot to the group")
			return
		}
	})
}

//isUserAdmin check if a user is an admin of a groups
func isUserAdmin(group *tb.Chat, user *tb.User) bool {

	admins, err := myBot.AdminsOf(group)
	if err != nil {
		return false
	}

	for _, admin := range admins {
		if admin.User.ID == user.ID {
			return true
		}
	}

	return false
}
