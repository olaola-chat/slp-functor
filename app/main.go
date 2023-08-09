package main

import (
	"context"

	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/olaola-chat/rbp-library/server/http"
	"github.com/olaola-chat/rbp-library/server/http/middleware"
	user_pb "github.com/olaola-chat/rbp-proto/gen_pb/rpc/user"
	user_rpc "github.com/olaola-chat/rbp-proto/rpcclient/user"

	"github.com/olaola-chat/rbp-functor/app/api"
)

func Auth(ctx context.Context, token string) (middleware.AuthUser, error) {
	respUserAuth, err := user_rpc.UserProfile.Auth(ctx, &user_pb.ReqUserAuth{Token: token})
	if err != nil {
		return middleware.AuthUser{}, gerror.New("auth error")
	}
	userData := middleware.AuthUser{
		UID:      respUserAuth.Uid,
		Time:     respUserAuth.Time,
		AppID:    uint8(respUserAuth.AppId),
		Salt:     respUserAuth.Salt,
		Platform: respUserAuth.Platform,
		Channel:  respUserAuth.Channel,
	}
	g.Log().Debugf("token=%s||userData=%+v", token, userData)
	return userData, nil
}

func route(server *ghttp.Server) {
	server.Group("/go/", func(group *ghttp.RouterGroup) {
		group.Middleware(
			middleware.Trace,
			middleware.CORS,                       //跨域请求
			middleware.Fire,                       //请求频率限制，简单的注入排除...
			middleware.NewCtxMiddleware(Auth).Ctx, //用户信息校验，多语言注入
		)
		group.Group("/func/", func(group *ghttp.RouterGroup) {
			group.Middleware(middleware.Auth) //登录校验
			group.Middleware(middleware.Error)
			group.ALL("voice_lover", api.VoiceLover)
		})
		group.Group("/func/admin/", func(group *ghttp.RouterGroup) {

		})
	})
}

func main() {
	http.Run(route)
}
