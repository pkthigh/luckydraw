package cmd

import (
	"fmt"
	"luckydraw/draw"
	"luckydraw/log"
	"luckydraw/model"
	"strings"
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
	PrintUser(users)

	log.Deubugf("adwards: \n")
	awardsList := udraw.GetLeftAdwards()
	PrintAwards(awardsList)

	for udraw.MakeNextDrawUser() != nil {
		udraw.DrawAdwards()
	}

	records := udraw.GetRecordList()
	PrintRecords(records)

}

func PrintUser(list []*model.User) string {
	var plans  []string
	for _, u := range list {
		str := fmt.Sprintf("ID: %d Name: %s", u.ID, u.Name)
		plans = append(plans, str)
	}
	totalinfo := fmt.Sprintf("total: %d", len(list))
	plans = append(plans, totalinfo)
	info := strings.Join(plans, "\n")
	log.Deubugf(info)
	log.Deubugf("\n")
	return info
}

func PrintAwards(list []*model.Awards) string {
	var awardsInfo = make(map[string]int)
	for _, a := range list {
		if num, ok := awardsInfo[a.Name]; ok {
			awardsInfo[a.Name] = num + a.Num
		}else {
			awardsInfo[a.Name] =  a.Num
		}
	}

	var plans []string
	var total int
	for name, num := range awardsInfo {
		str := fmt.Sprintf("%s : %d", name, num)
		plans = append(plans, str)
		total += num
	}
	sum := fmt.Sprintf("total: %d", total)
	plans = append(plans, sum)
	info := strings.Join(plans, "\n")
	log.Deubugf(info)
	log.Deubugf("\n")
	return info
}

func PrintRecords(list []*model.AwardsRecord) string {
	var rinfo = make(map[string][]string)
	for _, r := range list {
		item := rinfo[r.Awards.Name]
		item = append(item, r.User.Name)
		rinfo[r.Awards.Name] = item
	}

	var plans []string
	for aname, owners := range rinfo {
		str := fmt.Sprintf("%s: [ %s ]",aname, strings.Join(owners, ", "))
		plans = append(plans, str)
	}
	total := fmt.Sprintf("total: %d", len(list))
	plans = append(plans, total)
	info := strings.Join(plans, "\n")
	log.Deubugf(info)
	log.Deubugf("\n")
	return info
}
