package bot

import (
	"fmt"
	"github.com/DipandaAser/linker"
	"github.com/DipandaAser/linker-telegram/bot/groups"
	tb "github.com/DipandaAser/telegrambot"
	"strings"
)

const (
	ERR_USE_CMD_IN_GROUP  = "Please use this command on a group where Linker is in"
	ERR_BOT_NOT_SUPPORTED = "Sorry but for now other bot can't interact with linker"
	ERR_GLOBAL            = "Something going wrong during the operation.\nPlease retry" +
		"\nIf it is persistent, contact @iamdipanda"
	ERR_INVALID_CODE = "Please provide two valid Linker Group Code"
	FOOTER           = "\nRepo: https://github.com/DipandaAser/linker. \nBy @iamdipanda"
)

func botHelpCommandHandler(m *tb.Message) {
	header := "Welcome to Bot Help" +
		"\nLinker is a bot who allows you to link group between platforms, and exchange messages" +
		"\nYou can link:" +
		"\n\t   telegram group <---> telegram group" +
		"\n\t   telegram group <---> whatsapp group"

	var commandMessage string = "Commands\n"

	for _, command := range botCommands {
		commandMessage += fmt.Sprintf("--> %s\nDescription: %s\n\n", command.Text, command.Description)
	}

	helpMessage := fmt.Sprintf("%s\n\n%s\n\n%s", header, commandMessage, FOOTER)
	if m.FromGroup() {
		// We send the message
		_, _ = myBot.Send(m.Sender, helpMessage, &tb.SendOptions{
			DisableWebPagePreview: true,
		})
		return
	} else {
		// We send the message
		_, _ = myBot.Send(m.Sender, helpMessage, &tb.SendOptions{
			DisableWebPagePreview: true,
		})
		return
	}
}

func botConfigCommandHandler(m *tb.Message) {

	if m.Sender.IsBot {
		_, _ = myBot.Send(m.Sender, ERR_BOT_NOT_SUPPORTED)
		return
	}
	if !m.FromGroup() {
		_, _ = myBot.Send(m.Sender, ERR_USE_CMD_IN_GROUP)
		return
	}

	group, lErr := groups.VerifyGroupExistenceAndCreateIfNot(fmt.Sprintf("%d", m.Chat.ID))
	if lErr != nil {
		_, _ = myBot.Send(m.Sender, ERR_GLOBAL)
		return
	}

	notAdminError := fmt.Sprintf("It's sound like your are not an admin of %s", m.Chat.Title)
	isAdmin := isUserAdmin(m.Chat, m.Sender)

	if isAdmin {
		_, _ = myBot.Send(m.Sender, fmt.Sprintf("Hey Dude this is your linker chat id of the %s group. \nLinker Group Code: %s", m.Chat.Title, group.ShortCode))
		return
	} else {
		_, _ = myBot.Send(m.Sender, notAdminError)
		return
	}

}

func botLinkCommandHandler(m *tb.Message) {
	if m.Sender.IsBot {
		_, _ = myBot.Send(m.Sender, ERR_BOT_NOT_SUPPORTED)
		return
	}

	payloadData := strings.Split(m.Payload, " ")
	var firstCode, secondCode string
	for _, id := range payloadData {
		if strings.TrimSpace(id) != "" {
			if firstCode == "" {
				firstCode = id
				continue
			}

			if secondCode == "" {
				secondCode = id
				break
			}
		}
	}

	if strings.TrimSpace(firstCode) == "" || strings.TrimSpace(secondCode) == "" {
		_, _ = myBot.Send(m.Sender, "Please provide two Linker Group Code")
		return
	}

	var firstGroup, secondGroup *linker.Group
	var lErr error
	if firstGroup, lErr = linker.GetGroupByShortCode(firstCode); lErr != nil {
		_, _ = myBot.Send(m.Sender, ERR_INVALID_CODE)
		return
	}
	if secondGroup, lErr = linker.GetGroupByShortCode(secondCode); lErr != nil {
		_, _ = myBot.Send(m.Sender, ERR_INVALID_CODE)
		return
	}

	_, lErr = linker.CreateLink([2]string{firstGroup.ID, secondGroup.ID})
	if lErr != nil {
		_, _ = myBot.Send(m.Sender, ERR_GLOBAL)
		return
	}

	_, _ = myBot.Send(m.Sender, "Link successfully created. \nYou can start exchange message between these Groups.")
	return
}

func botDiffuseCommandHandler(m *tb.Message) {

	if m.Sender.IsBot {
		_, _ = myBot.Send(m.Sender, ERR_BOT_NOT_SUPPORTED)
		return
	}

	payloadData := strings.Split(m.Payload, " ")
	var firstCode, secondCode string
	for _, id := range payloadData {
		if strings.TrimSpace(id) != "" {
			if firstCode == "" {
				firstCode = id
				continue
			}

			if secondCode == "" {
				secondCode = id
				break
			}
		}
	}

	if strings.TrimSpace(firstCode) == "" || strings.TrimSpace(secondCode) == "" {
		_, _ = myBot.Send(m.Sender, "Please provide two Linker Group Code")
		return
	}

	broadcaster, err := linker.GetGroupByShortCode(firstCode)
	if err != nil {
		_, _ = myBot.Send(m.Sender, ERR_INVALID_CODE)
		return
	}

	receiver, err := linker.GetGroupByShortCode(secondCode)
	if err != nil {
		_, _ = myBot.Send(m.Sender, ERR_INVALID_CODE)
		return
	}

	_, err = linker.CreateDiffusion(broadcaster.ID, receiver.ID)
	if err != nil {
		_, _ = myBot.Send(m.Sender, ERR_GLOBAL)
		return
	}

	_, _ = myBot.Send(m.Sender, "Diffusion successfully created. \nYou can start diffuse message in the second group/channel.")
	return
}

func botListCommandHandler(m *tb.Message) {

	if m.Sender.IsBot {
		_, _ = myBot.Send(m.Sender, ERR_BOT_NOT_SUPPORTED)
		return
	}
	if !m.FromGroup() {
		_, _ = myBot.Send(m.Sender, ERR_USE_CMD_IN_GROUP)
		return
	}

	notAdminError := fmt.Sprintf("It's sound like your are not an admin of %s", m.Chat.Title)
	isAdmin := isUserAdmin(m.Chat, m.Sender)
	if !isAdmin {
		_, _ = myBot.Send(m.Sender, notAdminError)
		return
	}

	group, lErr := groups.VerifyGroupExistenceAndCreateIfNot(fmt.Sprintf("%d", m.Chat.ID))
	if lErr != nil {
		_, _ = myBot.Send(m.Sender, ERR_GLOBAL)
		return
	}

	links, err := linker.GetLinksByGroupID(group.ID)
	if err != nil {
		_, _ = myBot.Send(m.Sender, ERR_GLOBAL)
		return
	}

	if len(links) == 0 {
		_, _ = myBot.Send(m.Sender, "The group don't have any active links or diffusion")
		return
	}

	messageLinks := fmt.Sprintf("Link list of %s\n\n", m.Chat.Title)
	for _, link := range links {
		var otherGroupID string
		for _, id := range link.GroupsID {
			if id != group.ID {
				otherGroupID = id
			}
		}
		grp, err := linker.GetGroupByID(otherGroupID)
		if err != nil {
			_, _ = myBot.Send(m.Sender, ERR_GLOBAL)
			return
		}
		messageLinks += fmt.Sprintf("\t--> %s on %s\n\n", grp.ShortCode, grp.Service)
	}

	_, _ = myBot.Send(m.Sender, messageLinks)
	return
}
