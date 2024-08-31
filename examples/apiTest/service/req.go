package service

type UserLoginReq struct {
	Phone string `json:"phone"`
	Code  string `json:"code"`
}
