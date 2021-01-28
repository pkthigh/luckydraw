package tgbot

import (
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"luckydraw/log"
	"strings"
)
var Bot *tgbotapi.BotAPI
var ChatMap map[int64]*tgbotapi.Chat

type CommandHandle func(msg *tgbotapi.Message)

type Command struct {
	Name string
	Desc string
	Handle CommandHandle
}

var commands map[string]*Command

func InitBot(){

	token := "1541942531:AAGVCduRPaP2flKl64yt5kV3iXPNKSFMPfE"
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		panic(err)
	}
	bot.Debug = true
	log.Deubugf("Authorized on account %+v\n",bot.Self)
	Bot = bot

	//每个聊天频道不一样
	ChatMap = make(map[int64]*tgbotapi.Chat)

	//
	InitCommands()
}

func InitCommands() {
	commands = make(map[string]*Command)

	//drawinit - 初始化抽奖
	//info - 抽奖的整体情况
    //luckydraw - 抽奖
    //rollbackrecord 回滚抽奖

	regCommand("/help", "帮助", Help)
	regCommand("/drawinit", "初始化抽奖", DrawInit)
	regCommand("/info", "基本信息", Info)
	regCommand("/luckydraw", "抽奖", LuckyDraw)
	regCommand("/rollbackrecord", "回滚抽奖", RollbackRecord)
	regCommand("/makedrawuser", "选择抽奖用户", MakeDrawUser)
	regCommand("/adddrawuser", "新增抽奖名单", AddDrawUser)
	regCommand("/addawards", "新增奖品", AddAwards)
	regCommand("/sign", "签到", Sign)
	regCommand("/rollstart", "摇摇乐开始", RollStart)
	regCommand("/roll", "摇摇乐", Roll)
	regCommand("/rollinfo", "摇摇乐信息", RollInfo)

}

func regCommand(name, desc string, f CommandHandle) {
	c := &Command{
		Name:  name,
		Desc: desc,
		Handle: f,
	}
	commands[name] = c
}

func Run() {

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := Bot.GetUpdatesChan(u)
	if err != nil {
		log.Deubugf("GetUpdates err: %+v\n",err)
		return
	}

	for update := range updates {
		handle(&update)
	}
}

func handle(update *tgbotapi.Update) {
	//log.Deubugf("[%s] %s\n", update.Message.From.UserName, update.Message.Text)
	log.Deubugf("####################### UPDATE ########################\n")
	log.Deubugf("update:%+v\n", update)
	msg := update.Message
	if msg == nil  && update.CallbackQuery != nil && update.CallbackQuery.Message != nil {
		msg = update.CallbackQuery.Message
		msg.From = update.CallbackQuery.From  //from会被当做使用消息者
		log.Deubugf("update.CallbackQuery.Data:%+v\n", update.CallbackQuery.Data)
		msg.Text = update.CallbackQuery.Data
	}
	if  msg != nil  {
		log.Deubugf("update.Message:%+v\n", msg)
		if msg.Chat != nil {
			log.Deubugf("Chat :%+v\n", msg.Chat)
			ChatMap[msg.Chat.ID] = msg.Chat
			PrintChatList()
		}

		//分发命令
		if strings.HasPrefix(msg.Text, "/") {
			args := strings.Split(msg.Text, " ")

			name := strings.ToLower(args[0])
			if c, ok := commands[name]; ok {
				c.Handle(msg)
			} else {
				log.Deubugf("not found command: %s!\n", name)
			}
		} else {
			if msg.NewChatMembers != nil && len(*msg.NewChatMembers) > 0 {
				for _, u := range *msg.NewChatMembers {
					info := fmt.Sprintf("欢迎 %s 大驾光临",displayName(&u))
					m := tgbotapi.NewMessage(msg.Chat.ID, info)
					Bot.Send(m)

					mmark := tgbotapi.NewMessage(msg.Chat.ID, fmt.Sprintf("@%s 请签到准备抽奖",u.String()))
					button := tgbotapi.NewKeyboardButton("/sign 签到")
					row := tgbotapi.NewKeyboardButtonRow(button)
					markup := tgbotapi.NewReplyKeyboard(row)
					markup.OneTimeKeyboard = true
					markup.Selective = true
					mmark.ReplyMarkup = markup

					Bot.Send(mmark)
				}
			}
			log.Deubugf("%s : %s\n", msg.From.UserName, msg.Text)
		}
	}
}

func PrintChatList() {
	log.Deubugf("################ Chat List ##########\n")
	for _, c := range ChatMap {
		log.Deubugf("Chat: ID: %d Title: %s Type: %s\n",c.ID, c.Title, c.Type)
	}
}
