package tgbot

import (
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"luckydraw/draw"
	"luckydraw/log"
	"luckydraw/model"
	"math/rand"
	"sort"
	"strconv"
	"strings"
	"time"
)

var (
	draws map[int64]*ChatDraw
	defaultAwards  []*model.Awards
	defaultUsers   []string
	signUser map[int]*tgbotapi.User
	//是否是管理员
	admins map[string]bool
)



func init() {
	draws = make(map[int64]*ChatDraw)

	defaultAwards = []*model.Awards{
		&model.Awards{		Name:"Mackbook Pro(16)",  Num: 1},
		&model.Awards{		Name:"DJI Mavic 2 pro",   Num: 2},
		&model.Awards{		Name:"iPhone 12 Pro Max",  Num: 3 },
		&model.Awards{		Name:"Ipad Pro",  Num: 4},
		&model.Awards{		Name:"Apple Watch Series 6",  Num: 5},
		&model.Awards{		Name:"AirPods Max",  Num: 6},
		&model.Awards{		Name:"AirPods Pro",  Num: 7},
		&model.Awards{		Name:"USDT 200",  Num: 4},
	}
	signUser = make(map[int]*tgbotapi.User)

	admins = make(map[string]bool)
	admins["wehi_shafa"] =  true
	admins["branway"] = true
}

type ChatDraw struct {
	ChatID int64
	udraw *draw.UserDraw
}

func NewChatDraw(chatid int64, udraw *draw.UserDraw) *ChatDraw{
	return &ChatDraw{
		ChatID: chatid,
		udraw: udraw,
	}
}

//command DrawInit
func DrawInit(m *tgbotapi.Message) {
	if !IsAdmin(m.From) {
		err_msg := tgbotapi.NewMessage(m.Chat.ID, fmt.Sprintf("%s 年轻人需要淡定\n", displayName(m.From)))
		Bot.Send(err_msg)
		return
	}

	chat := m.Chat
	msg := tgbotapi.NewMessage(chat.ID, "抽奖程序已准备就绪\n")
	if _, ok := draws[chat.ID]; ok {
		msg.Text = "draw aready exist!\n"
	} else {
		udraw := draw.NewUserDraw()

		//增加抽奖用户
		for _, u := range signUser {

			udraw.AddDrawUser(u.ID, displayName(u),u.String())
		}

		//增加奖品
		for _, item := range defaultAwards {
			udraw.AddTotalAwards(item)
		}

		resetDraw(udraw)

		//确定抽奖人
		//udraw.MakeNextDrawUser()

		cdraw := NewChatDraw(chat.ID, udraw)
		draws[chat.ID] = cdraw
	}
	cd := draws[chat.ID]
	msg.Text += infoDraw(cd.udraw)
	msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(false)
	Bot.Send(msg)


}

func Info(m *tgbotapi.Message) {
	chat := m.Chat
	if cdraw , ok := draws[chat.ID]; ok {
		udraw := cdraw.udraw
		info :=  infoDraw(udraw)
		reply := tgbotapi.NewMessage(chat.ID, info)
		Bot.Send( reply )
	} else {
		msg := tgbotapi.NewMessage(chat.ID, "抽奖程序还没准备好!\n")
		Bot.Send(msg)
	}
}

//幸运大抽奖
func LuckyDraw(m *tgbotapi.Message) {
	chat := m.Chat
	if cdraw , ok := draws[chat.ID]; ok {
		udraw := cdraw.udraw
		drawUser := udraw.GetDrawUser()
		if drawUser == nil {
			msg := tgbotapi.NewMessage(chat.ID, fmt.Sprintf("需要集齐7颗龙珠，方能召唤 抽奖 达人 !"))
			Bot.Send(msg)
			return
		}

		if drawUser.ID != m.From.ID && !IsAdmin(m.From){
			msg := tgbotapi.NewMessage(chat.ID, fmt.Sprintf("有请 [%s]抽奖，  [%s]请您稍后再抽",drawUser.Name, displayName(m.From)))
			Bot.Send(msg)
			return
		}

		awards := udraw.GetLeftAdwards()
		if len(awards) == 0 {
			msg := tgbotapi.NewMessage(chat.ID, "奖品都抽完了，亲")
			Bot.Send(msg)
			return
		}

		mwho := tgbotapi.NewMessage(chat.ID, fmt.Sprintf("现在抽奖的是: %s\n", drawUser.Name))
		mwho.ReplyMarkup = tgbotapi.NewRemoveKeyboard(false)
		Bot.Send(mwho)
		time.Sleep(time.Duration(rand.Int63n(3000))*time.Millisecond)
		m1 := tgbotapi.NewMessage(chat.ID, fmt.Sprintf(" %s 正在召唤 幸运之神...\n", drawUser.Name))
		Bot.Send(m1)

		time.Sleep(time.Duration(rand.Int63n(3000))*time.Millisecond)
		m2 := tgbotapi.NewMessage(chat.ID, fmt.Sprintf("%s 碎碎念.. 吾以吾神之名 --> 开: \n", drawUser.Name))
		Bot.Send(m2)

		r := udraw.DrawAdwards()
		mr := tgbotapi.NewMessage(chat.ID, fmt.Sprintf("恭喜 %s  喜得神器: %s\n", r.User.Name, r.Awards.Name))
		Bot.Send(mr)

		time.Sleep(3000*time.Millisecond)
		info := tgbotapi.NewMessage(chat.ID, infoDraw(udraw))
		Bot.Send(info)


	} else {
		msg := tgbotapi.NewMessage(chat.ID, "抽奖程序还没准备好!\n")
		Bot.Send(msg)
	}
}

//回退一条抽奖记录
func RollbackRecord(m *tgbotapi.Message) {
	if !IsAdmin(m.From) {
		err_msg := tgbotapi.NewMessage(m.Chat.ID, fmt.Sprintf("%s 年轻人需要淡定\n", displayName(m.From)))
		Bot.Send(err_msg)
		return
	}

	chat := m.Chat
	if cdraw , ok := draws[chat.ID]; ok {
		udraw := cdraw.udraw

		info := strings.Split(m.Text, " ")
		log.Deubugf("CancleRecord %+v", info)
		reply := tgbotapi.NewMessage(chat.ID, "回退失败")
		if len(info) > 1 {

			if rid, err := strconv.Atoi(info[1]); err == nil {
				rd := udraw.RollbackRecord(rid)
				if rd != nil {
					reply.Text = fmt.Sprintf("回退成功: %s  撤回 奖品: %s",rd.User.Name,rd.Awards.Name )
				} else {
					reply.Text = fmt.Sprintf("没找到记录")
				}
			}
		} else {
			//选择id
			list := udraw.GetRecordList()
			if len(list) > 0 {
				markup := tgbotapi.NewReplyKeyboard()
				for _, r := range list {
					button := tgbotapi.NewKeyboardButton(fmt.Sprintf("/rollbackrecord %d [%s]-->[%s]",r.RID, r.User.Name, r.Awards.Name))
					row :=  tgbotapi.NewKeyboardButtonRow(button)
					markup.Keyboard = append(markup.Keyboard, row)
				}
				markup.OneTimeKeyboard = true //一次消失
				reply.ReplyMarkup = markup
				reply.Text = "请选择记录"
			} else {
				reply.Text = "还没有抽奖记录"
			}
		}
		Bot.Send(reply)
	} else {
		msg := tgbotapi.NewMessage(chat.ID, "抽奖程序还没准备好!\n")
		Bot.Send(msg)
	}
}

//确定下一个抽奖的用户
func MakeDrawUser(m *tgbotapi.Message) {
	if !IsAdmin(m.From) {
		err_msg := tgbotapi.NewMessage(m.Chat.ID, fmt.Sprintf("%s 年轻人需要淡定\n", displayName(m.From)))
		Bot.Send(err_msg)
		return
	}

	chat := m.Chat
	if cdraw , ok := draws[chat.ID]; ok {
		udraw := cdraw.udraw
		info := strings.Split(m.Text, " ")
		log.Deubugf("MakeDrawUser %+v", info)
		userinfo := "空"
		if drawuser := udraw.GetDrawUser(); drawuser != nil {
			userinfo = drawuser.Name
		}

		reply := tgbotapi.NewMessage(chat.ID, fmt.Sprintf("当前抽奖:[ %s ],请选择抽奖人", userinfo))

		if len(info) > 1 {


			if uid, err := strconv.Atoi(info[1]); err == nil {
				u := udraw.MakeNextDrawUser(uid)
				if u != nil {
					reply.Text = fmt.Sprintf("下个抽抽奖的是:  @%s",u.Account)

					//增加抽奖快捷按钮
					button := tgbotapi.NewKeyboardButton("/luckydraw 抽奖")
					row := tgbotapi.NewKeyboardButtonRow(button)
					markup := tgbotapi.NewReplyKeyboard(row)
					markup.OneTimeKeyboard = true
					markup.Selective = true
					reply.ReplyMarkup = markup
				} else {
					reply.Text = fmt.Sprintf("选择失败")
				}
			}

		} else {
			//选择id
			list := udraw.GetLeftUser()
			if len(list) > 0 {
				reply.Text += fmt.Sprintf(" @%s 剩余: %d", m.From.String(),len(list))
				rbutton := tgbotapi.NewKeyboardButton("/makedrawuser -1 随机")
				rrow :=  tgbotapi.NewKeyboardButtonRow(rbutton)
				markup := tgbotapi.NewReplyKeyboard(rrow)
				for _, u := range list {
					button := tgbotapi.NewKeyboardButton(fmt.Sprintf("/makedrawuser %d  %s",u.ID, u.Name))
					row :=  tgbotapi.NewKeyboardButtonRow(button)
					markup.Keyboard = append(markup.Keyboard, row)
				}
				markup.OneTimeKeyboard = true //一次隐藏
				markup.Selective = true

				reply.ReplyMarkup = markup

			} else {
				reply.Text = "都抽完啦 亲"
			}
		}
		Bot.Send(reply)
	} else {
		msg := tgbotapi.NewMessage(chat.ID, "抽奖程序还没准备好!\n")
		Bot.Send(msg)
	}
}

//额外添加抽奖的用户
func AddDrawUser(m *tgbotapi.Message) {
	if !IsAdmin(m.From) {
		err_msg := tgbotapi.NewMessage(m.Chat.ID, fmt.Sprintf("%s 年轻人需要淡定\n", displayName(m.From)))
		Bot.Send(err_msg)
		return
	}

	chat := m.Chat
	if cdraw , ok := draws[chat.ID]; ok {
		udraw := cdraw.udraw
		info := strings.Split(m.Text, " ")
		log.Deubugf("AddDrawUser %+v", info)

		if len(info) > 1 {

			udraw.AddDrawUserByName(info[1])
			msg := tgbotapi.NewMessage(chat.ID, fmt.Sprintf("添加新的 抽奖名单: %s", info[1]))
			Bot.Send(msg)

		} else {

			msg := tgbotapi.NewMessage(chat.ID, "添加抽奖名单参数错误\n")
			Bot.Send(msg)
		}

	} else {
		info := strings.Split(m.Text, " ")
		log.Deubugf("AddDrawUser %+v", info)

		msg := tgbotapi.NewMessage(chat.ID, "抽奖程序还没准备好!\n")
		if len(info) > 1 {
			defaultUsers = append(defaultUsers, info[1])
			msg.Text = fmt.Sprintf("新增默认抽奖名单: %s",info[1])
		}

		Bot.Send(msg)
	}
}

//签到
func Sign(m *tgbotapi.Message) {
	chat := m.Chat
	info := strings.Split(m.Text, " ")
	log.Deubugf("Sign %+v", info)

	msg := tgbotapi.NewMessage(chat.ID, "抽奖程序还没准备好!\n")
	name := displayName(m.From)
	msg.Text = fmt.Sprintf("[ %s ]已签到",name)
	signUser[m.From.ID] = m.From

	var names []string
	for _,u := range signUser {
		names = append(names,displayName(u))
	}
	msg.Text += "\n########################\n"
	msg.Text += formatSlice(names, 5)
	msg.Text += fmt.Sprintf("\n签到人数: %d",len(names))

	Bot.Send(msg)
}



//新增奖品
func AddAwards(m *tgbotapi.Message) {
	if !IsAdmin(m.From) {
		err_msg := tgbotapi.NewMessage(m.Chat.ID, fmt.Sprintf("%s 年轻人需要淡定\n", displayName(m.From)))
		Bot.Send(err_msg)
		return
	}

	chat := m.Chat
	if cdraw , ok := draws[chat.ID]; ok {
		udraw := cdraw.udraw
		info := strings.Split(m.Text, " ")
		log.Deubugf("AddAwards %+v\n", info)

		if len(info) > 1 {
			num := 1
			if len(info) >2 {
				if i, err := strconv.Atoi(info[2]); err == nil {
					num = i
				}else {
					log.Deubugf("######### atoi %s err: %+v\n",info[2], err)
				}
			}
			udraw.AddAwardsInfo(info[1],num, "")
			msg := tgbotapi.NewMessage(chat.ID, fmt.Sprintf("添加新的 奖品: %s 数量: %d", info[1], num))
			Bot.Send(msg)

		} else {

			msg := tgbotapi.NewMessage(chat.ID, "添加奖品参数错误\n")
			Bot.Send(msg)
		}

	} else {
		msg := tgbotapi.NewMessage(chat.ID, "抽奖程序还没准备好!\n")
		Bot.Send(msg)
	}
}

//
func RollStart(m *tgbotapi.Message) {
	if !IsAdmin(m.From) {
		err_msg := tgbotapi.NewMessage(m.Chat.ID, fmt.Sprintf("%s 年轻人需要淡定\n", displayName(m.From)))
		Bot.Send(err_msg)
		return
	}

	chat := m.Chat

	if cdraw , ok := draws[chat.ID]; ok {
		udraw := cdraw.udraw
		udraw.RollStart()
		msg := tgbotapi.NewMessage(chat.ID, "瑶瑶乐 以准备就绪 /roll \n")
		Bot.Send(msg)

	} else {
		msg := tgbotapi.NewMessage(chat.ID, "抽奖程序还没准备好!\n")
		Bot.Send(msg)
	}
}

//
func Roll(m *tgbotapi.Message) {
	chat := m.Chat

	if cdraw , ok := draws[chat.ID]; ok {
		udraw := cdraw.udraw
		name := displayName(m.From)

		luck, ok := udraw.HaveRoll(name)
		if luck == -1 {
			msg := tgbotapi.NewMessage(chat.ID, "瑶瑶乐还没准备好")
			Bot.Send(msg)
		}else if ok {
			msg := tgbotapi.NewMessage(chat.ID, fmt.Sprintf("[ %s ] 已经摇过了 以前的结果是:  %d  \n", name, luck))
			Bot.Send(msg)
		}else {
			luck := udraw.Roll(name)
			msg := tgbotapi.NewMessage(chat.ID, fmt.Sprintf("[ %s ] 摇出的幸运数字是 %d  \n", name, luck))
			Bot.Send(msg)
		}


	} else {
		msg := tgbotapi.NewMessage(chat.ID, "抽奖程序还没准备好!\n")
		Bot.Send(msg)
	}
}
//
func RollInfo(m *tgbotapi.Message) {
	chat := m.Chat

	if cdraw , ok := draws[chat.ID]; ok {
		udraw := cdraw.udraw

		info := udraw.RollInfo()
		msg := tgbotapi.NewMessage(chat.ID, fmt.Sprintf("目前的结果是: \n%s\n", info))
		Bot.Send(msg)

	} else {
		msg := tgbotapi.NewMessage(chat.ID, "抽奖程序还没准备好!\n")
		Bot.Send(msg)
	}
}


func resetDraw(udraw *draw.UserDraw) {
	log.Deubugf("resetDreaw\n")
	udraw.ResetLeftUser()
	udraw.ResetLeftAdwards()
	udraw.ResetRecord()
}

func infoDraw(udraw *draw.UserDraw) string {
	r := udraw.GetRecordList()
	a := udraw.GetLeftAdwards()
	u := udraw.GetLeftUser()

	info := fmt.Sprintf("####### 中奖信息 #######\n%s\n" +
		"########## 剩余奖品 ###########\n%s\n\n" +
		"########### 还未抽奖 ############\n%s\n\n",
		PrintRecords(r),
		PrintAwards(a),
		PrintUser(u))
	return info
}

func PrintUser(list []*model.User) string {

	if len(list) == 0 {
		return "空"
	}

	var plans  []string
	var index = 0
	var line []string
	for _, u := range list {
		str := u.Name//fmt.Sprintf("ID: %d Name: %s", u.ID, u.Name)
		line = append(line, str)
		index += 1
		if index % 5 == 0 {
			lstr := strings.Join(line, "   #   ")
			plans = append(plans,lstr)
			index = 0
			line = nil
		}
	}
	if line != nil {
		lstr := strings.Join(line, "   #   ")
		plans = append(plans,lstr)
	}

	totalinfo := fmt.Sprintf("total: 【 %d 】", len(list))
	plans = append(plans, totalinfo)
	info := strings.Join(plans, splitLine())
	log.Deubugf(info)
	log.Deubugf("\n")
	return info
}

func formatSlice(list []string, linenum int) string{
	var plans  []string
	var index = 0
	var line []string
	var joinstr = " , "
	for _, str := range list {

		line = append(line, fmt.Sprintf("【 %s 】",str))
		index += 1
		if index % linenum == 0 {
			lstr := strings.Join(line, joinstr)
			plans = append(plans,lstr)
			index = 0
			line = nil
		}
	}
	if line != nil {
		lstr := strings.Join(line, joinstr)
		plans = append(plans,lstr)
	}
	return strings.Join(plans,"\n")
}

func PrintAwards(list []*model.Awards) string {
	if len(list) == 0 {
		return "空"
	}

	sort.Slice(list, func(i, j int) bool {
		return list[i].Order < list[j].Order
	})

	var sortNames []string
	var awardsInfo = make(map[string]int)
	for _, a := range list {
		//log.Deubugf("%d %+v\n",index,a.Name)
		if _, ok := awardsInfo[a.Name]; ok {
			awardsInfo[a.Name] += a.Num
		}else {
			awardsInfo[a.Name] =  a.Num
			sortNames = append(sortNames, a.Name)
		}
	}
	log.Deubugf("\n#########################\nsortNames: %s\n",strings.Join(sortNames, ", "))

	var plans []string
	var total int
	for _, name := range sortNames {
		str := fmt.Sprintf("%s : 【 %d 】", name, awardsInfo[name])
		plans = append(plans, str)
		total += awardsInfo[name]
	}
	sum := fmt.Sprintf("total: 【 %d 】", total)
	plans = append(plans, sum)
	info := strings.Join(plans, splitLine())
	log.Deubugf(info)
	log.Deubugf("\n")
	return info
}

func PrintRecords(list []*model.AwardsRecord) string {
	if len(list) == 0 {
		return "空"
	}
	sort.Slice(list, func(i, j int) bool {
		return list[i].Order < list[j].Order
	})
	var sortNames []string //排序的名字
	var rinfo = make(map[string][]string)
	for _, r := range list {
		//log.Deubugf("%+v\n", r.Awards)
		item := rinfo[r.Awards.Name]
		item = append(item, r.User.Name)
		rinfo[r.Awards.Name] = item
		if len(item) == 1 {
			sortNames = append(sortNames,r.Awards.Name)
		}
	}

	var plans []string
	for _, aname := range sortNames {
		owners := rinfo[aname]
		str := fmt.Sprintf("%s -->: [ %s ]",aname, strings.Join(owners, ", "))
		plans = append(plans, str)
	}
	total := fmt.Sprintf("total: 【 %d 】", len(list))
	plans = append(plans, total)
	info := strings.Join(plans, splitLine())
	log.Deubugf(info)
	log.Deubugf("\n")
	return info
}

func splitLine() string {
	//return fmt.Sprintf("\n%s\n",strings.Repeat("-",40))
	return "\n"
}

func displayName(u *tgbotapi.User) string {
	return u.FirstName+u.LastName
}
