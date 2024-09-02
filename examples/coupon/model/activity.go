package model

import (
	mysqltool "github.com/yunduansing/gtools/database"
)

const (
	ActivityUserRangeAll    = "all"
	ActivityUserRangeImport = "import"
	ActivityUserRangeNew    = "new"
)

type Activity struct {
	Id         int64
	Name       string
	State      int //1启用、2-禁用
	AuditState int //1-通过、2-驳回
	AuditTime  int64
	StartTime  int64
	EndTime    int64
	CreateTime int64
	UpdateTime int64
	UserRange  string //all-全部、import-导入、new-新注册
}

func (t *Activity) TableName() string {
	return "activity"
}

type ActivityCoupon struct {
	ActivityId    int64
	TemplateCode  string
	Num           int //投放数量
	OneAccountNum int //单账号领取数量
}

func (t *ActivityCoupon) TableName() string {
	return "activity_coupon"
}

func FindActivityById(activityId int64) (res *Activity, err error) {
	err = mysqltool.Get(ReadMysql).DB.Where("id = ?", activityId).First(res).Error
	return
}

func FindActivityCouponList(activityId int64) (res []*ActivityCoupon, err error) {
	err = mysqltool.Get(ReadMysql).DB.Where("activity_id = ?", activityId).Find(res).Error
	return
}

type ActivityCouponGroupView struct {
	ActivityId   int64
	TemplateCode string
	Num          int
}

func FindActivityCouponGroup(activityId int64, templateCode string) (res *ActivityCouponGroupView, err error) {
	err = mysqltool.Get(ReadMysql).DB.Model(&ActivityCoupon{}).Where("activity_id = ? and template_code = ?",
		activityId, templateCode).Select("activity_id,template_code,count(1) as num").Find(&res).Error
	return

}
