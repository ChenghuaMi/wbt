package user

import (
	"context"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"wbt/client/errorx"
	"wbt/model"
	"wbt/util"
	"time"

	"wbt/client/internal/svc"
	"wbt/client/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req *types.LoginReq) (resp *types.CommonReply, err error) {

	user := model.User{}
	fmt.Println(req.Username,req.Password)
	result := l.svcCtx.GoDb.Where("username = ? and password = ?",req.Username,util.Md5(req.Password,"")).First(&user)
	fmt.Println("err:",result.Error,result.RowsAffected)
	err = result.Error
	if err != nil {
		fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")
		return nil,errorx.NewCodeDefaultError(errorx.UserNotFond)
	}
	iat := time.Now().Unix()
	auth := l.svcCtx.Config.Auth
	tokenString,err := l.getJwtToken(auth.AccessSecret,iat,auth.AccessExpire,user.Id)

	if err != nil {
		return nil,err
	}
	loginReply := types.LoginReply{
		Id:           user.Id,
		Username:     user.Username,
		TotalMoney:    user.TotalMoney,
		UseMoney:     user.UseMoney,
		Money:        user.Money,
		TotalScore:   user.TotalScore,
		UseScore:     user.UseScore,
		Score:        user.Score,
		AccessToken:  tokenString,
	}
	resp = &types.CommonReply{
		Code: 2000,
		Msg: "success",
		Data: loginReply,

	}
	return resp,nil
}

//生成token
func(l *LoginLogic) getJwtToken(secretKey string,iat,seconds,userId int64) (string,error){
	claims := make(jwt.MapClaims)
	claims["iat"] = iat
	claims["exp"] = iat + seconds
	claims["userId"] = userId
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = claims
	tokenString,err := token.SignedString([]byte(secretKey))
	return tokenString,err
}
