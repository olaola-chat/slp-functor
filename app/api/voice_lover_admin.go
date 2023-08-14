package api

import (
	"sort"

	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/util/gconv"
	"github.com/olaola-chat/rbp-library/response"
	"github.com/olaola-chat/rbp-proto/gen_pb/db/xianshi"
	user2 "github.com/olaola-chat/rbp-proto/gen_pb/rpc/user"
	voice_lover3 "github.com/olaola-chat/rbp-proto/gen_pb/rpc/voice_lover"
	"github.com/olaola-chat/rbp-proto/rpcclient/user"
	voice_lover2 "github.com/olaola-chat/rbp-proto/rpcclient/voice_lover"

	"github.com/olaola-chat/rbp-functor/app/consts"
	"github.com/olaola-chat/rbp-functor/app/pb"
	"github.com/olaola-chat/rbp-functor/app/query"
	"github.com/olaola-chat/rbp-functor/app/service/voice_lover"
	"github.com/olaola-chat/rbp-functor/app/utils"
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
		response.Output(r, &pb.RespVoiceLoverPost{
			Success: false,
			Msg:     err.Error(),
		})
		return
	}
	ctx := r.Context()
	g.Log().Infof("VoiceLoverAdmin audio_list param = %v", *req)
	res, total, err := voice_lover.VoiceLoverService.GetAudioList(ctx, req)
	if err != nil {
		OutputCustomError(r, consts.ERROR_SYSTEM)
		return
	}
	response.Output(r, &pb.RespAdminVoiceLoverAudioList{
		Success: true,
		Audios:  res,
		Total:   total,
	})
	return
}

// AudioDetail
// @Tags VoiceLoverAdmin
// @Summary 声恋后台audio detail
// @Description 声恋后台audio detail
// @Accept application/json
// @Produce json
// @Security ApiKeyAuth,OAuth2Implicit
// @Request query.ReqAdminVoiceLoverAudioDetail query
// @Success 200 {object} pb.RespAdminVoiceLoverAudioDetail
// @Router /go/func/admin/voice_lover/audio-detail [get]
func (a *voiceLoverAdminApi) AudioDetail(r *ghttp.Request) {
	var req *query.ReqAdminVoiceLoverAudioDetail
	if err := r.Parse(&req); err != nil {
		response.Output(r, &pb.RespVoiceLoverPost{
			Success: false,
			Msg:     err.Error(),
		})
		return
	}
	ctx := r.Context()
	g.Log().Infof("VoiceLoverAdmin audio_detail param = %v", *req)
	reply, err := voice_lover2.VoiceLoverAdmin.GetAudioDetail(ctx, &voice_lover3.ReqGetAudioDetail{Id: req.Id})
	if err != nil {
		OutputCustomError(r, consts.ERROR_SYSTEM)
		return
	}
	if reply.Audio == nil {
		response.Output(r, &pb.RespAdminVoiceLoverAudioDetail{
			Success: true,
		})
		return
	}
	uids := make([]uint32, 0)
	for _, r := range reply.Audio.EditContents {
		uids = append(uids, r.Uid)
	}
	for _, r := range reply.Audio.EditDubs {
		uids = append(uids, r.Uid)
	}
	for _, r := range reply.Audio.EditContents {
		uids = append(uids, r.Uid)
	}
	for _, r := range reply.Audio.EditPosts {
		uids = append(uids, r.Uid)
	}
	uids = utils.DistinctUint32Slice(uids)
	userMap := make(map[uint32]*xianshi.EntityXsUserProfile)
	userReply, err := user.UserProfile.Mget(ctx, &user2.ReqUserProfiles{
		Uids:   uids,
		Fields: []string{"name", "uid", "icon"},
	})
	for _, u := range userReply.Data {
		userMap[u.Uid] = u
	}
	audio := &pb.AdminVoiceLoverAudio{
		Id:           reply.Audio.Id,
		CreateTime:   reply.Audio.CreateTime,
		PubUid:       reply.Audio.Uid,
		Resource:     reply.Audio.Resource,
		Covers:       reply.Audio.Covers,
		Source:       reply.Audio.Source,
		Desc:         reply.Audio.Desc,
		Labels:       reply.Audio.Labels,
		AuditStatus:  reply.Audio.AuditStatus,
		OpUid:        reply.Audio.OpUid,
		EditDubs:     buildAudioEdit(reply.Audio.EditDubs, userMap),
		EditContents: buildAudioEdit(reply.Audio.EditContents, userMap),
		EditPosts:    buildAudioEdit(reply.Audio.EditPosts, userMap),
		EditCovers:   buildAudioEdit(reply.Audio.EditCovers, userMap),
	}
	response.Output(r, &pb.RespAdminVoiceLoverAudioDetail{
		Success: true,
		Audio:   audio,
	})
	return
}

func buildAudioEdit(edits []*voice_lover3.AudioEditData, userMap map[uint32]*xianshi.EntityXsUserProfile) []*pb.AdminVoiceLoverAudioEdit {
	editDubs := make([]*pb.AdminVoiceLoverAudioEdit, 0)
	for _, e := range edits {
		editDubs = append(editDubs, &pb.AdminVoiceLoverAudioEdit{
			Uid:    e.Uid,
			Name:   userMap[e.Uid].GetName(),
			Avatar: userMap[e.Uid].GetIcon(),
		})
	}
	return editDubs
}

// AudioUpdate
// @Tags VoiceLoverAdmin
// @Summary 声恋后台audio 列表
// @Description 声恋后台audio 列表
// @Accept application/json
// @Produce json
// @Security ApiKeyAuth,OAuth2Implicit
// @Request query.ReqAdminVoiceLoverAudioUpdate query
// @Success 200 {object} pb.RespAdminVoiceLoverAudioUpdate
// @Router /go/func/admin/voice_lover/audio-update [post]
func (a *voiceLoverAdminApi) AudioUpdate(r *ghttp.Request) {
	var req *query.ReqAdminVoiceLoverAudioUpdate
	if err := r.Parse(&req); err != nil {
		response.Output(r, &pb.RespAdminVoiceLoverAudioUpdate{
			Msg: err.Error(),
		})
		return
	}
	ctx := r.Context()
	_, err := voice_lover2.VoiceLoverAdmin.UpdateAudio(ctx, &voice_lover3.ReqUpdateAudio{Id: req.Id, Title: req.Title, Desc: req.Desc, Labels: req.Labels})
	if err != nil {
		response.Output(r, &pb.RespAdminVoiceLoverAudioUpdate{
			Msg: err.Error(),
		})
		return
	}
	OutputCustomData(r, &pb.RespAdminVoiceLoverAudioUpdate{
		Success: true,
	})
}

// AudioAudit
// @Tags VoiceLoverAdmin
// @Summary 声恋后台audio 列表
// @Description 声恋后台audio 列表
// @Accept application/json
// @Produce json
// @Security ApiKeyAuth,OAuth2Implicit
// @Request query.ReqAdminVoiceLoverAudioAudit query
// @Success 200 {object} pb.RespAdminVoiceLoverAudioAudit {
// @Router /go/func/admin/voice_lover/audio-audit [post]
func (a *voiceLoverAdminApi) AudioAudit(r *ghttp.Request) {
	var req *query.ReqAdminVoiceLoverAudioAudit
	if err := r.Parse(&req); err != nil {
		response.Output(r, &pb.RespAdminVoiceLoverAudioAudit{
			Msg: err.Error(),
		})
		return
	}
	ctx := r.Context()
	_, err := voice_lover2.VoiceLoverAdmin.AuditAudio(ctx, &voice_lover3.ReqAuditAudio{Id: req.Id, AuditStatus: req.AuditStatus, AuditReason: req.AuditReason, OpUid: req.OpUid})
	if err != nil {
		response.Output(r, &pb.RespAdminVoiceLoverAudioAudit{
			Msg: err.Error(),
		})
		return
	}
	OutputCustomData(r, &pb.RespAdminVoiceLoverAudioAudit{
		Success: true,
	})
}

// AudioAuditReason
// @Tags VoiceLoverAdmin
// @Summary 声恋后台audio 列表
// @Description 声恋后台audio 列表
// @Accept application/json
// @Produce json
// @Security ApiKeyAuth,OAuth2Implicit
// @Request query
// @Success 200 {object} pb.RespAdminVoiceLoverAudioAuditReason
// @Router /go/func/admin/voice_lover/audio-audit-reason [get]
func (a *voiceLoverAdminApi) AudioAuditReason(r *ghttp.Request) {
	data := &pb.RespAdminVoiceLoverAudioAuditReason{}
	for id, reason := range consts.AuditAudioReasonMap {
		data.Reasons = append(data.Reasons, &pb.AdminVoiceLoverAudioAuditReason{
			Id:     id,
			Reason: reason,
		})
	}
	sort.Slice(data.Reasons, func(i, j int) bool {
		return data.Reasons[i].Id < data.Reasons[j].Id
	})
	OutputCustomData(r, &pb.RespAdminVoiceLoverAudioAuditReason{
		Success: true,
	})
}

// AlbumCreate
// @Tags VoiceLoverAdmin
// @Summary 声恋后台audio 列表
// @Description 声恋后台audio 列表
// @Accept application/json
// @Produce json
// @Security ApiKeyAuth,OAuth2Implicit
// @Request query.ReqAdminVoiceLoverAlbumCreate query
// @Success 200 {object} pb.RespAdminVoiceLoverAlbumCreate
// @Router /go/func/admin/voice_lover/album-update [post]
func (a *voiceLoverAdminApi) AlbumCreate(r *ghttp.Request) {
	var req *query.ReqAdminVoiceLoverAlbumCreate
	if err := r.Parse(&req); err != nil {
		OutputCustomData(r, &pb.RespAdminVoiceLoverAlbumCreate{
			Msg: err.Error(),
		})
		return
	}
	ctx := r.Context()
	reply, err := voice_lover2.VoiceLoverAdmin.CreateAlbum(ctx, &voice_lover3.ReqCreateAlbum{Name: req.Name, Intro: req.Intro, Cover: req.Cover, OpUid: req.OpUid})
	if err != nil {
		OutputCustomData(r, &pb.RespAdminVoiceLoverAlbumCreate{
			Msg: err.Error(),
		})
		return
	}
	OutputCustomData(r, &pb.RespAdminVoiceLoverAlbumCreate{
		Success: true,
		Id:      reply.Id,
	})
	return
}

// AlbumUpdate
// @Tags VoiceLoverAdmin
// @Summary 声恋后台audio 列表
// @Description 声恋后台audio 列表
// @Accept application/json
// @Produce json
// @Security ApiKeyAuth,OAuth2Implicit
// @Request query.ReqAdminVoiceLoverAlbumUpdate query
// @Success 200 {object}
// @Router /go/func/admin/voice_lover/album-update [post]
func (a *voiceLoverAdminApi) AlbumUpdate(r *ghttp.Request) {
	var req *query.ReqAdminVoiceLoverAlbumUpdate
	if err := r.Parse(&req); err != nil {
		OutputCustomData(r, &pb.RespAdminVoiceLoverAlbumUpdate{
			Msg: err.Error(),
		})
		return
	}
	ctx := r.Context()
	_, err := voice_lover2.VoiceLoverAdmin.UpdateAlbum(ctx, &voice_lover3.ReqUpdateAlbum{Id: req.Id, Name: req.Name, Intro: req.Intro, Cover: req.Cover, OpUid: req.OpUid})
	if err != nil {
		OutputCustomData(r, &pb.RespAdminVoiceLoverAlbumUpdate{
			Msg: err.Error(),
		})
		return
	}
	OutputCustomData(r, &pb.RespAdminVoiceLoverAlbumUpdate{
		Success: true,
	})
}

// AlbumDel
// @Tags VoiceLoverAdmin
// @Summary 声恋后台audio 列表
// @Description 声恋后台audio 列表
// @Accept application/json
// @Produce json
// @Security ApiKeyAuth,OAuth2Implicit
// @Request query.ReqAdminVoiceLoverAlbumDel query
// @Success 200 {object} pb.RespAdminVoiceLoverAlbumDel
// @Router /go/func/admin/voice_lover/album-del [post]
func (a *voiceLoverAdminApi) AlbumDel(r *ghttp.Request) {
	var req *query.ReqAdminVoiceLoverAlbumDel
	if err := r.Parse(&req); err != nil {
		OutputCustomData(r, &pb.RespAdminVoiceLoverAlbumDel{
			Msg: err.Error(),
		})
		return
	}
	ctx := r.Context()
	_, err := voice_lover2.VoiceLoverAdmin.DelAlbum(ctx, &voice_lover3.ReqDelAlbum{Id: req.Id, OpUid: req.OpUid})
	if err != nil {
		OutputCustomData(r, &pb.RespAdminVoiceLoverAlbumDel{
			Msg: err.Error(),
		})
		return
	}
	OutputCustomData(r, &pb.RespAdminVoiceLoverAlbumDel{
		Success: true,
	})
	return
}

// AlbumDetail
// @Tags VoiceLoverAdmin
// @Summary 声恋后台audio 列表
// @Description 声恋后台audio 列表
// @Accept application/json
// @Produce json
// @Security ApiKeyAuth,OAuth2Implicit
// @Request query.ReqAdminVoiceLoverAlbumDetail query
// @Success 200 {object} pb.RespAdminVoiceLoverAlbumDetail
// @Router /go/func/admin/voice_lover/album-detail [get]
func (a *voiceLoverAdminApi) AlbumDetail(r *ghttp.Request) {
	var req *query.ReqAdminVoiceLoverAlbumDetail
	if err := r.Parse(&req); err != nil {
		OutputCustomData(r, &pb.RespAdminVoiceLoverAlbumDetail{
			Msg: err.Error(),
		})
		return
	}
	ctx := r.Context()
	reply, err := voice_lover2.VoiceLoverAdmin.GetAlbumDetail(ctx, &voice_lover3.ReqGetAlbumDetail{AlbumStr: []string{gconv.String(req.Id)}})
	if err != nil {
		OutputCustomData(r, &pb.RespAdminVoiceLoverAlbumDetail{
			Msg: err.Error(),
		})
		return
	}
	if reply.Albums == nil || reply.Albums[req.Id] == nil {
		OutputCustomData(r, &pb.RespAdminVoiceLoverAlbumDetail{
			Success: true,
		})
		return
	}
	album := &pb.AdminVoiceLoverAlbum{
		Id:         reply.Albums[req.Id].Id,
		Name:       reply.Albums[req.Id].Name,
		Intro:      reply.Albums[req.Id].Intro,
		Cover:      reply.Albums[req.Id].Cover,
		AudioCount: int32(reply.Albums[req.Id].AudioCount),
		CreateTime: reply.Albums[req.Id].CreateTime,
		OpUid:      reply.Albums[req.Id].OpUid,
	}
	OutputCustomData(r, &pb.RespAdminVoiceLoverAlbumDetail{
		Success: true,
		Album:   album,
	})
}

// AlbumList
// @Tags VoiceLoverAdmin
// @Summary 声恋后台audio 列表
// @Description 声恋后台audio 列表
// @Accept application/json
// @Produce json
// @Security ApiKeyAuth,OAuth2Implicit
// @Request query.ReqAdminVoiceLoverAlbumList query
// @Success 200 {object} pb.RespAdminVoiceLoverAlbumList
// @Router /go/func/admin/voice_lover/album-list [get]
func (a *voiceLoverAdminApi) AlbumList(r *ghttp.Request) {
	var req *query.ReqAdminVoiceLoverAlbumList
	if err := r.Parse(&req); err != nil {
		OutputCustomData(r, &pb.RespAdminVoiceLoverAlbumList{
			Msg: err.Error(),
		})
		return
	}
	ctx := r.Context()
	reply, err := voice_lover2.VoiceLoverAdmin.GetAlbumList(ctx, &voice_lover3.ReqGetAlbumList{Name: req.Name, StartTime: req.StartTime, EndTime: req.EndTime,
		Limit: int32(req.Limit), Page: int32(req.Page)})
	if err != nil {
		OutputCustomData(r, &pb.RespAdminVoiceLoverAlbumList{
			Msg: err.Error(),
		})
		return
	}
	var albums []*pb.AdminVoiceLoverAlbum
	for _, l := range reply.Albums {
		albums = append(albums, &pb.AdminVoiceLoverAlbum{
			Id:         l.Id,
			Name:       l.Name,
			Intro:      l.Intro,
			Cover:      l.Cover,
			OpUid:      l.OpUid,
			AudioCount: int32(l.AudioCount),
			CreateTime: l.CreateTime,
		})
	}
	OutputCustomData(r, &pb.RespAdminVoiceLoverAlbumList{
		Success: true,
		Total:   reply.Total,
		Albums:  albums,
	})
}

// AudioCollectList
// @Tags VoiceLoverAdmin
// @Summary 声恋后台audio 列表
// @Description 声恋后台audio 列表
// @Accept application/json
// @Produce json
// @Security ApiKeyAuth,OAuth2Implicit
// @Request query.ReqAdminVoiceLoverAudioCollectList query
// @Success 200 {object} pb.RespAdminVoiceLoverAlbumList
// @Router /go/func/admin/voice_lover/audio-collect-list [get]
func (a *voiceLoverAdminApi) AudioCollectList(r *ghttp.Request) {
	var req *query.ReqAdminVoiceLoverAudioCollectList
	if err := r.Parse(&req); err != nil {
		OutputCustomData(r, &pb.RespAdminVoiceLoverAlbumList{
			Msg: err.Error(),
		})
		return
	}
	ctx := r.Context()
	g.Log().Infof("VoiceLoverAdmin audio_collect_list param = %v", *req)
	res, total, err := voice_lover.VoiceLoverService.GetAudioCollectList(ctx, req)
	if err != nil {
		OutputCustomData(r, &pb.RespAdminVoiceLoverAudioCollectList{
			Msg: err.Error(),
		})
		return
	}
	OutputCustomData(r, &pb.RespAdminVoiceLoverAudioCollectList{
		Success: true,
		Total:   total,
		Audios:  res,
	})
}

// AudioCollect
// @Tags VoiceLoverAdmin
// @Summary 声恋后台audio 列表
// @Description 声恋后台audio 列表
// @Accept application/json
// @Produce json
// @Security ApiKeyAuth,OAuth2Implicit
// @Request query.ReqAdminVoiceLoverAudioCollect query
// @Success 200 {object} pb.RespAdminVoiceLoverAudioCollect
// @Router /go/func/admin/voice_lover/audio-collect [post]
func (a *voiceLoverAdminApi) AudioCollect(r *ghttp.Request) {
	var req *query.ReqAdminVoiceLoverAudioCollect
	if err := r.Parse(&req); err != nil {
		OutputCustomData(r, &pb.RespAdminVoiceLoverAudioCollect{
			Msg: err.Error(),
		})
		return
	}
	ctx := r.Context()
	g.Log().Infof("VoiceLoverAdmin audio_collect_list param = %v", *req)
	_, err := voice_lover2.VoiceLoverAdmin.AudioCollect(ctx, &voice_lover3.ReqAudioCollect{
		AudioId: req.AudioId,
		AlbumId: req.AlbumId,
		Type:    req.Type,
	})
	if err != nil {
		OutputCustomData(r, &pb.RespAdminVoiceLoverAudioCollect{
			Msg: err.Error(),
		})
		return
	}
	OutputCustomData(r, &pb.RespAdminVoiceLoverAudioCollect{
		Success: true,
	})
}
