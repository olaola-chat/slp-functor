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
// @Param request body  query.ReqAdminVoiceLoverAudioList false "request"
// @Success 200 {object} pb.RespAdminVoiceLoverAudioList
// @Router /go/func/admin/voice_lover/audioList [get]
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
		response.Output(r, consts.ERROR_SYSTEM)
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
// @Param request body query.ReqAdminVoiceLoverAudioDetail false "request"
// @Success 200 {object} pb.RespAdminVoiceLoverAudioDetail
// @Router /go/func/admin/voice_lover/audioDetail [get]
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
		response.Output(r, consts.ERROR_SYSTEM)
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
// @Param request body query.ReqAdminVoiceLoverAudioUpdate false "request"
// @Success 200 {object} pb.RespAdminVoiceLoverAudioUpdate
// @Router /go/func/admin/voice_lover/audioUpdate [post]
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
	response.Output(r, &pb.RespAdminVoiceLoverAudioUpdate{
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
// @Param request body  query.ReqAdminVoiceLoverAudioAudit false "request"
// @Success 200 {object} pb.RespAdminVoiceLoverAudioAudit
// @Router /go/func/admin/voice_lover/audioAudit [post]
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
	response.Output(r, &pb.RespAdminVoiceLoverAudioAudit{
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
// @Success 200 {object} pb.RespAdminVoiceLoverAudioAuditReason
// @Router /go/func/admin/voice_lover/audioAuditReason [get]
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
	response.Output(r, &pb.RespAdminVoiceLoverAudioAuditReason{
		Success: true,
		Reasons: data.Reasons,
	})
}

// AlbumCreate
// @Tags VoiceLoverAdmin
// @Summary 声恋后台audio 列表
// @Description 声恋后台audio 列表
// @Accept application/json
// @Produce json
// @Security ApiKeyAuth,OAuth2Implicit
// @Param request body  query.ReqAdminVoiceLoverAlbumCreate false "request"
// @Success 200 {object} pb.RespAdminVoiceLoverAlbumCreate
// @Router /go/func/admin/voice_lover/albumCreate [post]
func (a *voiceLoverAdminApi) AlbumCreate(r *ghttp.Request) {
	var req *query.ReqAdminVoiceLoverAlbumCreate
	if err := r.Parse(&req); err != nil {
		response.Output(r, &pb.RespAdminVoiceLoverAlbumCreate{
			Msg: err.Error(),
		})
		return
	}
	ctx := r.Context()
	reply, err := voice_lover2.VoiceLoverAdmin.CreateAlbum(ctx, &voice_lover3.ReqCreateAlbum{Name: req.Name, Intro: req.Intro, Cover: req.Cover, OpUid: req.OpUid})
	if err != nil {
		response.Output(r, &pb.RespAdminVoiceLoverAlbumCreate{
			Msg: err.Error(),
		})
		return
	}
	response.Output(r, &pb.RespAdminVoiceLoverAlbumCreate{
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
// @Param request body query.ReqAdminVoiceLoverAlbumUpdate false "request"
// @Success 200 {object} pb.RespAdminVoiceLoverAlbumUpdate
// @Router /go/func/admin/voice_lover/albumUpdate [post]
func (a *voiceLoverAdminApi) AlbumUpdate(r *ghttp.Request) {
	var req *query.ReqAdminVoiceLoverAlbumUpdate
	if err := r.Parse(&req); err != nil {
		response.Output(r, &pb.RespAdminVoiceLoverAlbumUpdate{
			Msg: err.Error(),
		})
		return
	}
	ctx := r.Context()
	_, err := voice_lover2.VoiceLoverAdmin.UpdateAlbum(ctx, &voice_lover3.ReqUpdateAlbum{Id: req.Id, Name: req.Name, Intro: req.Intro, Cover: req.Cover, OpUid: req.OpUid})
	if err != nil {
		response.Output(r, &pb.RespAdminVoiceLoverAlbumUpdate{
			Msg: err.Error(),
		})
		return
	}
	response.Output(r, &pb.RespAdminVoiceLoverAlbumUpdate{
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
// @Param request body query.ReqAdminVoiceLoverAlbumDel false "request"
// @Success 200 {object} pb.RespAdminVoiceLoverAlbumDel
// @Router /go/func/admin/voice_lover/albumDel [post]
func (a *voiceLoverAdminApi) AlbumDel(r *ghttp.Request) {
	var req *query.ReqAdminVoiceLoverAlbumDel
	if err := r.Parse(&req); err != nil {
		response.Output(r, &pb.RespAdminVoiceLoverAlbumDel{
			Msg: err.Error(),
		})
		return
	}
	ctx := r.Context()
	_, err := voice_lover2.VoiceLoverAdmin.DelAlbum(ctx, &voice_lover3.ReqDelAlbum{Id: req.Id, OpUid: req.OpUid})
	if err != nil {
		response.Output(r, &pb.RespAdminVoiceLoverAlbumDel{
			Msg: err.Error(),
		})
		return
	}
	response.Output(r, &pb.RespAdminVoiceLoverAlbumDel{
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
// @Param request body  query.ReqAdminVoiceLoverAlbumDetail false "request"
// @Success 200 {object} pb.RespAdminVoiceLoverAlbumDetail
// @Router /go/func/admin/voice_lover/albumDetail [get]
func (a *voiceLoverAdminApi) AlbumDetail(r *ghttp.Request) {
	var req *query.ReqAdminVoiceLoverAlbumDetail
	if err := r.Parse(&req); err != nil {
		response.Output(r, &pb.RespAdminVoiceLoverAlbumDetail{
			Msg: err.Error(),
		})
		return
	}
	ctx := r.Context()
	reply, err := voice_lover2.VoiceLoverAdmin.GetAlbumDetail(ctx, &voice_lover3.ReqGetAlbumDetail{AlbumStr: []string{gconv.String(req.Id)}})
	if err != nil {
		response.Output(r, &pb.RespAdminVoiceLoverAlbumDetail{
			Msg: err.Error(),
		})
		return
	}
	if reply.Albums == nil || reply.Albums[req.Id] == nil {
		response.Output(r, &pb.RespAdminVoiceLoverAlbumDetail{
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
	response.Output(r, &pb.RespAdminVoiceLoverAlbumDetail{
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
// @Param request body query.ReqAdminVoiceLoverAlbumList false "request"
// @Success 200 {object} pb.RespAdminVoiceLoverAlbumList
// @Router /go/func/admin/voice_lover/albumList [get]
func (a *voiceLoverAdminApi) AlbumList(r *ghttp.Request) {
	var req *query.ReqAdminVoiceLoverAlbumList
	if err := r.Parse(&req); err != nil {
		response.Output(r, &pb.RespAdminVoiceLoverAlbumList{
			Msg: err.Error(),
		})
		return
	}
	ctx := r.Context()
	reply, err := voice_lover2.VoiceLoverAdmin.GetAlbumList(ctx, &voice_lover3.ReqGetAlbumList{Name: req.Name, StartTime: req.StartTime, EndTime: req.EndTime,
		Limit: int32(req.Limit), Page: int32(req.Page)})
	if err != nil {
		response.Output(r, &pb.RespAdminVoiceLoverAlbumList{
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
	response.Output(r, &pb.RespAdminVoiceLoverAlbumList{
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
// @Param  request body query.ReqAdminVoiceLoverAudioCollectList false "request"
// @Success 200 {object} pb.RespAdminVoiceLoverAudioCollectList
// @Router /go/func/admin/voice_lover/audioCollectList [get]
func (a *voiceLoverAdminApi) AudioCollectList(r *ghttp.Request) {
	var req *query.ReqAdminVoiceLoverAudioCollectList
	if err := r.Parse(&req); err != nil {
		response.Output(r, &pb.RespAdminVoiceLoverAlbumList{
			Msg: err.Error(),
		})
		return
	}
	ctx := r.Context()
	g.Log().Infof("VoiceLoverAdmin audio_collect_list param = %v", *req)
	res, total, err := voice_lover.VoiceLoverService.GetAudioCollectList(ctx, req)
	if err != nil {
		response.Output(r, &pb.RespAdminVoiceLoverAudioCollectList{
			Msg: err.Error(),
		})
		return
	}
	response.Output(r, &pb.RespAdminVoiceLoverAudioCollectList{
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
// @Param request body  query.ReqAdminVoiceLoverAudioCollect false "request"
// @Success 200 {object} pb.RespAdminVoiceLoverAudioCollect
// @Router /go/func/admin/voice_lover/audioCollect [post]
func (a *voiceLoverAdminApi) AudioCollect(r *ghttp.Request) {
	var req *query.ReqAdminVoiceLoverAudioCollect
	if err := r.Parse(&req); err != nil {
		response.Output(r, &pb.RespAdminVoiceLoverAudioCollect{
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
		response.Output(r, &pb.RespAdminVoiceLoverAudioCollect{
			Msg: err.Error(),
		})
		return
	}
	response.Output(r, &pb.RespAdminVoiceLoverAudioCollect{
		Success: true,
	})
}

// SubjectCreate
// @Tags VoiceLoverAdmin
// @Summary 声恋后台audio 列表
// @Description 声恋后台audio 列表
// @Accept application/json
// @Produce json
// @Security ApiKeyAuth,OAuth2Implicit
// @Param request body query.ReqAdminVoiceLoverSubjectCreate false "request"
// @Success 200 {object} pb.RespAdminVoiceLoverSubjectCreate
// @Router /go/func/admin/voice_lover/subjectCreate [post]
func (a *voiceLoverAdminApi) SubjectCreate(r *ghttp.Request) {
	var req *query.ReqAdminVoiceLoverSubjectCreate
	if err := r.Parse(&req); err != nil {
		response.Output(r, &pb.RespAdminVoiceLoverSubjectCreate{
			Msg: err.Error(),
		})
		return
	}
	ctx := r.Context()
	reply, err := voice_lover2.VoiceLoverAdmin.CreateSubject(ctx, &voice_lover3.ReqCreateSubject{
		Name:  req.Name,
		OpUid: req.OpUid,
	})
	if err != nil {
		response.Output(r, &pb.RespAdminVoiceLoverSubjectCreate{
			Msg: err.Error(),
		})
		return
	}
	response.Output(r, &pb.RespAdminVoiceLoverSubjectCreate{
		Success: true,
		Id:      reply.Id,
	})
}

// SubjectUpdate
// @Tags VoiceLoverAdmin
// @Summary 声恋后台audio 列表
// @Description 声恋后台audio 列表
// @Accept application/json
// @Produce json
// @Security ApiKeyAuth,OAuth2Implicit
// @Param request body query.ReqAdminVoiceLoverSubjectUpdate false "request"
// @Success 200 {object} pb.RespAdminVoiceLoverSubjectUpdate
// @Router /go/func/admin/voice_lover/subjectUpdate [post]
func (a *voiceLoverAdminApi) SubjectUpdate(r *ghttp.Request) {
	var req *query.ReqAdminVoiceLoverSubjectUpdate
	if err := r.Parse(&req); err != nil {
		response.Output(r, &pb.RespAdminVoiceLoverSubjectUpdate{
			Msg: err.Error(),
		})
		return
	}
	ctx := r.Context()
	_, err := voice_lover2.VoiceLoverAdmin.UpdateSubject(ctx, &voice_lover3.ReqUpdateSubject{
		Id:    req.Id,
		Name:  req.Name,
		OpUid: req.OpUid,
	})
	if err != nil {
		response.Output(r, &pb.RespAdminVoiceLoverSubjectUpdate{
			Msg: err.Error(),
		})
		return
	}
	response.Output(r, &pb.RespAdminVoiceLoverSubjectUpdate{
		Success: true,
	})
}

// SubjectDel
// @Tags VoiceLoverAdmin
// @Summary 声恋后台audio 列表
// @Description 声恋后台audio 列表
// @Accept application/json
// @Produce json
// @Security ApiKeyAuth,OAuth2Implicit
// @Param request body  query.ReqAdminVoiceLoverSubjectDel false "request"
// @Success 200 {object} pb.RespAdminVoiceLoverSubjectDel
// @Router /go/func/admin/voice_lover/subjectDel [post]
func (a *voiceLoverAdminApi) SubjectDel(r *ghttp.Request) {
	var req *query.ReqAdminVoiceLoverSubjectDel
	if err := r.Parse(&req); err != nil {
		response.Output(r, &pb.RespAdminVoiceLoverSubjectDel{
			Msg: err.Error(),
		})
		return
	}
	ctx := r.Context()
	_, err := voice_lover2.VoiceLoverAdmin.DelSubject(ctx, &voice_lover3.ReqDelSubject{
		Id:    req.Id,
		OpUid: req.OpUid,
	})
	if err != nil {
		response.Output(r, &pb.RespAdminVoiceLoverSubjectDel{
			Msg: err.Error(),
		})
		return
	}
	response.Output(r, &pb.RespAdminVoiceLoverSubjectDel{
		Success: true,
	})
}

// SubjectList
// @Tags VoiceLoverAdmin
// @Summary 声恋后台audio 列表
// @Description 声恋后台audio 列表
// @Accept application/json
// @Produce json
// @Security ApiKeyAuth,OAuth2Implicit
// @Param request body  query.ReqAdminVoiceLoverSubjectList false "request"
// @Success 200 {object} pb.RespAdminVoiceLoverSubjectList
// @Router /go/func/admin/voice_lover/subjectList [get]
func (a *voiceLoverAdminApi) SubjectList(r *ghttp.Request) {
	var req *query.ReqAdminVoiceLoverSubjectList
	if err := r.Parse(&req); err != nil {
		response.Output(r, &pb.RespAdminVoiceLoverSubjectList{
			Msg: err.Error(),
		})
		return
	}
	ctx := r.Context()
	reply, err := voice_lover2.VoiceLoverAdmin.GetSubjectList(ctx, &voice_lover3.ReqGetSubjectList{
		Name:      req.Name,
		Page:      int32(req.Page),
		Limit:     int32(req.Limit),
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
	})
	if err != nil {
		response.Output(r, &pb.RespAdminVoiceLoverSubjectList{
			Msg: err.Error(),
		})
		return
	}
	list := make([]*pb.SubjectData, 0)
	for _, subject := range reply.Subjects {
		list = append(list, &pb.SubjectData{
			Id:         subject.Id,
			Title:      subject.Name,
			AlbumTotal: uint32(subject.AlbumCount),
		})
	}
	response.Output(r, &pb.RespAdminVoiceLoverSubjectList{
		Success: true,
		List:    list,
		Total:   reply.Total,
	})
}

// AlbumCollect
// @Tags VoiceLoverAdmin
// @Summary 声恋后台audio 列表
// @Description 声恋后台audio 列表
// @Accept application/json
// @Produce json
// @Security ApiKeyAuth,OAuth2Implicit
// @Param request body  query.ReqAdminVoiceLoverAlbumCollect false "request"
// @Success 200 {object} pb.RespAdminVoiceLoverAlbumCollect
// @Router /go/func/admin/voice_lover/albumCollect [get]
func (a *voiceLoverAdminApi) AlbumCollect(r *ghttp.Request) {
	var req *query.ReqAdminVoiceLoverAlbumCollect
	if err := r.Parse(&req); err != nil {
		response.Output(r, &pb.RespAdminVoiceLoverAlbumCollect{
			Msg: err.Error(),
		})
		return
	}
	ctx := r.Context()
	_, err := voice_lover2.VoiceLoverAdmin.AlbumCollect(ctx, &voice_lover3.ReqAlbumCollect{
		AlbumId:   req.AlbumId,
		SubjectId: req.SubjectId,
	})
	if err != nil {
		response.Output(r, &pb.RespAdminVoiceLoverAlbumCollect{
			Msg: err.Error(),
		})
		return
	}
	response.Output(r, &pb.RespAdminVoiceLoverAlbumCollect{
		Success: true,
	})
}

// AlbumCollectList
// @Tags VoiceLoverAdmin
// @Summary 声恋后台audio 列表
// @Description 声恋后台audio 列表
// @Accept application/json
// @Produce json
// @Security ApiKeyAuth,OAuth2Implicit
// @Param request body query.ReqAdminVoiceLoverAlbumCollectList false "request"
// @Success 200 {object} pb.RespAdminVoiceLoverAlbumCollectList
// @Router /go/func/admin/voice_lover/albumCollectList [get]
func (a *voiceLoverAdminApi) AlbumCollectList(r *ghttp.Request) {
	var req *query.ReqAdminVoiceLoverAlbumCollectList
	if err := r.Parse(&req); err != nil {
		response.Output(r, &pb.RespAdminVoiceLoverAlbumCollectList{
			Msg: err.Error(),
		})
		return
	}
	ctx := r.Context()
	reply, err := voice_lover2.VoiceLoverAdmin.GetAlbumCollect(ctx, &voice_lover3.ReqGetAlbumCollect{
		AlbumStr:   req.AlbumStr,
		SubjectStr: req.SubjectStr,
	})
	if err != nil {
		response.Output(r, &pb.RespAdminVoiceLoverAlbumCollectList{
			Msg: err.Error(),
		})
		return
	}
	list := make([]*pb.AdminVoiceLoverAlbumSubject, 0)
	for _, a := range reply.AlbumCollects {
		list = append(list, &pb.AdminVoiceLoverAlbumSubject{
			Id:          a.Id,
			AlbumName:   a.AlbumName,
			SubjectName: a.SubjectName,
			AlbumId:     a.AlbumId,
			SubjectId:   a.SubjectId,
		})
	}
	response.Output(r, &pb.RespAdminVoiceLoverAlbumCollectList{
		Success: true,
		List:    list,
		Total:   reply.Total,
	})
}

// SubjectDetail
// @Tags VoiceLoverAdmin
// @Summary 声恋后台audio 列表
// @Description 声恋后台audio 列表
// @Accept application/json
// @Produce json
// @Security ApiKeyAuth,OAuth2Implicit
// @Param request body query.ReqAdminVoiceLoverSubjectDetail false "request"
// @Success 200 {object} pb.RespAdminVoiceLoverSubjectDetail
// @Router /go/func/admin/voice_lover/subjectDetail [get]
func (a *voiceLoverAdminApi) SubjectDetail(r *ghttp.Request) {
	var req *query.ReqAdminVoiceLoverSubjectDetail
	if err := r.Parse(&req); err != nil {
		response.Output(r, &pb.RespAdminVoiceLoverSubjectDetail{
			Msg: err.Error(),
		})
		return
	}
	ctx := r.Context()
	reply, err := voice_lover2.VoiceLoverAdmin.GetSubjectDetail(ctx, &voice_lover3.ReqGetSubjectDetail{
		Ids: []uint64{req.Id},
	})
	if err != nil {
		response.Output(r, &pb.RespAdminVoiceLoverSubjectDetail{
			Msg: err.Error(),
		})
		return
	}
	if reply.Subjects[req.Id] == nil {
		response.Output(r, &pb.RespAdminVoiceLoverSubjectDetail{
			Success: true,
		})
		return
	}

	response.Output(r, &pb.RespAdminVoiceLoverSubjectDetail{
		Success: true,
		Subject: &pb.SubjectData{
			Id:         reply.Subjects[req.Id].GetId(),
			Title:      reply.Subjects[req.Id].GetName(),
			AlbumTotal: uint32(reply.Subjects[req.Id].GetAlbumCount()),
		},
	})
}

// AlbumChoice
// @Tags VoiceLoverAdmin
// @Summary 声恋后台audio 列表
// @Description 声恋后台audio 列表
// @Accept application/json
// @Produce json
// @Security ApiKeyAuth,OAuth2Implicit
// @Param request body query.ReqAdminVoiceLoverAlbumChoice false "request"
// @Success 200 {object} pb.RespAdminVoiceLoverAlbumChoice
// @Router /go/func/admin/voice_lover/albumChoice [post]
func (a *voiceLoverAdminApi) AlbumChoice(r *ghttp.Request) {
	var req *query.ReqAdminVoiceLoverAlbumChoice
	if err := r.Parse(&req); err != nil {
		response.Output(r, &pb.RespAdminVoiceLoverAlbumChoice{
			Msg: err.Error(),
		})
		return
	}
	ctx := r.Context()
	_, err := voice_lover2.VoiceLoverAdmin.AlbumChoice(ctx, &voice_lover3.ReqAlbumChoice{
		Id:   req.Id,
		Type: req.Choice,
	})
	if err != nil {
		response.Output(r, &pb.RespAdminVoiceLoverAlbumChoice{
			Msg: err.Error(),
		})
		return
	}
	response.Output(r, &pb.RespAdminVoiceLoverAlbumChoice{
		Success: true,
	})
}

// AlbumChoiceList
// @Tags VoiceLoverAdmin
// @Summary 声恋后台audio 列表
// @Description 声恋后台audio 列表
// @Accept application/json
// @Produce json
// @Security ApiKeyAuth,OAuth2Implicit
// @Success 200 {object} pb.RespAdminVoiceLoverAlbumChoiceList
// @Router /go/func/admin/voice_lover/albumChoiceList [get]
func (a *voiceLoverAdminApi) AlbumChoiceList(r *ghttp.Request) {
	ctx := r.Context()
	reply, err := voice_lover2.VoiceLoverAdmin.GetAlbumChoice(ctx, &voice_lover3.ReqGetAlbumChoice{})
	if err != nil {
		response.Output(r, &pb.RespAdminVoiceLoverAlbumChoiceList{
			Msg: err.Error(),
		})
		return
	}
	res := make([]*pb.AdminAlbumData, 0)
	level := 1
	for _, r := range reply.Albums {
		res = append(res, &pb.AdminAlbumData{
			Id:         r.Id,
			Name:       r.Name,
			CreateTime: r.CreateTime,
			Level:      int32(level),
		})
		level = level + 1
	}
	response.Output(r, &pb.RespAdminVoiceLoverAlbumChoiceList{
		Success: true,
		Albums:  res,
	})
}

// BannerList
// @Tags VoiceLoverAdmin
// @Summary 声恋后台audio 列表
// @Description 声恋后台audio 列表
// @Accept application/json
// @Produce json
// @Security ApiKeyAuth,OAuth2Implicit
// @Param request body query.ReqAdminVoiceLoverBannerList false "request"
// @Success 200 {object} pb.RespAdminVoiceLoverBannerList
// @Router /go/func/admin/voice_lover/bannerList [get]
func (a *voiceLoverAdminApi) BannerList(r *ghttp.Request) {
	var req *query.ReqAdminVoiceLoverBannerList
	if err := r.Parse(&req); err != nil {
		response.Output(r, &pb.RespAdminVoiceLoverBannerList{
			Msg: err.Error(),
		})
		return
	}
	ctx := r.Context()
	reply, err := voice_lover2.VoiceLoverAdmin.GetBannerList(ctx, &voice_lover3.ReqGetBannerList{
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		Status:    req.Status,
		Title:     req.Title,
	})
	if err != nil {
		response.Output(r, &pb.RespAdminVoiceLoverBannerList{
			Msg: err.Error(),
		})
		return
	}
	res := make([]*pb.AdminBannerData, 0)
	for _, b := range reply.Banners {
		res = append(res, &pb.AdminBannerData{
			Id:         b.Id,
			Title:      b.Title,
			Cover:      b.Cover,
			Schema:     b.Schema,
			OpUid:      b.OpUid,
			Sort:       b.Sort,
			StartTime:  b.StartTime,
			EndTime:    b.EndTime,
			CreateTime: b.CreateTime,
		})
	}
	response.Output(r, &pb.RespAdminVoiceLoverBannerList{
		Success: true,
		Total:   reply.Total,
		List:    res,
	})
}

// BannerCreate
// @Tags VoiceLoverAdmin
// @Summary 声恋后台audio 列表
// @Description 声恋后台audio 列表
// @Accept application/json
// @Produce json
// @Security ApiKeyAuth,OAuth2Implicit
// @Param request body query.ReqAdminVoiceLoverBannerCreate false "request"
// @Success 200 {object} pb.RespAdminVoiceLoverBannerCreate
// @Router /go/func/admin/voice_lover/bannerCreate [post]
func (a *voiceLoverAdminApi) BannerCreate(r *ghttp.Request) {
	var req *query.ReqAdminVoiceLoverBannerCreate
	if err := r.Parse(&req); err != nil {
		response.Output(r, &pb.RespAdminVoiceLoverBannerCreate{
			Msg: err.Error(),
		})
		return
	}
	ctx := r.Context()
	reply, err := voice_lover2.VoiceLoverAdmin.CreateBanner(ctx, &voice_lover3.ReqCreateBanner{
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		Title:     req.Title,
		OpUid:     req.OpUid,
		Sort:      req.Sort,
		Cover:     req.Cover,
		Schema:    req.Schema,
	})
	if err != nil {
		response.Output(r, &pb.RespAdminVoiceLoverBannerCreate{
			Msg: err.Error(),
		})
		return
	}
	response.Output(r, &pb.RespAdminVoiceLoverBannerCreate{
		Success: true,
		Id:      reply.Id,
	})
}

// BannerUpdate
// @Tags VoiceLoverAdmin
// @Summary 声恋后台audio 列表
// @Description 声恋后台audio 列表
// @Accept application/json
// @Produce json
// @Security ApiKeyAuth,OAuth2Implicit
// @Param request body query.ReqAdminVoiceLoverBannerUpdate false "request"
// @Success 200 {object} pb.RespAdminVoiceLoverBannerUpdate
// @Router /go/func/admin/voice_lover/bannerUpdate [post]
func (a *voiceLoverAdminApi) BannerUpdate(r *ghttp.Request) {
	var req *query.ReqAdminVoiceLoverBannerUpdate
	if err := r.Parse(&req); err != nil {
		response.Output(r, &pb.RespAdminVoiceLoverBannerUpdate{
			Msg: err.Error(),
		})
		return
	}
	ctx := r.Context()
	_, err := voice_lover2.VoiceLoverAdmin.UpdateBanner(ctx, &voice_lover3.ReqUpdateBanner{
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		Title:     req.Title,
		OpUid:     req.OpUid,
		Sort:      req.Sort,
		Cover:     req.Cover,
		Schema:    req.Schema,
		Id:        req.Id,
	})
	if err != nil {
		response.Output(r, &pb.RespAdminVoiceLoverBannerUpdate{
			Msg: err.Error(),
		})
		return
	}
	response.Output(r, &pb.RespAdminVoiceLoverBannerUpdate{
		Success: true,
	})
}

// BannerDetail
// @Tags VoiceLoverAdmin
// @Summary 声恋后台audio 列表
// @Description 声恋后台audio 列表
// @Accept application/json
// @Produce json
// @Security ApiKeyAuth,OAuth2Implicit
// @Param request body query.ReqAdminVoiceLoverBannerDetail false "request"
// @Success 200 {object} pb.RespAdminVoiceLoverBannerDetail
// @Router /go/func/admin/voice_lover/bannerDetail [get]
func (a *voiceLoverAdminApi) BannerDetail(r *ghttp.Request) {
	var req *query.ReqAdminVoiceLoverBannerDetail
	if err := r.Parse(&req); err != nil {
		response.Output(r, &pb.RespAdminVoiceLoverBannerDetail{
			Msg: err.Error(),
		})
		return
	}
	ctx := r.Context()
	reply, err := voice_lover2.VoiceLoverAdmin.GetBannerDetail(ctx, &voice_lover3.ReqGetBannerDetail{

		Id: req.Id,
	})
	if err != nil {
		response.Output(r, &pb.RespAdminVoiceLoverBannerDetail{
			Msg: err.Error(),
		})
		return
	}
	if reply.Banner == nil {
		response.Output(r, &pb.RespAdminVoiceLoverBannerDetail{
			Success: true,
		})
	}
	banner := &pb.AdminBannerData{
		Id:         reply.Banner.Id,
		Title:      reply.Banner.Title,
		Cover:      reply.Banner.Cover,
		Schema:     reply.Banner.Schema,
		OpUid:      reply.Banner.OpUid,
		Sort:       reply.Banner.Sort,
		StartTime:  reply.Banner.StartTime,
		EndTime:    reply.Banner.EndTime,
		CreateTime: reply.Banner.CreateTime,
	}
	response.Output(r, &pb.RespAdminVoiceLoverBannerDetail{
		Success: true,
		Banner:  banner,
	})
}
