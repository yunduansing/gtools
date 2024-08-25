package model

type ErrorCode struct {
	Code int
	Msg  string
}

func (e ErrorCode) Error() string {
	return e.Msg
}

func NewErrorCodeErr(code int, err error) ErrorCode {
	return ErrorCode{
		Code: code,
		Msg:  err.Error(),
	}
}

func NewErrorCode(code int, msg string) ErrorCode {
	return ErrorCode{
		Code: code,
		Msg:  msg,
	}
}

var (
	Success                  = ErrorCode{Code: 0, Msg: "success"}
	ErrActivityNotExists     = ErrorCode{Code: 10001, Msg: "活动不存在"}
	ErrCouponNotInActivity   = ErrorCode{Code: 10002, Msg: "优惠券不在活动范围内"}
	ErrUserNotInActivity     = ErrorCode{Code: 10003, Msg: "用户不在活动范围内"}
	ErrUserAlreadyGetCoupon  = ErrorCode{Code: 10004, Msg: "用户已领取优惠券"}
	ErrActivityAlreadyEnd    = ErrorCode{Code: 10005, Msg: "活动已结束"}
	ErrActivityNotStart      = ErrorCode{Code: 10006, Msg: "活动未开始"}
	ErrActivityStateOffline  = ErrorCode{Code: 10007, Msg: "活动已下线"}
	ErrActivityAuditReject   = ErrorCode{Code: 10008, Msg: "活动审核未通过"}
	ErrTemplateCodeInvalid   = ErrorCode{Code: 10009, Msg: "优惠券模板编码无效"}
	ErrActivityParamsError   = ErrorCode{Code: 10010, Msg: "活动参数错误"}
	ErrActivityCouponTooHot  = ErrorCode{Code: 10011, Msg: "活动优惠券太火爆"}
	ErrActivityCouponNoStock = ErrorCode{Code: 10012, Msg: "活动优惠券已抢完"}
)
