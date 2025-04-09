package service

import (
	"apiTest/config"
	"apiTest/model"
	context2 "context"
	"encoding/json"
	"errors"
	"github.com/redis/go-redis/v9"
	"github.com/yunduansing/gtools/context"
	"github.com/yunduansing/gtools/crypto"
	"github.com/yunduansing/gtools/utils"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"
	"strconv"
	"time"
)

type UserLoginService struct {
	ctx *context.Context
}

func NewUserLoginService(ctx *context.Context) *UserLoginService {
	return &UserLoginService{ctx: ctx}
}

type UserLoginDataCache struct {
	*UserLoginRes
	LoginTime    int64  `json:"loginTime"` //登录时间
	ClientType   string `json:"clientType"`
	LoginTimeout int64  `json:"loginTimeout"` //登录过期时间
}

var (
	loginTimeoutDuration = time.Hour * 24 * 30
	userLoginKeyPrefix   = "userLogin:Key:"  //登录缓存key--->  userLoginKey:token
	userLoginDataPrefix  = "userLogin:Data:" //登录用户信息--->  userLoginData:userId
)

func (s *UserLoginService) UserLogin(req *UserLoginReq) (res *UserLoginRes, err error) {
	user, err := model.FindUserByPhone(s.ctx, req.Phone)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		s.ctx.Log.Error(s.ctx.Ctx, "find user err:", err)
		return nil, model.NewMyError(500, "db error")
	}
	if user == nil || user.UserId == 0 {
		return nil, model.NewMyError(400, "user not exist")
	}
	//if user.State != model.UserStatusNormal {
	//	return nil, model.NewMyError(400, "user banned")
	//}
	token := crypto.GenHmacSha256(req.Phone, utils.UUID())
	res = &UserLoginRes{
		UserId:   user.UserId,
		Phone:    user.Phone,
		UserName: user.UserName,
		Token:    token,
	}
	err = s.saveUserLoginInfoToCache(res)
	if err != nil {
		return nil, err
	}
	//实际场景中还应该保存登录日志
	return
}

func (s *UserLoginService) saveUserLoginInfoToCache(userInfo *UserLoginRes) error {
	now := time.Now()
	var dataCache = UserLoginDataCache{
		UserLoginRes: userInfo,
		LoginTime:    now.Unix(),
		ClientType:   "1",
		LoginTimeout: now.Add(loginTimeoutDuration).Unix(),
	}
	var err error
	config.Redis.WrapDoWithTracing(
		s.ctx.Ctx, "redis.Set", func(ctx context2.Context, span trace.Span) error {
			_, err = config.Redis.Pipelined(
				ctx, func(pipe redis.Pipeliner) error {
					//首先查一下是否已经有登录信息，如果有就清除掉旧的登录信息
					userIdKey := userLoginDataPrefix + strconv.FormatInt(userInfo.UserId, 10)

					result, err := pipe.Get(ctx, userIdKey).Result()
					if err != nil && !errors.Is(err, redis.Nil) {
						s.ctx.Log.Error(ctx, "redis.Get loginData error:", err)
						return err
					}
					if len(result) > 0 {
						userOldTokenKey := userLoginKeyPrefix + result
						err = pipe.Del(ctx, userOldTokenKey).Err()
						if err != nil {
							s.ctx.Log.Error(ctx, "redis.Del loginData error:", err)
							return err
						}
					}
					err = pipe.Set(
						ctx, userLoginKeyPrefix+userInfo.Token, utils.ToJsonString(dataCache), loginTimeoutDuration,
					).Err()
					if err != nil {
						s.ctx.Log.Error(ctx, "redis.Set loginData error:", err)
						return err
					}

					err = pipe.Set(
						ctx, userIdKey, dataCache.Token,
						loginTimeoutDuration,
					).Err()
					if err != nil {
						s.ctx.Log.Error(ctx, "redis.Set loginData error:", err)
						return err
					}
					return nil
				},
			)
			return err
		},
	)
	if err != nil {
		return model.NewMyError(500, "redis error")
	}
	return nil
}

func (s *UserLoginService) GetUserLoginInfoByTokenFromCache(token string) (*UserLoginDataCache, error) {
	var userLoginData UserLoginDataCache
	var err error
	config.Redis.WrapDoWithTracing(
		s.ctx.Ctx, "redis.Get", func(ctx context2.Context, span trace.Span) error {
			var result string
			result, err = config.Redis.Get(ctx, userLoginKeyPrefix+token)
			if err != nil {
				s.ctx.Log.Error(ctx, "redis.Get loginData error:", err)
				return err
			}
			err = json.Unmarshal(utils.StringToByte(result), &userLoginData)
			if err != nil {
				s.ctx.Log.Error(ctx, "json.Unmarshal loginData from redis error:", err)
				return err
			}

			if userLoginData.LoginTimeout < time.Now().Unix() {
				err = model.NewMyError(401, "login timeout")
				s.ctx.Log.Error(ctx, "login timeout error:", err)
				return err
			}
			return nil
		},
	)
	return &userLoginData, err
}
