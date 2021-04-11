package bot

import (
	"fmt"
	"github.com/DipandaAser/linker"
	"github.com/DipandaAser/linker-telegram/app"
	"github.com/DipandaAser/linker-telegram/bot/groups"
	tb "github.com/DipandaAser/telegrambot"
	"strings"
)

func botTextMessageHandler(m *tb.Message) {
	if m.IsService() || m.Sender.IsBot {
		// we not transfer bot or service message
		return
	}
	if !m.FromGroup() {
		return
	}

	if strings.HasPrefix(m.Text, "/") {
		return
	}

	group, getErr := groups.VerifyGroupExistenceAndCreateIfNot(fmt.Sprintf("%d", m.Chat.ID))
	if getErr != nil {
		return
	}

	_ = group.IncrementMessage()

	go linkSend(group, m)
	go diffusionSend(group, m)
}

func botAudioMessageHandler(m *tb.Message) {
	if m.IsService() || m.Sender.IsBot {
		// we not transfer bot or service message
		return
	}
	if !m.FromGroup() {
		return
	}

	if strings.HasPrefix(m.Text, "/") {
		return
	}

	group, getErr := groups.VerifyGroupExistenceAndCreateIfNot(fmt.Sprintf("%d", m.Chat.ID))
	if getErr != nil {
		return
	}

	_ = group.IncrementMessage()

	go linkSend(group, m)
	go diffusionSend(group, m)
}

func botOnPhotoHandler(m *tb.Message) {
	if m.IsService() || m.Sender.IsBot {
		// we not transfer bot or service message
		return
	}
	if !m.FromGroup() {
		return
	}

	if strings.HasPrefix(m.Text, "/") {
		return
	}

	group, getErr := groups.VerifyGroupExistenceAndCreateIfNot(fmt.Sprintf("%d", m.Chat.ID))
	if getErr != nil {
		return
	}

	_ = group.IncrementMessage()

	go linkSend(group, m)
	go diffusionSend(group, m)
}

func botOnVideoHAndler(m *tb.Message) {
	if m.IsService() || m.Sender.IsBot {
		// we not transfer bot or service message
		return
	}
	if !m.FromGroup() {
		return
	}

	if strings.HasPrefix(m.Text, "/") {
		return
	}

	group, getErr := groups.VerifyGroupExistenceAndCreateIfNot(fmt.Sprintf("%d", m.Chat.ID))
	if getErr != nil {
		return
	}

	_ = group.IncrementMessage()

	go linkSend(group, m)
	go diffusionSend(group, m)
}

func botOnDocumentHandler(m *tb.Message) {
	if m.IsService() || m.Sender.IsBot {
		// we not transfer bot or service message
		return
	}
	if !m.FromGroup() {
		return
	}

	if strings.HasPrefix(m.Text, "/") {
		return
	}

	group, getErr := groups.VerifyGroupExistenceAndCreateIfNot(fmt.Sprintf("%d", m.Chat.ID))
	if getErr != nil {
		return
	}

	_ = group.IncrementMessage()

	go linkSend(group, m)
	go diffusionSend(group, m)

}

func linkSend(group *linker.Group, m *tb.Message) {

	links, err := linker.GetLinksByGroupID(group.ID)
	if err != nil {
		return
	}

	for _, link := range links {
		var otherGroupID string
		for _, id := range link.GroupsID {
			if id != group.ID {
				otherGroupID = id
			}
		}
		grp, err := linker.GetGroupByID(otherGroupID)
		if err != nil {
			continue
		}

		if grp.Service == app.Config.ServiceName {
			_, _ = myBot.Forward(CustomReceiver{id: otherGroupID}, m)
			_ = link.IncrementMessage()
			continue
		}

		// TODO implement send message to other service
	}
}

func diffusionSend(group *linker.Group, m *tb.Message) {

	diffusions, err := linker.GetDiffusionsByBroadcaster(group.ID)
	if err != nil {
		return
	}

	for _, diffusion := range diffusions {
		grp, err := linker.GetGroupByID(diffusion.Receiver)
		if err != nil {
			continue
		}

		if grp.Service == app.Config.ServiceName {
			_, _ = myBot.Forward(CustomReceiver{id: grp.ID}, m)
			_ = diffusion.IncrementMessage()
			continue
		}

		// TODO implement send message to other service
	}
}
