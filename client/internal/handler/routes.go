// Code generated by goctl. DO NOT EDIT.
package handler

import (
	"net/http"

	order "wbt/client/internal/handler/order"
	user "wbt/client/internal/handler/user"
	"wbt/client/internal/svc"

	"github.com/zeromicro/go-zero/rest"
)

func RegisterHandlers(server *rest.Server, serverCtx *svc.ServiceContext) {
	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodPost,
				Path:    "/login",
				Handler: user.LoginHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/register",
				Handler: user.RegisterHandler(serverCtx),
			},
		},
		rest.WithPrefix("/api/v1/user"),
	)

	server.AddRoutes(
		rest.WithMiddlewares(
			[]rest.Middleware{serverCtx.CheckToken},
			[]rest.Route{
				{
					Method:  http.MethodPost,
					Path:    "/createOrder",
					Handler: order.CreateOrderHandler(serverCtx),
				},
				{
					Method:  http.MethodGet,
					Path:    "/orderList",
					Handler: order.OrderListHandler(serverCtx),
				},
			}...,
		),
		rest.WithJwt(serverCtx.Config.Auth.AccessSecret),
		rest.WithPrefix("/api/v1/order"),
	)
}
