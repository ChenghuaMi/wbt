package sev

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"
	"wbt/server/internal/svc"
)

type RunServerLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRunServerLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RunServerLogic {
	return &RunServerLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RunServerLogic) RunServer() error {
	// todo: add your logic here and delete this line

	return nil
}
