package user

import (
	"wbt/client/errorx"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"wbt/client/internal/logic/user"
	"wbt/client/internal/svc"
	"wbt/client/internal/types"
)

func LoginHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.LoginReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.WriteJson(w,http.StatusOK,errorx.NewCodeDefaultError(err.Error()))
			//httpx.Error(w, err)
			return
		}

		l := user.NewLoginLogic(r.Context(), svcCtx)
		resp, err := l.Login(&req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
