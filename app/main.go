package main

import (
	"github.com/gogf/gf/net/ghttp"
	"github.com/olaola-chat/rbp-functor/app/api"
	"github.com/olaola-chat/rbp-library/server/http"
)

func route(server *ghttp.Server) {
	server.Group("/go/", func(group *ghttp.RouterGroup) {
		// 		group.Middleware(
		// 			service.Middleware.Trace,
		// 			service.Middleware.CORS, //跨域请求
		// 			service.Middleware.Fire, //请求频率限制，简单的注入排除...
		// 			service.Middleware.Ctx,  //用户信息校验，多语言注入
		// 		)
		//这里的是不需要登录验证的
		group.Group("/func/", func(group *ghttp.RouterGroup) {
			// group.Middleware(service.Middleware.Error)
			group.ALL("/voice_lover", api.VoiceLover{})
		})
	})
}

func main() {
	http.Run(route)
}
