package model

type Awards struct {
	Name string  //奖品名称
	Img string // 图片地址
	Num int //数量
}

func (item *Awards) Copy() *Awards{
	newItem := &Awards{
		Name: item.Name,
		Img: item.Img,
		Num: item.Num,
	}
	return newItem
}
