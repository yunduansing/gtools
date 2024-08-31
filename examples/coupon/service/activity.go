package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/yunduansing/gtools/context"
	"github.com/yunduansing/gtools/examples/apiTest/config"
	"github.com/yunduansing/gtools/examples/coupon/model"
	"github.com/yunduansing/gtools/redistool"
	"github.com/yunduansing/gtools/utils"
	"gorm.io/gorm"
	"strconv"
	"strings"
	"time"
)

type UserGetCouponReq struct {
	UserId       int64
	ActivityId   int64
	TemplateCode string
}

func UserGetCoupon(ctx context.Context, req *UserGetCouponReq) (err error) {
	if len(req.TemplateCode) == 0 || req.ActivityId == 0 {
		return model.ErrActivityParamsError
	}

	ctx.Log.Infof(ctx.Ctx, "领取优惠券，开始执行   req======%+v", req)

	var activity UserActivityCouponView
	if err = loadActivity(ctx, req, &activity); err != nil {
		return
	}
	if err = checkActivityCoupon(ctx, req, &activity); err != nil {
		return
	}

	successNum, err := userGetCouponFromRedis(ctx, req, &activity)
	if err != nil {
		return
	}
	notifyGetCouponSuccess(ctx, req, successNum)
	return
}

var (
	activityCouponKeyFormat = "{activity:%d:coupon:%s}:config" +
		"" //存储活动、优惠券库存和优惠券配置数据
	userActivityCouponKeyFormat = "{activity:%d:coupon:%s}:user:%d" //存储用户领取记录
	activityFromDbKeyFormat     = "distlock:activity:%d:coupon:%s"

	activityFields = []string{"activity", "coupon", "stock", "userImports"}
)

func loadActivity(ctx context.Context, req *UserGetCouponReq, activity *UserActivityCouponView) (err error) {
	ctx.Log.Infof(ctx.Ctx, "领取优惠券，step-loadActivity:加载活动数据   req======%+v", req)
	//第一步，先尝试直接从redis加载
	activityCouponKey := fmt.Sprintf(activityCouponKeyFormat, req.ActivityId, req.TemplateCode)
	//userActivityCouponKey:=fmt.Sprintf(userActivityCouponKeyFormat, req.ActivityId, req.TemplateCode, req.UserId)
	err = loadActivityFromRedis(ctx, req, activity)
	if err != nil || activity.Activity != nil {
		return
	}

	distLockKey := fmt.Sprintf(activityFromDbKeyFormat, req.ActivityId, req.TemplateCode)
	distLock := redistool.NewRedisLock(config.Redis.UniversalClient, distLockKey)
	distLockSuccess, err := distLock.AcquireBackoff(1000, 300*time.Microsecond, 10*time.Millisecond)
	if err != nil {
		ctx.Log.Errorf(ctx.Ctx, "领取优惠券，分布式锁获取失败：req=%+v,err=%+v", req, err)
		return model.NewErrorCodeErr(model.ErrActivityCouponTooHot.Code, err)
	}
	if !distLockSuccess {
		ctx.Log.Errorf(ctx.Ctx, "领取优惠券，分布式锁获取失败：req=%+v,err=%+v", req, err)
		return model.NewErrorCodeErr(model.ErrActivityCouponTooHot.Code, err)
	}
	defer distLock.Release()
	//拿到锁再查一次，防止有其他请求已经处理过
	err = loadActivityFromRedis(ctx, req, activity)
	if err != nil || activity.Activity != nil {
		return
	}
	//第二步，加载活动数据
	activityData, err := model.FindActivityById(req.ActivityId)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		ctx.Log.Errorf(ctx.Ctx, "领取优惠券，从db查询活动数据失败：req=%+v,err=%+v", req, err)
		return model.NewErrorCodeErr(model.ErrActivityCouponTooHot.Code, err)
	}
	if errors.Is(err, gorm.ErrRecordNotFound) || activityData.Id == 0 {
		ctx.Log.Errorf(ctx.Ctx, "领取优惠券，活动数据不存在：req=%+v", req)
		return model.ErrActivityNotExists
	}
	couponListData, err := model.FindActivityCouponList(req.ActivityId)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		ctx.Log.Errorf(ctx.Ctx, "领取优惠券，从db查询优惠券数据失败：req=%+v,err=%+v", req, err)
		return model.NewErrorCodeErr(model.ErrActivityCouponTooHot.Code, err)
	}
	if len(couponListData) == 0 {
		ctx.Log.Errorf(ctx.Ctx, "领取优惠券，优惠券数据不存在：req=%+v", req)
		return model.ErrCouponNotInActivity
	}
	activity.Activity = activityData
	for i := 0; i < len(couponListData); i++ {
		coupon := couponListData[i]
		if coupon.TemplateCode == req.TemplateCode {
			activity.Coupon = coupon
			break
		}
	}
	activityCouponGroup, err := model.FindActivityCouponGroup(req.ActivityId, req.TemplateCode)
	err = config.Redis.Set(ctx.Ctx, activityCouponKey, utils.ToJsonString(activity), -1).Err()
	if err != nil {
		ctx.Log.Errorf(ctx.Ctx, "领取优惠券，存储活动优惠券配置数据到redis失败：req=%+v,err=%+v", req, err)
		return model.NewErrorCodeErr(model.ErrActivityCouponTooHot.Code, err)
	}
	if activityCouponGroup != nil && activityCouponGroup.Num > 0 {
		activity.Stock = activity.Coupon.Num - activityCouponGroup.Num
	}
	if activityData.UserRange == model.ActivityUserRangeImport {
		userImport, err := model.FindActivityUserImportList(req.ActivityId)
		if err != nil {
			ctx.Log.Errorf(ctx.Ctx, "领取优惠券，从db查询活动数据失败：req=%+v,err=%+v", req, err)
			return model.NewErrorCodeErr(model.ErrActivityCouponTooHot.Code, err)
		}
		var userIds []string
		for _, user := range userImport {
			userIds = append(userIds, fmt.Sprint(user.UserId))
		}
		activity.UserImports = "," + strings.Join(userIds, ",")
	}
	err = config.Redis.HMSet(ctx.Ctx, activityCouponKey, map[string]interface{}{
		"activity":    utils.ToJsonString(activity.Activity),
		"coupon":      utils.ToJsonString(activity.Coupon),
		"stock":       activity.Stock,
		"userImports": activity.UserImports,
	}).Err()
	if err != nil {
		ctx.Log.Errorf(ctx.Ctx, "领取优惠券，存储活动优惠券配置数据到redis失败：req=%+v,err=%+v", req, err)
		return model.ErrActivityCouponTooHot
	}
	return
}

func loadActivityFromRedis(ctx context.Context, req *UserGetCouponReq, activity *UserActivityCouponView) (err error) {
	ctx.Log.Infof(ctx.Ctx, "领取优惠券，step-loadActivityFromRedis:加载活动数据   req======%+v", req)
	activityCouponKey := fmt.Sprintf(activityCouponKeyFormat, req.ActivityId, req.TemplateCode)
	//userActivityCouponKey:=fmt.Sprintf(userActivityCouponKeyFormat, req.ActivityId, req.TemplateCode, req.UserId)
	activityCache, err := config.Redis.HMGet(ctx.Ctx, activityCouponKey, activityFields...).Result()
	if err != nil && err != redis.Nil {
		ctx.Log.Errorf(ctx.Ctx, "领取优惠券，从redis查询失败：req=%+v,err=%+v", req, err)
		return model.NewErrorCodeErr(model.ErrActivityCouponTooHot.Code, err)
	}
	if len(activityCache) == 4 {

		// Assuming activity data is stored in the first index
		if err := json.Unmarshal([]byte(activityCache[0].(string)), &activity.Activity); err != nil {
			ctx.Log.Errorf(ctx.Ctx, "领取优惠券，从redis反序列化活动数据失败：req=%+v,err=%+v", req, err)
			return model.NewErrorCodeErr(model.ErrActivityCouponTooHot.Code, err)
		}
		if err := json.Unmarshal([]byte(activityCache[1].(string)), &activity.Activity); err != nil {
			ctx.Log.Errorf(ctx.Ctx, "领取优惠券，从redis反序列化活动优惠券数据失败：req=%+v,err=%+v", req, err)
			return model.NewErrorCodeErr(model.ErrActivityCouponTooHot.Code, err)
		}
		// Assuming stock data is stored in the second index
		activity.Stock, err = strconv.Atoi(activityCache[2].(string))
		if err != nil {
			ctx.Log.Errorf(ctx.Ctx, "领取优惠券，从redis获取库存数据失败：req=%+v,err=%+v", req, err)
			return model.NewErrorCodeErr(model.ErrActivityCouponTooHot.Code, err)
		}
		// Assuming user imports data is stored in the third index
		activity.UserImports = activityCache[3].(string)

	}
	return nil
}

func checkActivityCoupon(ctx context.Context, req *UserGetCouponReq, activity *UserActivityCouponView) (err error) {
	ctx.Log.Infof(ctx.Ctx, "领取优惠券，step-checkActivityCoupon:验证活动规则   req======%+v", req)
	if activity.Activity.State == 2 {
		return model.ErrActivityStateOffline
	}
	if activity.Activity.AuditState == 2 {
		return model.ErrActivityAuditReject
	}

	now := time.Now().Unix()
	if activity.Activity.StartTime > now {
		return model.ErrActivityNotStart
	}
	if activity.Activity.EndTime < now {
		return model.ErrActivityAlreadyEnd
	}
	if activity.Activity.UserRange == model.ActivityUserRangeImport && !strings.Contains(activity.UserImports, ","+fmt.Sprint(req.UserId)) {
		return model.ErrUserNotInActivity
	}

	if activity.Stock <= 0 {
		return model.ErrActivityCouponNoStock
	}

	return nil
}

var userGetCouponScript = `local couponKey = KEYS[1]  -- 消费券库存键
local userClaimsKey = KEYS[2]  -- 用户领取记录键
local maxClaims = tonumber(ARGV[1])  -- 用户允许的最大领取数量
local acquireCount = tonumber(ARGV[2])  -- 用户请求领取的数量

-- 检查用户是否已达到最大领取数量
local currentClaims = tonumber(redis.call('GET', userClaimsKey) or 0)
if currentClaims >= maxClaims then
    return -1  -- 已达到最大领取数量
end

-- 检查库存
local stock = tonumber(redis.call('HGET', couponKey, 'stock') or 0)
if stock <= 0 then
    return 0  -- 库存已领完
end

-- 计算用户实际可领取的数量
local remainingClaims = maxClaims - currentClaims
local claimable = math.min(remainingClaims, acquireCount, stock)

-- 更新库存和用户领取记录
redis.call('HINCRBY', couponKey, 'stock', -claimable)
redis.call('INCRBY', userClaimsKey, claimable)

return claimable  -- 返回用户实际领取的数量`

func userGetCouponFromRedis(ctx context.Context, req *UserGetCouponReq, activity *UserActivityCouponView) (successNum int, err error) {
	ctx.Log.Infof(ctx.Ctx, "领取优惠券，step-userGetCouponFromRedis:执行从redis领取优惠券   req======%+v", req)
	couponKey := fmt.Sprintf(activityCouponKeyFormat, req.ActivityId, req.TemplateCode)
	//activityCouponKey := fmt.Sprintf("activity_coupon:%d:%s", req.ActivityId, req.TemplateCode)
	userClaimsKey := fmt.Sprintf(userActivityCouponKeyFormat, req.ActivityId, req.TemplateCode, req.UserId)
	//activityCouponGroupKey := fmt.Sprintf("activity_coupon_group:%d:%s", req.ActivityId, req.TemplateCode)

	result, err := config.Redis.Eval(ctx.Ctx, userGetCouponScript, []string{couponKey, userClaimsKey},
		[]string{fmt.Sprint(activity.Coupon.Num), fmt.Sprint(activity.Coupon.OneAccountNum)}).Result()
	if err != nil {
		ctx.Log.Errorf(ctx.Ctx, "领取优惠券，从redis领取优惠券失败：req=%+v,err=%+v", req, err)
		return 0, model.ErrActivityCouponTooHot
	}
	resultStr := result.(string)
	if resultStr == "0" {
		return 0, model.ErrActivityCouponNoStock
	}
	if resultStr == "-1" {
		return 0, model.ErrUserAlreadyGetCoupon
	}
	successNum, _ = strconv.Atoi(resultStr)
	return successNum, nil

}

func notifyGetCouponSuccess(ctx context.Context, req *UserGetCouponReq, successNum int) (err error) {
	ctx.Log.Infof(ctx.Ctx, "领取优惠券，step-notifyGetCouponSuccess:领取成功，执行通知   req======%+v", req)
	//var result = UserGetActivityCouponSuccess{
	//	ActivityId:   req.ActivityId,
	//	UserId:       req.UserId,
	//	TemplateCode: req.TemplateCode,
	//	SuccessNum:   successNum,
	//	Time:         time.Now().Unix(),
	//}
	return nil
}
