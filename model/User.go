package model

type User struct {
	ID int     //用户的id
	Name string   //用户的名字
	Account string //账号名
}

func (u *User)Copy() *User{
	return &User{
		ID: u.ID,
		Name: u.Name,
		Account: u.Account,
	}
}
