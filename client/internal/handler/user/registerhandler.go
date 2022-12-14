package user

import (
	"fmt"
	"wbt/client/errorx"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"wbt/client/internal/logic/user"
	"wbt/client/internal/svc"
	"wbt/client/internal/types"
)

func RegisterHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.RegisterReq
		if err := httpx.Parse(r, &req); err != nil {
			//httpx.Error(w, err)
			fmt.Println(r)
			httpx.WriteJson(w,http.StatusOK,errorx.NewCodeDefaultError(err.Error()))
			return
		}

		l := user.NewRegisterLogic(r.Context(), svcCtx)
		resp, err := l.Register(&req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
