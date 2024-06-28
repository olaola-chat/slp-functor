package api

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/olaola-chat/rbp-functor/app/consts"
	"github.com/olaola-chat/rbp-functor/app/pb"
	"github.com/olaola-chat/rbp-functor/app/query"
	vl_serv "github.com/olaola-chat/rbp-functor/app/service/voice_lover"
	"github.com/olaola-chat/rbp-library/response"
	context2 "github.com/olaola-chat/rbp-library/server/http/context"
	"github.com/olaola-chat/rbp-library/tool"
	vl_pb "github.com/olaola-chat/rbp-proto/gen_pb/rpc/voice_lover"
	vl_rpc "github.com/olaola-chat/rbp-proto/rpcclient/voice_lover"
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
// @Param request body query.ReqVoiceLoverMain false "request"
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
	ctx := r.GetCtx()
	ctxUser := context2.ContextSrv.GetUserCtx(ctx)
	data, err := vl_serv.VoiceLoverService.GetMainData(ctx, ctxUser.UID)
	if err != nil {
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
// @Param request body query.ReqAlbumList false "request"
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
	data, err := vl_serv.VoiceLoverService.GetAlbumList(r.GetCtx(), req)
	if err != nil {
		response.Output(r, &pb.RespAlbumList{
			Success: false,
			Msg:     consts.ERROR_SYSTEM.Msg(),
		})
		return
	}
	response.Output(r, data)
}

// RecUserList
// @Tags VoiceLover
// @Summary 获取更多推荐用户接口
// @Description 获取更多推荐用户接口
// @Accept application/json
// @Produce json
// @Security ApiKeyAuth,OAuth2Implicit
// @Param request body query.ReqRecUserList false "request"
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
	response.Output(r, &pb.RespRecUserList{Success: true, Msg: ""})
}

// AlbumDetail
// @Tags VoiceLover
// @Summary 查看专辑详情
// @Description 查看专辑详情
// @Accept application/json
// @Produce json
// @Security ApiKeyAuth,OAuth2Implicit
// @Param request body query.ReqAlbumDetail false "request"
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
	ctx := r.GetCtx()
	ctxUser := context2.ContextSrv.GetUserCtx(ctx)
	data, err := vl_serv.VoiceLoverService.GetAlbumDetail(ctx, ctxUser.UID, req.AlbumId)
	if err != nil {
		response.Output(r, &pb.RespAlbumDetail{
			Success: false,
			Msg:     consts.ERROR_SYSTEM.Msg(),
		})
		return
	}
	response.Output(r, data)
}

// AlbumComments
// @Tags VoiceLover
// @Summary 查看专辑评论列表
// @Description 查看专辑评论列表
// @Accept application/json
// @Produce json
// @Security ApiKeyAuth,OAuth2Implicit
// @Param request body query.ReqAlbumComments false "request"
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
	ret := vl_serv.VoiceLoverService.GetAlbumCommentList(r.GetCtx(), req.AlbumId, req.Paginator.Page, req.Paginator.Limit)
	response.Output(r, ret)
}

// CommentAlbum
// @Tags VoiceLover
// @Summary 发表专辑评论
// @Description 发表专辑评论
// @Accept application/json
// @Produce json
// @Security ApiKeyAuth,OAuth2Implicit
// @Param request body query.ReqCommentAlbum false "request"
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

	ctx := r.GetCtx()
	ctxUser := context2.ContextSrv.GetUserCtx(ctx)
	postData := &vl_pb.ReqAlbumSubmitComment{
		AlbumId: req.AlbumId,
		Content: req.Comment,
		Uid:     ctxUser.UID,
	}
	region, err := tool.IP.GetAddr(r.RemoteAddr)
	if err == nil && region.Province != "" && region.Province != "0" {
		postData.Address = region.Province
	}
	ret := vl_serv.VoiceLoverService.SubmitAlbumComment(ctx, postData)
	response.Output(r, ret)
}

// AudioDetail
// @Tags VoiceLover
// @Summary 查看音频详情
// @Description 查看音频详情
// @Accept application/json
// @Produce json
// @Security ApiKeyAuth,OAuth2Implicit
// @Param request body query.ReqAudioDetail false "request"
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
	ctx := r.GetCtx()
	ctxUser := context2.ContextSrv.GetUserCtx(ctx)
	ret := vl_serv.VoiceLoverService.GetAudioDetail(ctx, ctxUser.UID, req.AudioId)
	response.Output(r, ret)
}

// AudioComments
// @Tags VoiceLover
// @Summary 查看音频评论&弹幕列表
// @Description 查看音频评论&弹幕列表
// @Accept application/json
// @Produce json
// @Security ApiKeyAuth,OAuth2Implicit
// @Param request body query.ReqAudioComments false "request"
// @Success 200 {object} pb.RespAudioComments
// @Router /go/func/voice_lover/audioComments [get]
func (a *voiceLoverAPI) AudioComments(r *ghttp.Request) {
	var req *query.ReqAudioComments
	if err := r.ParseQuery(&req); err != nil {
		response.Output(r, &pb.RespAudioComments{
			Success: false,
			Msg:     consts.ERROR_PARAM.Msg(),
			//Msg:     err.Error(),
		})
		return
	}
	ret := vl_serv.VoiceLoverService.GetAudioCommentList(r.GetCtx(), req.AudioId, req.Page, req.Limit)
	response.Output(r, ret)
}

// CommentAudio
// @Tags VoiceLover
// @Summary 发表音频评论
// @Description 发表音频评论
// @Accept application/json
// @Produce json
// @Security ApiKeyAuth,OAuth2Implicit
// @Param request body query.ReqCommentAudio false "request"
// @Success 200 {object} pb.RespCommentAudio
// @Router /go/func/voice_lover/commentAudio [post]
func (a *voiceLoverAPI) CommentAudio(r *ghttp.Request) {
	var req *query.ReqCommentAudio
	if err := r.ParseForm(&req); err != nil {
		response.Output(r, &pb.RespCommentAudio{
			Success: false,
			Msg:     consts.ERROR_PARAM.Msg(),
		})
		return
	}
	ctx := r.GetCtx()
	ctxUser := context2.ContextSrv.GetUserCtx(ctx)
	postData := &vl_pb.ReqAudioSubmitComment{
		AudioId: req.AudioId,
		Content: req.Comment,
		Uid:     ctxUser.UID,
		Type:    req.Type,
	}
	region, err := tool.IP.GetAddr(r.RemoteAddr)
	if err == nil && region.Province != "" && region.Province != "0" {
		postData.Address = region.Province
	}
	ret := vl_serv.VoiceLoverService.SubmitAudioComment(ctx, postData)

	response.Output(r, ret)
}

// Collect
// @Tags VoiceLover
// @Summary 收藏接口
// @Description 收藏接口
// @Accept application/json
// @Produce json
// @Security ApiKeyAuth,OAuth2Implicit
// @Param request body query.ReqCollectVoiceLover false "request"
// @Success 200 {object} pb.RespCollectVoiceLover
// @Router /go/func/voice_lover/collect [post]
func (a *voiceLoverAPI) Collect(r *ghttp.Request) {
	var req *query.ReqCollectVoiceLover
	if err := r.ParseQuery(&req); err != nil {
		response.Output(r, &pb.RespCollectVoiceLover{
			Success: false,
			Msg:     consts.ERROR_PARAM.Msg(),
		})
		return
	}
	ctx := r.GetCtx()
	ctxUser := context2.ContextSrv.GetUserCtx(ctx)
	_, err := vl_rpc.VoiceLoverMain.Collect(ctx, &vl_pb.ReqCollect{
		Uid:  ctxUser.UID,
		Id:   req.Id,
		Type: req.Type,
		From: req.From})
	if err != nil {
		g.Log().Errorf("voiceLoverAPI Collect error=%v", err)
		response.Output(r, &pb.RespCollectVoiceLover{
			Success: false,
			Msg:     consts.ERROR_SYSTEM.Msg(),
		})
		return
	}
	response.Output(r, &pb.RespCollectVoiceLover{Success: true, Msg: ""})
}

// CollectAlbumList
// @Tags VoiceLover
// @Summary 专辑收藏列表接口
// @Description 专辑收藏列表接口
// @Accept application/json
// @Produce json
// @Security ApiKeyAuth,OAuth2Implicit
// @Param request body query.ReqCollectAlbumList false "request"
// @Success 200 {object} pb.RespCollectAlbumList
// @Router /go/func/voice_lover/collectAlbumList [get]
func (a *voiceLoverAPI) CollectAlbumList(r *ghttp.Request) {
	var req *query.ReqCollectAlbumList
	if err := r.ParseQuery(&req); err != nil {
		response.Output(r, &pb.RespCollectAlbumList{
			Success: false,
			Msg:     consts.ERROR_PARAM.Msg(),
		})
		return
	}
	ctx := r.GetCtx()
	ctxUser := context2.ContextSrv.GetUserCtx(ctx)
	data, err := vl_serv.VoiceLoverService.GetCollectAlbumList(ctx, ctxUser.UID, req)
	if err != nil {
		response.Output(r, &pb.RespCollectAlbumList{
			Success: false,
			Msg:     consts.ERROR_SYSTEM.Msg(),
		})
		return
	}
	response.Output(r, data)
}

// CollectAudioList
// @Tags VoiceLover
// @Summary 音频收藏列表接口
// @Description 音频收藏列表接口
// @Accept application/json
// @Produce json
// @Security ApiKeyAuth,OAuth2Implicit
// @Param request body query.ReqCollectAudioList false "request"
// @Success 200 {object} pb.RespCollectAudioList
// @Router /go/func/voice_lover/collectAudioList [get]
func (a *voiceLoverAPI) CollectAudioList(r *ghttp.Request) {
	var req *query.ReqCollectAudioList
	if err := r.ParseQuery(&req); err != nil {
		response.Output(r, &pb.RespCollectAudioList{
			Success: false,
			Msg:     consts.ERROR_PARAM.Msg(),
		})
		return
	}
	ctx := r.GetCtx()
	ctxUser := context2.ContextSrv.GetUserCtx(ctx)
	data, err := vl_serv.VoiceLoverService.GetCollectAudioList(ctx, ctxUser.UID, req)
	if err != nil {
		response.Output(r, &pb.RespCollectAudioList{
			Success: false,
			Msg:     consts.ERROR_SYSTEM.Msg(),
		})
		return
	}
	response.Output(r, data)
}

// Post
// @Tags VoiceLover
// @Summary 声恋投稿
// @Description 声恋投稿
// @Accept application/json
// @Produce json
// @Security ApiKeyAuth,OAuth2Implicit
// @Param request body query.ReqVoiceLoverPost false "request"
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
	ctx := r.GetCtx()
	ctxUser := context2.ContextSrv.GetUserCtx(ctx)
	g.Log().Infof("VoiceLoverPost param = %v", *req)
	_, err := vl_rpc.VoiceLoverMain.Post(ctx, &vl_pb.ReqPost{
		Uid:         ctxUser.UID,
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
	response.Output(r, &pb.RespVoiceLoverPost{Success: true, Msg: ""})
}

// Report
// @Tags VoiceLover
// @Summary 评论举报
// @Description 评论举报
// @Accept application/json
// @Produce json
// @Security ApiKeyAuth,OAuth2Implicit
// @Param request body query.ReqReport false "request"
// @Success 200 {object} pb.CommonResp
// @Router /go/func/voice_lover/report [post]
func (a *voiceLoverAPI) Report(r *ghttp.Request) {
	var req *query.ReqReport
	if err := r.ParseQuery(&req); err != nil {
		response.Output(r, &pb.CommonResp{
			Success: false,
			Msg:     consts.ERROR_PARAM.Msg(),
		})
		return
	}
	//_, err := vl_rpc.VoiceLoverMain.UpdateReportStatus(r.GetCtx(), &vl_pb.ReqUpdateStatus{
	//	Id:   req.Id,
	//	Type: req.Type,
	//})
	//if err != nil {
	//	response.Output(r, &pb.RespVoiceLoverPost{
	//		Success: false,
	//		Msg:     consts.ERROR_PARAM.Msg(),
	//	})
	//}
	response.Output(r, &pb.CommonResp{Success: true})
}

// PlayStatReport
// @Tags VoiceLover
// @Summary 专辑播放事件上报
// @Description 专辑播放事件上报
// @Accept application/json
// @Produce json
// @Security ApiKeyAuth,OAuth2Implicit
// @Param request body query.ReqPlayStatReport false "request"
// @Success 200 {object} pb.RespPlayStatReport
// @Router /go/func/voice_lover/playStatReport [post]
func (a *voiceLoverAPI) PlayStatReport(r *ghttp.Request) {
	var req *query.ReqPlayStatReport
	if err := r.Parse(&req); err != nil {
		g.Log().Errorf("voiceLoverAPI PlayStatReport param error=%v", err)
		response.Output(r, &pb.RespPlayStatReport{
			Success: true,
			Msg:     "",
		})
		return
	}
	ctx := r.GetCtx()
	_, _ = vl_rpc.VoiceLoverMain.PlayStatReport(ctx, &vl_pb.ReqPlayStatReport{
		AlbumId: req.AlbumId,
		AudioId: req.AudioId,
	})
	response.Output(r, &pb.RespPlayStatReport{Success: true, Msg: ""})
}

// ShareAlbum
// @Tags VoiceLover
// @Summary 专辑分享接口
// @Description 专辑分享接口
// @Accept application/json
// @Produce json
// @Security ApiKeyAuth,OAuth2Implicit
// @Param request body query.ReqShareAlbum false "request"
// @Success 200 {object} pb.RespShareInfo
// @Router /go/func/voice_lover/shareAlbum [get]
func (a *voiceLoverAPI) ShareAlbum(r *ghttp.Request) {
	var req *query.ReqShareAlbum
	if err := r.ParseQuery(&req); err != nil {
		response.Output(r, &pb.RespShareInfo{
			Success: false,
			Msg:     consts.ERROR_PARAM.Msg(),
		})
		return
	}
	ctx := r.GetCtx()
	ctxUser := context2.ContextSrv.GetUserCtx(ctx)
	data, err := vl_serv.VoiceLoverService.ShareAlbumInfo(ctx, ctxUser.UID, req)
	if err != nil {
		response.Output(r, &pb.RespShareInfo{
			Success: false,
			Msg:     consts.ERROR_SYSTEM.Msg(),
		})
		return
	}
	response.Output(r, data)
}

// ShareAudio
// @Tags VoiceLover
// @Summary 音频分享接口
// @Description 音频分享接口
// @Accept application/json
// @Produce json
// @Security ApiKeyAuth,OAuth2Implicit
// @Param request body query.ReqShareAudio false "request"
// @Success 200 {object} pb.RespShareInfo
// @Router /go/func/voice_lover/shareAudio [get]
func (a *voiceLoverAPI) ShareAudio(r *ghttp.Request) {
	var req *query.ReqShareAudio
	if err := r.ParseQuery(&req); err != nil {
		response.Output(r, &pb.RespShareInfo{
			Success: false,
			Msg:     consts.ERROR_PARAM.Msg(),
		})
		return
	}
	ctx := r.GetCtx()
	ctxUser := context2.ContextSrv.GetUserCtx(ctx)
	data, err := vl_serv.VoiceLoverService.ShareAudioInfo(ctx, ctxUser.UID, req)
	if err != nil {
		response.Output(r, &pb.RespShareInfo{
			Success: false,
			Msg:     consts.ERROR_SYSTEM.Msg(),
		})
		return
	}
	response.Output(r, data)
}

// ShareAlbumFans
// @Tags VoiceLover
// @Summary 分享专辑通知粉丝接口
// @Description 分享专辑通知粉丝接口
// @Accept application/json
// @Produce json
// @Security ApiKeyAuth,OAuth2Implicit
// @Param request body query.ReqShareAlbumFans false "request"
// @Success 200 {object} pb.CommonResp
// @Router /go/func/voice_lover/shareAlbumFans [post]
func (a *voiceLoverAPI) ShareAlbumFans(r *ghttp.Request) {
	var req *query.ReqShareAlbumFans
	if err := r.ParseQuery(&req); err != nil {
		response.Output(r, &pb.CommonResp{
			Success: false,
			Msg:     consts.ERROR_PARAM.Msg(),
		})
		return
	}
	response.Output(r, &pb.CommonResp{
		Success: true,
	})
}

// ShareAudioFans
// @Tags VoiceLover
// @Summary 分享音频通知粉丝接口
// @Description 分享音频通知粉丝接口
// @Accept application/json
// @Produce json
// @Security ApiKeyAuth,OAuth2Implicit
// @Param request body query.ReqShareAudioFans false "request"
// @Success 200 {object} pb.CommonResp
// @Router /go/func/voice_lover/shareAudioFans [post]
func (a *voiceLoverAPI) ShareAudioFans(r *ghttp.Request) {
	var req *query.ReqShareAudioFans
	if err := r.ParseQuery(&req); err != nil {
		response.Output(r, &pb.CommonResp{
			Success: false,
			Msg:     consts.ERROR_PARAM.Msg(),
		})
		return
	}
	response.Output(r, &pb.CommonResp{
		Success: true,
	})
}

// ActivityMain
// @Tags VoiceLover
// @Summary 声恋挑战详情接口
// @Description 声恋挑战详情接口
// @Accept application/json
// @Produce json
// @Security ApiKeyAuth,OAuth2Implicit
// @Param request body query.ReqActivityMain false "request"
// @Success 200 {object} pb.RespVoiceLoverActivityMain
// @Router /go/func/voice_lover/activityMain [get]
func (a *voiceLoverAPI) ActivityMain(r *ghttp.Request) {
	var req query.ReqActivityMain
	if err := r.ParseQuery(&req); err != nil {
		response.Output(r, &pb.RespVoiceLoverActivityMain{Msg: consts.ERROR_PARAM.Msg()})
	}

	data, err := vl_serv.ActivitySrv.GetInfo(r.Context(), req.ActivityId)
	if err != nil {
		response.Output(r, &pb.RespVoiceLoverActivityMain{Msg: "服务开小差了，请稍后重试"})
	}
	response.Output(r, &pb.RespVoiceLoverActivityMain{Success: true, Data: data})
}
