// Code generated by goctl. DO NOT EDIT.
package handler

import (
	"net/http"

	sev "wbt/server/internal/handler/sev"
	"wbt/server/internal/svc"

	"github.com/zeromicro/go-zero/rest"
)

func RegisterHandlers(server *rest.Server, serverCtx *svc.ServiceContext) {
	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodPost,
				Path:    "/runServer",
				Handler: sev.RunServerHandler(serverCtx),
			},
		},
		rest.WithPrefix("/bak"),
	)
}
