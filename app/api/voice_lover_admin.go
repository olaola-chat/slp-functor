package api

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"

	"github.com/olaola-chat/rbp-functor/app/consts"
	"github.com/olaola-chat/rbp-functor/app/pb"
	"github.com/olaola-chat/rbp-functor/app/query"
	"github.com/olaola-chat/rbp-functor/app/service/voice_lover"
)

type voiceLoverAdminApi struct {
}

var VoiceLoverAdmin = &voiceLoverAdminApi{}

// AudioList
// @Tags VoiceLoverAdmin
// @Summary 声恋后台audio 列表
// @Description 声恋后台audio 列表
// @Accept application/json
// @Produce json
// @Security ApiKeyAuth,OAuth2Implicit
// @Request query.ReqAdminVoiceLoverAudioList query
// @Success 200 {object} pb.RespAdminVoiceLoverAudioList
// @Router /go/func/admin/voice_lover/audio-list [get]
func (a *voiceLoverAdminApi) AudioList(r *ghttp.Request) {
	var req *query.ReqAdminVoiceLoverAudioList
	if err := r.Parse(&req); err != nil {
		OutputCustomError(r, consts.ERROR_PARAM)
		return
	}
	ctx := r.Context()
	g.Log().Infof("VoiceLoverAdmin audio_list param = %v", *req)
	res, total, err := voice_lover.VoiceLoverService.GetAudioList(ctx, req)
	if err != nil {
		OutputCustomError(r, consts.ERROR_SYSTEM)
		return
	}
	data := &pb.RespAdminVoiceLoverAudioList{
		Audios: res,
		Total:  total,
	}
	OutputCustomData(r, data)
}
