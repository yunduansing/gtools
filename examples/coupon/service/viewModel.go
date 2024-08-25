package service

import "github.com/yunduansing/gtools/examples/coupon/model"

type UserActivityCouponView struct {
	Activity    *model.Activity
	Coupon      *model.ActivityCoupon
	Stock       int
	UserImports string
}

type UserGetActivityCouponSuccess struct {
	ActivityId   int64
	UserId       int64
	TemplateCode string
	SuccessNum   int
	Time         int64
}
