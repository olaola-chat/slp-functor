package api

import (
	"sort"

	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
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
		OutputCustomError(r, consts.ERROR_PARAM)
		return
	}
	ctx := r.Context()
	g.Log().Infof("VoiceLoverAdmin audio_detail param = %v", *req)
	data := &pb.RespAdminVoiceLoverAudioDetail{}
	reply, err := voice_lover2.VoiceLoverAdmin.GetAudioDetail(ctx, &voice_lover3.ReqGetAudioDetail{Id: req.Id})
	if err != nil {
		OutputCustomError(r, consts.ERROR_SYSTEM)
		return
	}
	if reply.Audio == nil {
		OutputCustomData(r, data)
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
	data.Audio = &pb.AdminVoiceLoverAudio{
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
	OutputCustomData(r, data)
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
// @Success 200 {object}
// @Router /go/func/admin/voice_lover/audio-update [post]
func (a *voiceLoverAdminApi) AudioUpdate(r *ghttp.Request) {
	var req *query.ReqAdminVoiceLoverAudioUpdate
	if err := r.Parse(&req); err != nil {
		OutputCustomError(r, consts.ERROR_PARAM)
		return
	}
	ctx := r.Context()
	_, err := voice_lover2.VoiceLoverAdmin.UpdateAudio(ctx, &voice_lover3.ReqUpdateAudio{Id: req.Id, Title: req.Title, Desc: req.Desc, Labels: req.Labels})
	if err != nil {
		OutputCustomError(r, consts.ERROR_SYSTEM)
		return
	}
	OutputCustomData(r, nil)
}

// AudioAudit
// @Tags VoiceLoverAdmin
// @Summary 声恋后台audio 列表
// @Description 声恋后台audio 列表
// @Accept application/json
// @Produce json
// @Security ApiKeyAuth,OAuth2Implicit
// @Request query.ReqAdminVoiceLoverAudioAudit query
// @Success 200 {object}
// @Router /go/func/admin/voice_lover/audio-audit [post]
func (a *voiceLoverAdminApi) AudioAudit(r *ghttp.Request) {
	var req *query.ReqAdminVoiceLoverAudioAudit
	if err := r.Parse(&req); err != nil {
		OutputCustomError(r, consts.ERROR_PARAM)
		return
	}
	ctx := r.Context()
	_, err := voice_lover2.VoiceLoverAdmin.AuditAudio(ctx, &voice_lover3.ReqAuditAudio{Id: req.Id, AuditStatus: req.AuditStatus, AuditReason: req.AuditReason, OpUid: req.OpUid})
	if err != nil {
		OutputCustomError(r, consts.ERROR_SYSTEM)
		return
	}
	OutputCustomData(r, nil)
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
	OutputCustomData(r, data)
}
