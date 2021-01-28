package model

type Awards struct {
	Name string  //奖品名称
	Img string // 图片地址
	Num int //数量
	Order int //排序值，越大越后
}

func (item *Awards) Copy() *Awards{
	newItem := &Awards{
		Name: item.Name,
		Img: item.Img,
		Num: item.Num,
		Order: item.Order,
	}
	return newItem
}
