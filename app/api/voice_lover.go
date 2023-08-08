package api

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/olaola-chat/rbp-library/response"
	voice_lover2 "github.com/olaola-chat/rbp-proto/gen_pb/rpc/voice_lover"
	"github.com/olaola-chat/rbp-proto/rpcclient/voice_lover"

	"github.com/olaola-chat/rbp-functor/app/pb"
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
// @Success 200 {object} pb.RespVoiceLoverMain
// @Router /go/func/voice_lover/main [get]
func (a *voiceLoverAPI) Main(r *ghttp.Request) {
	var req *query.ReqVoiceLoverMain
	if err := r.ParseQuery(&req); err != nil {
		response.Output(r, &pb.RespVoiceLoverMain{})
	}
	response.Output(r, &pb.RespVoiceLoverMain{})
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
		response.Output(r, &pb.RespAlbumList{})
	}
	response.Output(r, &pb.RespAlbumList{})
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
		response.Output(r, &pb.RespVoiceLoverPost{
			Msg: err.Error(),
		})
		return
	}
	ctx := r.Context()
	g.Log().Infof("VoiceLoverPost param = %v", *req)
	_, err := voice_lover.VoiceLoverMain.Post(ctx, &voice_lover2.ReqVoiceLoverPost{
		Uid:         uint64(1),
		Resource:    req.Resource,
		Source:      req.Source,
		Cover:       req.Cover,
		Title:       req.Title,
		Desc:        req.Desc,
		EditDub:     req.EditDub,
		EditContent: req.EditContent,
		EditPost:    req.EditPost,
		EditCover:   req.Cover,
		Labels:      req.Labels,
	})
	if err != nil {
		g.Log().Errorf("VoiceLover Post error, err = %v", err)
		response.Output(r, &pb.RespVoiceLoverPost{Msg: err.Error()})
		return
	}
	response.Output(r, &pb.RespVoiceLoverPost{Success: true})
}
