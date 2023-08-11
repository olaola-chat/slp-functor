package api

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/olaola-chat/rbp-library/response"
	context2 "github.com/olaola-chat/rbp-library/server/http/context"
	"github.com/olaola-chat/rbp-library/server/http/middleware"
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
		response.Output(r, &pb.RespVoiceLoverMain{
			Success: false,
			Msg:     consts.ERROR_PARAM.Msg(),
		})
		return
	}
	ctxUser, _ := r.GetCtxVar(middleware.ContextUserKey).Interface().(*context2.ContextUser)
	g.Log().Debugf("ctxUser=%+v", ctxUser)
	data, err := vl_serv.VoiceLoverService.GetMainData(r.GetCtx(), ctxUser.UID)
	if err != nil {
		g.Log().Errorf("voiceLoverAPI Main GetMainData error=%v", err)
		response.Output(r, &pb.RespVoiceLoverMain{
			Success: false,
			Msg:     consts.ERROR_SYSTEM.Msg(),
		})
		return
	}
	response.Output(r, data)
}

// AlbumList
// @Tags VoiceLover
// @Summary 获取更多专辑列表 精选&非精选&专题都走该接口
// @Description 获取更多专辑列表 精选&非精选&专题都走该接口
// @Accept application/json
// @Produce json
// @Security ApiKeyAuth,OAuth2Implicit
// @Request query.ReqAlbumList query
// @Success 200 {object} pb.RespAlbumList
// @Router /go/func/voice_lover/albumList [get]
func (a *voiceLoverAPI) AlbumList(r *ghttp.Request) {
	var req *query.ReqAlbumList
	if err := r.ParseQuery(&req); err != nil {
		response.Output(r, &pb.RespAlbumList{
			Success: false,
			Msg:     consts.ERROR_PARAM.Msg(),
		})
		return
	}
	OutputCustomData(r, &pb.RespAlbumList{Success: true, Msg: ""})
}

// RecUserList
// @Tags VoiceLover
// @Summary 获取更多推荐用户接口
// @Description 获取更多推荐用户接口
// @Accept application/json
// @Produce json
// @Security ApiKeyAuth,OAuth2Implicit
// @Request query.ReqRecUserList query
// @Success 200 {object} pb.RespRecUserList
// @Router /go/func/voice_lover/recUserList [get]
func (a *voiceLoverAPI) RecUserList(r *ghttp.Request) {
	var req *query.ReqRecUserList
	if err := r.ParseQuery(&req); err != nil {
		response.Output(r, &pb.RespRecUserList{
			Success: false,
			Msg:     consts.ERROR_PARAM.Msg(),
		})
		return
	}
	OutputCustomData(r, &pb.RespRecUserList{Success: true, Msg: ""})
}

// AlbumDetail
// @Tags VoiceLover
// @Summary 查看专辑详情
// @Description 查看专辑详情
// @Accept application/json
// @Produce json
// @Security ApiKeyAuth,OAuth2Implicit
// @Request query.ReqAlbumDetail query
// @Success 200 {object} pb.RespAlbumDetail
// @Router /go/func/voice_lover/albumDetail [get]
func (a *voiceLoverAPI) AlbumDetail(r *ghttp.Request) {
	var req *query.ReqAlbumDetail
	if err := r.ParseQuery(&req); err != nil {
		response.Output(r, &pb.RespAlbumDetail{
			Success: false,
			Msg:     consts.ERROR_PARAM.Msg(),
		})
		return
	}
	OutputCustomData(r, &pb.RespAlbumDetail{Success: true, Msg: ""})
}

// AlbumComments
// @Tags VoiceLover
// @Summary 查看专辑评论列表
// @Description 查看专辑评论列表
// @Accept application/json
// @Produce json
// @Security ApiKeyAuth,OAuth2Implicit
// @Request query.ReqAlbumComments query
// @Success 200 {object} pb.RespAlbumComments
// @Router /go/func/voice_lover/albumComments [get]
func (a *voiceLoverAPI) AlbumComments(r *ghttp.Request) {
	var req *query.ReqAlbumComments
	if err := r.ParseQuery(&req); err != nil {
		response.Output(r, &pb.RespAlbumComments{
			Success: false,
			Msg:     consts.ERROR_PARAM.Msg(),
		})
		return
	}
	OutputCustomData(r, &pb.RespAlbumComments{Success: true, Msg: ""})
}

// CommentAlbum
// @Tags VoiceLover
// @Summary 发表专辑评论
// @Description 发表专辑评论
// @Accept application/json
// @Produce json
// @Security ApiKeyAuth,OAuth2Implicit
// @Request query.ReqCommentAlbum query
// @Success 200 {object} pb.RespCommentAlbum
// @Router /go/func/voice_lover/commentAlbum [post]
func (a *voiceLoverAPI) CommentAlbum(r *ghttp.Request) {
	var req *query.ReqCommentAlbum
	if err := r.ParseQuery(&req); err != nil {
		response.Output(r, &pb.RespCommentAlbum{
			Success: false,
			Msg:     consts.ERROR_PARAM.Msg(),
		})
		return
	}
	OutputCustomData(r, &pb.RespCommentAlbum{Success: true, Msg: ""})
}

// AudioDetail
// @Tags VoiceLover
// @Summary 查看音频详情
// @Description 查看音频详情
// @Accept application/json
// @Produce json
// @Security ApiKeyAuth,OAuth2Implicit
// @Request query.ReqAudioDetail query
// @Success 200 {object} pb.RespAudioDetail
// @Router /go/func/voice_lover/audioDetail [get]
func (a *voiceLoverAPI) AudioDetail(r *ghttp.Request) {
	var req *query.ReqAudioDetail
	if err := r.ParseQuery(&req); err != nil {
		response.Output(r, &pb.RespAudioDetail{
			Success: false,
			Msg:     consts.ERROR_PARAM.Msg(),
		})
		return
	}
	OutputCustomData(r, &pb.RespAudioDetail{Success: true, Msg: ""})
}

// AudioComments
// @Tags VoiceLover
// @Summary 查看音频评论&弹幕列表
// @Description 查看音频评论&弹幕列表
// @Accept application/json
// @Produce json
// @Security ApiKeyAuth,OAuth2Implicit
// @Request query.ReqAudioComments query
// @Success 200 {object} pb.RespAudioComments
// @Router /go/func/voice_lover/audioComments [get]
func (a *voiceLoverAPI) AudioComments(r *ghttp.Request) {
	var req *query.ReqAudioComments
	if err := r.ParseQuery(&req); err != nil {
		response.Output(r, &pb.RespAudioComments{
			Success: false,
			Msg:     consts.ERROR_PARAM.Msg(),
		})
		return
	}
	OutputCustomData(r, &pb.RespAudioComments{Success: true, Msg: ""})
}

// CommentAudio
// @Tags VoiceLover
// @Summary 发表音频评论
// @Description 发表音频评论
// @Accept application/json
// @Produce json
// @Security ApiKeyAuth,OAuth2Implicit
// @Request query.ReqCommentAudio query
// @Success 200 {object} pb.RespCommentAudio
// @Router /go/func/voice_lover/commentAudio [post]
func (a *voiceLoverAPI) CommentAudio(r *ghttp.Request) {
	var req *query.ReqCommentAudio
	if err := r.ParseQuery(&req); err != nil {
		response.Output(r, &pb.RespCommentAudio{
			Success: false,
			Msg:     consts.ERROR_PARAM.Msg(),
		})
		return
	}
	OutputCustomData(r, &pb.RespCommentAudio{Success: true, Msg: ""})
}

// Post
// @Tags VoiceLover
// @Summary 声恋投稿
// @Description 声恋投稿
// @Accept application/json
// @Produce json
// @Security ApiKeyAuth,OAuth2Implicit
// @Request query.ReqVoiceLoverPost query
// @Success 200 {object} pb.RespVoiceLoverPost
// @Router /go/func/voice_lover/post [post]
func (a *voiceLoverAPI) Post(r *ghttp.Request) {
	var req *query.ReqVoiceLoverPost
	if err := r.Parse(&req); err != nil {
		response.Output(r, &pb.RespVoiceLoverPost{
			Success: false,
			Msg:     consts.ERROR_PARAM.Msg(),
		})
		return
	}
	ctx := r.Context()
	g.Log().Infof("VoiceLoverPost param = %v", *req)
	_, err := vl_rpc.VoiceLoverMain.Post(ctx, &vl_pb.ReqPost{
		Uid:         101000002,
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
		Seconds:     req.Seconds,
	})
	if err != nil {
		g.Log().Errorf("VoiceLover Post error, err = %v", err)
		response.Output(r, &pb.RespVoiceLoverPost{
			Success: false,
			Msg:     consts.ERROR_SYSTEM.Msg(),
		})
		return
	}
	OutputCustomData(r, &pb.RespVoiceLoverPost{Success: true, Msg: ""})
}
