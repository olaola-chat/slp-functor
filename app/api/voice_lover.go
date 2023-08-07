package api

import (
	"github.com/gogf/gf/net/ghttp"
	"github.com/olaola-chat/rbp-functor/app/query"
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
// @Success 200 resp.EmptyResp{data=pb.RespVoiceLoverMain}
// @Router /go/banban/account/pushReport [get]
func (a *voiceLoverAPI) Main(r *ghttp.Request) {
	var req *query.ReqVoiceLoverMain
	if err := r.ParseQuery(&req); err != nil {
		response.Output(r, &pb.ResAccountPushReport{
			Success: false,
			Msg:     err.Error(),
		})
	}

	res := account.PushReportService.PushReport(r.Context(), req.Status)

	response.Output(r, res)
}