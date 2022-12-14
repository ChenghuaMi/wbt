package order

import (
	"wbt/client/errorx"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"wbt/client/internal/logic/order"
	"wbt/client/internal/svc"
	"wbt/client/internal/types"
)

func CreateOrderHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CreateOrderReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.WriteJson(w,http.StatusOK,errorx.NewCodeDefaultError(err.Error()))
			return
		}

		l := order.NewCreateOrderLogic(r.Context(), svcCtx)
		resp, err := l.CreateOrder(&req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
