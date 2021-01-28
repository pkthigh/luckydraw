package tgbot

import "github.com/go-telegram-bot-api/telegram-bot-api"



func Help(m *tgbotapi.Message) {
	chat := m.Chat


	reply := tgbotapi.NewMessage(chat.ID, "这是一个抽奖机器人")
	markup := tgbotapi.NewInlineKeyboardMarkup()

	cmds := make(map[string]string)
	if _ , ok := draws[chat.ID]; ok {
		cmds["抽奖"]="/luckydraw"

	} else {
		cmds["签到"] = "/sign"

	}

	for k, v := range cmds {
		button := tgbotapi.NewInlineKeyboardButtonData(k, v)
		row := tgbotapi.NewInlineKeyboardRow(button)
		markup.InlineKeyboard = append(markup.InlineKeyboard, row)
	}
	reply.ReplyMarkup = markup

	Bot.Send( reply )

}

//是否是管理员
func IsAdmin(u *tgbotapi.User) bool {
	if flag, ok := admins[u.String()]; ok && flag {
		return true
	}
	return false
}
