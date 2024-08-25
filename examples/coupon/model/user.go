package model

type User struct {
	Id            int64
	Account       string
	Name          string
	Phone         string
	Password      string
	State         int //1-正常、2-锁定、3-已注销
	CreateTime    int64
	UpdateTime    int64
	LastLoginTime int64
	RealNameState int //0-未实名、1-已实名
	IdNumber      string
}

func (t *User) TableName() string {
	return "user"
}
