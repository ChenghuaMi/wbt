package order

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
	"wbt/client/errorx"
	"wbt/client/internal/svc"
	"wbt/client/internal/types"
	"wbt/model"
	"wbt/util"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateOrderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateOrderLogic {
	return &CreateOrderLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateOrderLogic) CreateOrder(req *types.CreateOrderReq) (resp *types.CommonReply, err error) {
	order := model.Order{}
	user := model.User{}
	userId := l.ctx.Value("userId")
	uid,_ := userId.(json.Number).Int64()
	l.svcCtx.GoDb.Debug().Where("id = ?",uid).First(&user)
	total := req.Price * float64(req.Num)
	if total >= user.Money {
		return nil,errorx.NewCodeDefaultError(errorx.MoneyNotEnough)
	}
	orderNo := util.CreateOrder(util.OrderLength)
	result := l.svcCtx.GoDb.Debug().Where("order_no = ?",orderNo).First(&order)

	if result.Error == gorm.ErrRecordNotFound {
		now := time.Now().Unix()
		order.Num = req.Num
		order.Price = req.Price
		order.OrderNo = orderNo
		order.Type = req.Type
		order.CreatedTime = now
		order.UpdatedTime = now
		order.UserId = uid

		tx := l.svcCtx.GoDb.Begin()
		if err := tx.Create(&order).Error;err != nil {
			tx.Rollback()
			return nil,errorx.NewCodeDefaultError(err.Error())
		}
		money := user.Money - total
		if money < 0 {
			money = 0
		}
		if err = tx.Model(&model.User{Id: uid}).Updates(map[string]interface{}{"Money":money,"UseMoney": user.UseMoney + total,"UpdatedTime":now}).Error;err != nil {
			tx.Rollback()
			return nil,errorx.NewCodeDefaultError(err.Error())
		}
		tx.Commit()
		key := "buy"
		if req.Type == 2 {
			key = "sell"
		}
		byt,_ := util.JsonData(order)
		member := &redis.Z{
			Score:  float64(now),
			Member: byt,
		}
		l.svcCtx.Redis.C.ZAdd(l.svcCtx.Redis.Ctx,key,member)

		createRepl := types.CreateOrderReply{OrderId: order.Id}
		resp = &types.CommonReply{
			Code: 2000,
			Msg:  "success",
			Data: createRepl,
		}
		return resp,nil

	}

	return nil,errorx.NewCodeDefaultError(errorx.OrderExist)
}