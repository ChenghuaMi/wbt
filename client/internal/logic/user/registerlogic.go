package user

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"wbt/client/errorx"
	"wbt/client/internal/svc"
	"wbt/client/internal/types"
	"wbt/model"
	"wbt/util"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RegisterLogic) Register(req *types.RegisterReq) (resp *types.RegisterReply, err error) {
	user := model.User{}
	result := l.svcCtx.GoDb.Where("username = ?",req.Username).First(&user)
	fmt.Println(result.Error,result.RowsAffected)
	if result.Error == gorm.ErrRecordNotFound {
		fmt.Println(">>>>")
		user.Username = req.Username
		user.Password = util.Md5(req.Password,"")
		now := time.Now().Unix()
		user.CreatedTime = now
		user.UpdatedTime = now
		result := l.svcCtx.GoDb.Create(&user)
		if result.RowsAffected > 0 {
			resp = &types.RegisterReply{
				Code: 2000,
				Msg: errorx.RegisterSuccess,
				Ok:  true,
			}
			return
		}
		err = result.Error
		return nil,err
	}
	return nil,errorx.NewCodeDefaultError(errorx.UserExist)




}
