package service

type UserLoginRes struct {
	UserId   int64  `json:"userId"`
	Phone    string `json:"phone"`
	UserName string `json:"userName"`
	Token    string `json:"token"`
}
