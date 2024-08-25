package model

import mysqltool "github.com/yunduansing/gtools/database/mysql"

type UseCoupon struct {
	UserId       int64
	ActivityId   int64
	TemplateCode string
	CouponCode   string
	State        int //1-正常、2-已锁定、3-已核销、4-已作废
	CreateTime   int64
	StartTime    int64
	EndTime      int64
	Source       int //1-自主领取、2-系统发放、3-运营发放
}

func (t *UseCoupon) TableName() string {
	return "use_coupon"
}

type ActivityUserImport struct {
	ActivityId int64
	UserId     int64
}

func (t *ActivityUserImport) TableName() string {
	return "activity_user_import"
}

type UserActivityCouponGroupView struct {
	ActivityId   int64
	UserId       int64
	TemplateCode string
	Num          int
}

func FindUserActivityCouponGroup(activityId, userId int64,
	templateCode string) (res *UserActivityCouponGroupView,
	err error) {
	err = mysqltool.Get(ReadMysql).Model(&UseCoupon{}).Where("activity_id = ? and user_id and template_code = ?",
		activityId, userId, templateCode).Group("user_id," +
		"template_code").Select("user_id,count(*) as num").Find(&res).Error
	return

}

func FindActivityUserImportList(activityId int64) (res []*ActivityUserImport, err error) {
	err = mysqltool.Get(ReadMysql).Where("activity_id = ?", activityId).Find(res).Error
	return
}
