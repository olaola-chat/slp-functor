package api

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	vl_pb "github.com/olaola-chat/rbp-proto/gen_pb/rpc/voice_lover"
	vl_rpc "github.com/olaola-chat/rbp-proto/rpcclient/voice_lover"

	"github.com/olaola-chat/rbp-functor/app/consts"
	"github.com/olaola-chat/rbp-functor/app/pb"
	"github.com/olaola-chat/rbp-functor/app/query"
	vl_serv "github.com/olaola-chat/rbp-functor/app/service/voice_lover"
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
		OutputCustomError(r, consts.ERROR_PARAM)
		return
	}
	data, err := vl_serv.VoiceLoverService.GetMainData(r.GetCtx(), 1)
	if err != nil {
		g.Log().Errorf("voiceLoverAPI Main GetMainData error=%v", err)
		OutputCustomError(r, consts.ERROR_SYSTEM)
		return
	}
	OutputCustomData(r, data)
}

// AlbumList
// @Tags VoiceLover
// @Summary 获取更多专辑列表
// @Description 获取更多专辑列表
// @Accept application/json
// @Produce json
// @Security ApiKeyAuth,OAuth2Implicit
// @Request query.ReqAlbumList query
// @Success 200 {object} pb.RespAlbumList
// @Router /go/func/voice_lover/album_list [get]
func (a *voiceLoverAPI) AlbumList(r *ghttp.Request) {
	var req *query.ReqAlbumList
	if err := r.ParseQuery(&req); err != nil {
		OutputCustomError(r, consts.ERROR_PARAM)
		return
	}
	OutputCustomData(r, &pb.RespAlbumList{})
}

// Post
// @Tags VoiceLover
// @Summary 声恋投稿
// @Description 声恋投稿
// @Accept application/json
// @Produce json
// @Security ApiKeyAuth,OAuth2Implicit
// @Request query.ReqVoiceLoverPost query
// @Success 200 {object} pb.RespVoiceLoverMain
// @Router /go/func/voice_lover/post [post]
func (a *voiceLoverAPI) Post(r *ghttp.Request) {
	var req *query.ReqVoiceLoverPost
	if err := r.Parse(&req); err != nil {
		OutputCustomError(r, consts.ERROR_PARAM)
		return
	}
	ctx := r.Context()
	g.Log().Infof("VoiceLoverPost param = %v", *req)
	_, err := vl_rpc.VoiceLoverMain.Post(ctx, &vl_pb.ReqPost{
		Uid:         uint64(1),
		Resource:    req.Resource,
		Source:      req.Source,
		Cover:       req.Cover,
		Title:       req.Title,
		Desc:        req.Desc,
		EditDub:     req.EditDub,
		EditContent: req.EditContent,
		EditPost:    req.EditPost,
		EditCover:   req.EditCover,
		Labels:      req.Labels,
	})
	if err != nil {
		g.Log().Errorf("VoiceLover Post error, err = %v", err)
		OutputCustomError(r, consts.ERROR_SYSTEM)
		return
	}
	OutputCustomData(r, &pb.RespVoiceLoverPost{})
}
