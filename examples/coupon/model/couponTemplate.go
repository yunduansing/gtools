package model

type CouponTemplate struct {
	Id           int64
	TemplateCode string
	Type         string
	Amount       int64
	Threshold    int64
	ExpireType   int
	StartTime    int64
	EndTime      int64
	Days         int
	Stock        int
}

func (t *CouponTemplate) TableName() string {
	return "coupon_template"
}
