package order

import (
	"context"
	"encoding/json"
	"wbt/client/errorx"
	"wbt/model"

	"wbt/client/internal/svc"
	"wbt/client/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type OrderListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewOrderListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OrderListLogic {
	return &OrderListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OrderListLogic) OrderList(req *types.OrderListReq) (resp *types.CommonReply, err error) {
	userId := l.ctx.Value("userId")
	uid,er := userId.(json.Number).Int64()
	if  er != nil {
		return nil,errorx.NewCodeDefaultError(errorx.UserNotLogin)
	}
	offset := (int(req.Page) - 1) * l.svcCtx.Config.PageSize
	var tmp []types.OrderList
	var count int64
	l.svcCtx.GoDb.Debug().Model(&model.Order{}).Where(map[string]interface{}{"user_id":uid,"type": req.Type}).Count(&count).Offset(offset).Limit(l.svcCtx.Config.PageSize).Order("created_time asc").Find(&tmp)
	orderListReply := types.OrderListReply{
		OrderList: tmp,
		PageList:  types.PageList{
			Total:    count,
			Page:     req.Page,
			PageSize: int64(l.svcCtx.Config.PageSize),
		},
	}
	resp = &types.CommonReply{
		Code: 2000,
		Msg:  "success",
		Data: orderListReply,
	}

	return
}
