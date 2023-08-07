package api

import (
	"github.com/gogf/gf/net/ghttp"
	"github.com/olaola-chat/rbp-functor/app/pb"
	"github.com/olaola-chat/rbp-functor/app/query"
	"github.com/olaola-chat/rbp-library/response"
)

var VoiceLover = &voiceLoverAPI{}

type voiceLoverAPI struct {
}

// Main
// @Tags VoiceLover
// @Summary 声恋首页
// @Description 声恋首页
// @Accept application/json
// @Produce json
// @Security ApiKeyAuth,OAuth2Implicit
// @Request query.ReqVoiceLoverMain query
// @Success 200 {object} pb.RespVoiceLoverMain
// @Router /go/func/voice_lover/main [get]
func (a *voiceLoverAPI) Main(r *ghttp.Request) {
	var req *query.ReqVoiceLoverMain
	if err := r.ParseQuery(&req); err != nil {
		response.Output(r, &pb.RespVoiceLoverMain{})
	}
	response.Output(r, &pb.RespVoiceLoverMain{})
}
