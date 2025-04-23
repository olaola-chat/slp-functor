package main

import (
	"github.com/gogf/gf/net/ghttp"
	"github.com/olaola-chat/slp-library/server/http"
	user_rpc "github.com/olaola-chat/slp-proto/rpcclient/user"

	"github.com/olaola-chat/slp-functor/app/api"

	"github.com/olaola-chat/slp-library/server/http/middleware"
)

//func Auth(ctx context.Context, token string) (middleware.AuthUser, error) {
//	respUserAuth, err := user_rpc.UserProfile.Auth(ctx, &user_pb.ReqUserAuth{Token: token})
//	if err != nil {
//		return middleware.AuthUser{}, gerror.New("auth error")
//	}
//	userData := middleware.AuthUser{
//		UID:      respUserAuth.Uid,
//		Time:     respUserAuth.Time,
//		AppID:    uint8(respUserAuth.AppId),
//		Salt:     respUserAuth.Salt,
//		Platform: respUserAuth.Platform,
//		Channel:  respUserAuth.Channel,
//	}
//	return userData, nil
//}

func route(server *ghttp.Server) {
	server.Group("/go/", func(group *ghttp.RouterGroup) {
		group.Middleware(
			middleware.Trace,
			middleware.CORS, //跨域请求
			//middleware.Fire, //请求频率限制，简单的注入排除...
			// middleware.NewCtxMiddleware(user_rpc.Auth).Ctx, //用户信息校验，多语言注入
			middleware.NewCtxMiddleware(user_rpc.Auth2).Ctx, //用户信息校验，多语言注入
		)
		group.Group("/func/", func(group *ghttp.RouterGroup) {
			group.Middleware(middleware.Auth) //登录校验
			group.Middleware(middleware.Error)
			group.ALL("voice_lover", api.VoiceLover)
		})
	})
}

func main() {
	http.Run(route)
}
