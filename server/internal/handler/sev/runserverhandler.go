package sev

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"wbt/server/internal/logic/sev"
	"wbt/server/internal/svc"
)

func RunServerHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := sev.NewRunServerLogic(r.Context(), svcCtx)
		err := l.RunServer()
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.Ok(w)
		}
	}
}
