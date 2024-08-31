package model

import (
	"github.com/yunduansing/gtools/context"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"
)

type User struct {
	UserId   int64  `json:"userId"`
	UserName string `json:"userName"`
}

func (*User) TableName() string {
	return "t_app_user"
}

func FindUserById(ctx *context.Context, id int64) (*User, error) {
	var user User
	err := dbContext.Read.First(ctx.Ctx, &user, func(db *gorm.DB, span trace.Span) *gorm.DB {
		db = db.Where("user_id=?", id)
		return db
	}).Error
	return &user, err
}
