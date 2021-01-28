package draw

import (
	"errors"
	"fmt"
	"luckydraw/log"
	"luckydraw/model"
	"luckydraw/random"
	"math/rand"
	"strings"
	"sync"
)

const RANDOM_UID = -1

//以用户为中心的抽奖
// 1. 确定获奖用户
// 2. 随机奖品

type UserDraw struct {
	totalAwards    map[string]*model.Awards  //总的奖品数
	totalLock sync.RWMutex

	leftAwards []*model.Awards //未抽的奖品池

	toatalUser map[int]*model.User //所有用户
	userLock   sync.RWMutex

	leftUser []int //剩余没抽奖的用户id

	recordList []*model.AwardsRecord //已经中奖的记录
	rid int //记录的id

	drawUser *model.User //下个抽奖的用户id
	uid int //抽奖用户的id

	rollinfo map[string]int //用户id /
}

func NewUserDraw() *UserDraw {
	udraw := &UserDraw{
		totalAwards: make(map[string]*model.Awards),
		toatalUser:  make(map[int]*model.User),
	}
	return udraw
}

//添加基本信息添加一个奖品
func (udraw *UserDraw)AddAwardsInfo(name string, num int, img string) {

	item := &model.Awards{
		Name: name,
		Num: num,
		Img: img,
	}
	udraw.AddTotalAwards(item)
	newItem := udraw.totalAwards[name]
	item.Order = newItem.Order
	if num > 0 {
		for i := 0; i < num; i++ {
			a := item.Copy()
			a.Num = 1
			udraw.AddLeftAwards(a)
		}
	}
}

//增加一个奖品清单
func (udraw *UserDraw)AddTotalAwards(item *model.Awards) error {
	udraw.totalLock.Lock()
	defer udraw.totalLock.Unlock()
	if udraw.totalAwards == nil {
		return errors.New("not found total awards")
	}
	old, ok := udraw.totalAwards[item.Name]
	if ok {
		old.Num += item.Num
	} else {
		udraw.totalAwards[item.Name] = item.Copy()
		udraw.totalAwards[item.Name].Order = len(udraw.totalAwards)
	}
	return nil
}

//减少清单的数量
func (udraw *UserDraw)DelTotalAwardsByNum(name string, num int) error {
	udraw.totalLock.Lock()
	defer udraw.totalLock.Unlock()
	if udraw.totalAwards == nil {
		return errors.New("not found total awards")
	}
	old, ok := udraw.totalAwards[name]
	if ok {
		old.Num -= num
		if old.Num <= 0 {
			delete(udraw.totalAwards, name)
		}
	} else {
		return errors.New("not found awards: " + name)
	}
	return nil
}

//获取整个奖品的清单
func (udraw *UserDraw)GetTotalAwardsList()[]*model.Awards {
	udraw.totalLock.RLock()
	udraw.totalLock.RUnlock()
	list := make([]*model.Awards, 0, len(udraw.totalAwards))
	for _, item := range udraw.totalAwards {
		list = append(list, item)
	}
	return list
}


//初始奖品池
func (udraw *UserDraw) ResetLeftAdwards() []*model.Awards{
	var pool []*model.Awards

	udraw.totalLock.Lock()
	for _, item := range udraw.totalAwards {
		if item.Num > 0 {
			for i := 0; i < item.Num; i++ {
				a := item.Copy()
				a.Num = 1
				pool = append(pool, a)
			}
		}
	}
	udraw.totalLock.Unlock()

	//log.d
	udraw.leftAwards = pool

	return udraw.GetLeftAdwards()
}

//获取剩余的奖品
func (udraw *UserDraw)GetLeftAdwards() []*model.Awards {
	return udraw.leftAwards
}

//增加用户
func (udraw *UserDraw) AddTotalUser(u *model.User) error {
	udraw.userLock.Lock()
	udraw.userLock.Unlock()
	if _, ok := udraw.toatalUser[u.ID]; ok {
		return errors.New("already has user: " + u.Name)
	}
	udraw.toatalUser[u.ID] = u.Copy()
	return nil
}

//添加一个可抽奖人的名字
func (udraw *UserDraw)AddDrawUserByName(name string) {
	u := &model.User{
		ID: udraw.nextuid(),
		Name: name,
		Account: name,
	}
	udraw.AddTotalUser(u)
	udraw.AddLeftUser(u.ID)
}

//添加一个可抽奖人的名字
func (udraw *UserDraw)AddDrawUser(id int, name, account string) {
	u := &model.User{
		ID: id,
		Name: name,
		Account: account,
	}
	udraw.AddTotalUser(u)
	udraw.AddLeftUser(u.ID)
}

func (udraw *UserDraw)GetTotalUserList() []*model.User {
	var list []*model.User
	udraw.userLock.RLock()
	for _, u := range udraw.toatalUser {
		list = append(list, u)
	}
	udraw.userLock.RUnlock()
	return list
}

func (udraw *UserDraw) ResetLeftUser() []*model.User{
	var list []*model.User
	var ids []int
	udraw.userLock.RLock()
	for _, u := range udraw.toatalUser {
		list = append(list, u)
		ids = append(ids, u.ID)
	}
	udraw.leftUser = ids
	udraw.userLock.RUnlock()
	return list
}

func (udraw *UserDraw) GetLeftUser() []*model.User {
	var list []*model.User
	udraw.userLock.RLock()
	for _, id := range udraw.leftUser {
		u := udraw.toatalUser[id]
		list = append(list, u)
	}
	udraw.userLock.RUnlock()
	return list
}

//随机确定下一个抽奖用户
func (udraw *UserDraw) RandomNextDrawUser() *model.User {
	return udraw.MakeNextDrawUser(RANDOM_UID)
}

//下一个抽奖人
func (udraw *UserDraw) MakeNextDrawUser(uid int) *model.User {
	if count := len(udraw.leftUser); count > 0 {
		index, _ := random.Uint64Range(0, uint64(count))
		id := udraw.leftUser[index]

		if uid != RANDOM_UID {
			var found bool = false
			for _, i := range udraw.leftUser {
				if uid == i {
					id = uid
					found = true
					break
				}
			}
			if !found {
				log.Deubugf("not found left user: \n", uid)
				return nil
			}
		}



		udraw.drawUser = udraw.toatalUser[id].Copy()

		log.Deubugf("new draw user: %+v\n", udraw.drawUser)
		return udraw.drawUser.Copy()
	}
	log.Deubugf("no left user\n")
	udraw.drawUser = nil
	return nil
}

//获取正在抽奖的人
func (udraw *UserDraw) GetDrawUser() *model.User{
	return udraw.drawUser
}
//抽奖
func (udraw *UserDraw)DrawAdwards() *model.AwardsRecord {
	if udraw.drawUser == nil {
		log.Deubugf("not draw user \n")
		return nil
	}
	if count := len(udraw.leftAwards); count > 0 {
		index,_:= random.Uint64Range(0, uint64(count))

		awards := udraw.leftAwards[index]
		record := &model.AwardsRecord{
			RID : udraw.nextRID(),
			User:udraw.drawUser,
			Awards:awards,
		}
		udraw.recordList = append(udraw.recordList, record)
		log.Deubugf("DrawAdwards record: [%d] %s -> %s\n",record.RID, record.User.Name, record.Awards.Name)

		//删除中奖的用户
		var leftUser []int
		for _, id := range udraw.leftUser {
			if id != udraw.drawUser.ID {
				leftUser = append(leftUser, id)
			}
		}
		udraw.leftUser = leftUser

		//删除中奖的的奖品
		var leftAwards []*model.Awards
		for i, a := range udraw.leftAwards {
			if i != int(index) {
				leftAwards = append(leftAwards, a)
			}
		}
		udraw.leftAwards = leftAwards

		//自动变换下一个抽奖用户
		//udraw.MakeNextDrawUser()
		udraw.drawUser = nil

		return record
	} else {
		return nil
	}
}

func (udraw *UserDraw)ResetRecord() {
	udraw.recordList = nil
	udraw.rid = 0
}
func (udraw *UserDraw)nextRID() int {
	udraw.rid += 1
	return udraw.rid
}

func (udraw *UserDraw) GetRecordList() []*model.AwardsRecord {
	return udraw.recordList
}

//回滚一条中奖信息
func (udraw *UserDraw) RollbackRecord(id int) (dr *model.AwardsRecord) {
	if len(udraw.recordList) > 0 {
		var list []*model.AwardsRecord
		for _, r := range udraw.recordList {
			if r.RID != id {
				list = append(list, r)
			} else {
				//删除成功，把奖品添加到剩余的当中去
				udraw.AddLeftUser(r.User.ID)
				udraw.AddLeftAwards(r.Awards.Copy())
				dr = r
				log.Deubugf("delete record %+v\n",r)
			}
		}
		udraw.recordList = list
	}
	return
}

func (udraw *UserDraw) AddLeftUser(id int) {
	udraw.leftUser = append(udraw.leftUser, id)
}

func (udraw *UserDraw) AddLeftAwards(a *model.Awards) {
	udraw.leftAwards = append(udraw.leftAwards, a)
}

func (udraw *UserDraw)RollStart() {
	udraw.rollinfo = make(map[string]int)
}

func (udraw *UserDraw)Roll(name string) int {
	if udraw.rollinfo == nil {
		return -1
	}
	luck := rand.Intn(1000)
	udraw.rollinfo[name] = luck
	return luck
}

func (udraw *UserDraw)HaveRoll(name string)(int, bool) {
	if udraw.rollinfo == nil {
		return -1, false
	}
	luck, ok := udraw.rollinfo[name]
	return luck, ok
}

func (udraw *UserDraw)RollInfo() string {
	if udraw.rollinfo == nil {
		return ""
	}
	var maxinfo string = ""
	var max = 0
	var lines []string
	for name, luck := range udraw.rollinfo {
		s := fmt.Sprintf("[%s] -> %d",name, luck)
		if luck > max {
			max = luck
			maxinfo = s
		}
		lines = append(lines, s)
	}
	str := strings.Join(lines, "\n")
	str += fmt.Sprintf("\n最高的是: %s" , maxinfo)
	return str
}

func (udraw *UserDraw) nextuid() int{
	udraw.uid += 1
	return udraw.uid
}


