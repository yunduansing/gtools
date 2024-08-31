package model

import (
	"github.com/yunduansing/gtools/context"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"
)

type UserStatus string

const (
	UserStatusNormal  UserStatus = "normal"
	UserStatusDeleted UserStatus = "deleted"
	UserStatusBanned  UserStatus = "banned"
)

type User struct {
	UserId   int64  `json:"userId"`
	UserName string `json:"userName"`
	Phone    string `json:"phone"`
	//State    UserStatus `json:"state" gorm:"type:enum('active','inactive','banned')"`
}

func (*User) TableName() string {
	return "t_app_user"
}

type UserLoginLog struct {
	Id           int64  `json:"id"`
	UserId       int64  `json:"userId"`
	Phone        string `json:"phone"`
	LoginTime    int64  `json:"loginTime"`
	ClientType   string `json:"clientType"`
	LoginTimeout int64  `json:"loginTimeout"`
	RequestId    string `json:"requestId"`
	Msg          string `json:"msg"`
}

func (u *UserLoginLog) TableName() string {
	return "t_app_user_login_log"
}

func FindUserById(ctx *context.Context, id int64) (*User, error) {
	var user User
	err := dbContext.Read.First(ctx.Ctx, &user, func(db *gorm.DB, span trace.Span) *gorm.DB {
		db = db.Where("user_id=?", id)
		return db
	}).Error
	return &user, err
}

func FindUserByPhone(ctx *context.Context, phone string) (*User, error) {
	var user User
	err := dbContext.Read.First(ctx.Ctx, &user, func(db *gorm.DB, span trace.Span) *gorm.DB {
		db = db.Where("phone=?", phone)
		return db
	}).Error
	return &user, err
}
