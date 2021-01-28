package cmd

import (
	"luckydraw/draw"
	"luckydraw/log"
	"luckydraw/model"
	"luckydraw/tgbot"

)

var awards = []*model.Awards{
	&model.Awards{		Name:"Mackbook Pro(16)",  Num: 1},
	&model.Awards{		Name:"DJI Mavic 2 pro",   Num: 2},
	&model.Awards{		Name:"iPhone 12 Pro Max",  Num: 3},
	&model.Awards{		Name:"Ipad Pro",  Num: 4},
	&model.Awards{		Name:"Apple Watch Series 6",  Num: 5},
	&model.Awards{		Name:"AirPods Max",  Num: 6},
	&model.Awards{		Name:"AirPods Pro",  Num: 7},
	&model.Awards{		Name:"USDT 200",  Num: 4},
}

var names = []string {
	"Alexis",
	"ali",
	"angelo",
	"Clark Zeng",
	"Edward",
	"ember",
	"evin",
	"Gale",
	"hansen",
	"hunter",
	"James",
	"John(Jonz)",
	"Johnson",
	"Justin",
	"Kei.Kuan",
	"Landon",
	"Lirika",
	"Luke",
	"luotuo",
	"Mage",
	"Marco",
	"Seven(Unity)",
	"shafa",
	"Sky",
	"Sukie",
	"Sven",
	"Tanya",
	"Wayne",
	"Will",
	"william",
	"Zen",
	"miantiao",
}

func LetsDraw() {
	udraw := draw.NewUserDraw()

	//增加用户
	for _, name := range names {
		udraw.AddDrawUserByName(name)
	}
	udraw.ResetLeftUser()

	//增加奖品
	for _, item := range awards {
		udraw.AddTotalAwards(item)
	}
	udraw.ResetLeftAdwards()

	users := udraw.GetLeftUser()
	tgbot.PrintUser(users)


	awardsList := udraw.GetLeftAdwards()
	log.Deubugf("adwards: %d\n", len(awardsList))
	tgbot.PrintAwards(awardsList)

	for udraw.MakeNextDrawUser(draw.RANDOM_UID) != nil {
		udraw.DrawAdwards()
	}

	records := udraw.GetRecordList()
	tgbot.PrintRecords(records)

}

