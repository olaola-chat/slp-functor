package voice_lover

import (
	"context"

	vl_pb "github.com/olaola-chat/slp-proto/gen_pb/rpc/voice_lover"

	"github.com/olaola-chat/slp-functor/rpc/voice_lover/internal/logic"
)

func NewVoiceLoverMain() interface{} {
	return &VoiceLoverMain{}
}

type VoiceLoverMain struct {
}

func (v *VoiceLoverMain) Post(ctx context.Context, req *vl_pb.ReqPost, reply *vl_pb.ResBase) error {
	return logic.MainLogic.Post(ctx, req, reply)
}

func (v *VoiceLoverMain) GetAlbumInfoById(ctx context.Context, req *vl_pb.ReqGetAlbumInfoById, reply *vl_pb.ResGetAlbumInfoById) error {
	return logic.MainLogic.GetAlbumInfoById(ctx, req, reply)
}

func (v *VoiceLoverMain) GetAlbumCommentCount(ctx context.Context, req *vl_pb.ReqGetAlbumCommentCount, reply *vl_pb.ResGetAlbumCommentCount) error {
	return logic.MainLogic.GetAlbumCommentCount(ctx, req, reply)
}

func (v *VoiceLoverMain) GetRecAlbums(ctx context.Context, req *vl_pb.ReqGetRecAlbums, reply *vl_pb.ResGetRecAlbums) error {
	return logic.MainLogic.GetRecAlbums(ctx, req, reply)
}

func (v *VoiceLoverMain) GetRecBanners(ctx context.Context, req *vl_pb.ReqGetRecBanners, reply *vl_pb.ResGetRecBanners) error {
	return logic.MainLogic.GetRecBanners(ctx, req, reply)
}

func (v *VoiceLoverMain) GetRecCommonAlbums(ctx context.Context, req *vl_pb.ReqGetRecCommonAlbums, reply *vl_pb.ResGetRecAlbums) error {
	return logic.MainLogic.GetRecCommonAlbums(ctx, req, reply)
}

func (v *VoiceLoverMain) GetAlbumsByPage(ctx context.Context, req *vl_pb.ReqGetAlbumsByPage, reply *vl_pb.ResGetAlbumsByPage) error {
	return logic.MainLogic.GetAlbumsByPage(ctx, req, reply)
}

func (v *VoiceLoverMain) GetSubjectAlbumsByPage(ctx context.Context, req *vl_pb.ReqGetSubjectAlbumsByPage, reply *vl_pb.ResGetAlbumsByPage) error {
	return logic.MainLogic.GetSubjectAlbumsByPage(ctx, req, reply)
}

func (v *VoiceLoverMain) GetRecSubjects(ctx context.Context, req *vl_pb.ReqGetRecSubjects, reply *vl_pb.ResGetRecSubjects) error {
	return logic.MainLogic.GetRecSubjects(ctx, req, reply)
}

func (v *VoiceLoverMain) BatchGetAlbumAudioCount(ctx context.Context, req *vl_pb.ReqBatchGetAlbumAudioCount, reply *vl_pb.ResBatchGetAlbumAudioCount) error {
	return logic.MainLogic.BatchGetAlbumAudioCount(ctx, req, reply)
}

func (v *VoiceLoverMain) IsUserCollectAlbum(ctx context.Context, req *vl_pb.ReqIsUserCollectAlbum, reply *vl_pb.ResIsUserCollectAlbum) error {
	return logic.MainLogic.IsUserCollectAlbum(ctx, req, reply)
}

func (v *VoiceLoverMain) IsUserCollectAlbums(ctx context.Context, req *vl_pb.ReqIsUserCollectAlbums, reply *vl_pb.ResIsUserCollectAlbums) error {
	return logic.MainLogic.IsUserCollectAlbums(ctx, req, reply)
}

func (v *VoiceLoverMain) Collect(ctx context.Context, req *vl_pb.ReqCollect, reply *vl_pb.ResCollect) error {
	return logic.MainLogic.Collect(ctx, req, reply)
}

func (v *VoiceLoverMain) GetAlbumCollectList(ctx context.Context, req *vl_pb.ReqGetAlbumCollectList, reply *vl_pb.ResGetAlbumCollectList) error {
	return logic.MainLogic.GetAlbumCollectList(ctx, req, reply)
}

func (v *VoiceLoverMain) GetAudioCollectList(ctx context.Context, req *vl_pb.ReqGetAudioCollectList, reply *vl_pb.ResGetAudioCollectList) error {
	return logic.MainLogic.GetAudioCollectList(ctx, req, reply)
}

func (v *VoiceLoverMain) GetAudioListByAlbumId(ctx context.Context, req *vl_pb.ReqGetAudioListByAlbumId, reply *vl_pb.ResGetAudioListByAlbumId) error {
	return logic.MainLogic.GetAudioListByAlbumId(ctx, req, reply)
}

func (v *VoiceLoverMain) SubmitAudioComment(ctx context.Context, req *vl_pb.ReqAudioSubmitComment, reply *vl_pb.ResCommonPost) error {
	return logic.MainLogic.SubmitAudioComment(ctx, req, reply)
}

func (v *VoiceLoverMain) GetAudioCommentList(ctx context.Context, req *vl_pb.ReqGetAudioCommentList, reply *vl_pb.ResCommentList) error {
	return logic.MainLogic.GetAudioCommentList(ctx, req, reply)
}

func (v *VoiceLoverMain) SubmitAlbumComment(ctx context.Context, req *vl_pb.ReqAlbumSubmitComment, reply *vl_pb.ResCommonPost) error {
	return logic.MainLogic.SubmitAlbumComment(ctx, req, reply)
}

func (v *VoiceLoverMain) GetAlbumCommentList(ctx context.Context, req *vl_pb.ReqGetAlbumCommentList, reply *vl_pb.ResCommentList) error {
	return logic.MainLogic.GetAlbumCommentList(ctx, req, reply)
}

func (v *VoiceLoverMain) GetAudioInfoById(ctx context.Context, req *vl_pb.ReqGetAudioDetail, reply *vl_pb.ResGetAudioDetail) error {
	return logic.MainLogic.GetAudioInfoById(ctx, req, reply)
}

func (v *VoiceLoverMain) UpdateReportStatus(ctx context.Context, req *vl_pb.ReqUpdateStatus, reply *vl_pb.ResUpdateStatus) error {
	return logic.MainLogic.UpdateReportStatus(ctx, req, reply)
}

func (v *VoiceLoverMain) PlayStatReport(ctx context.Context, req *vl_pb.ReqPlayStatReport, reply *vl_pb.ResPlayStatReport) error {
	return logic.MainLogic.PlayStatReport(ctx, req, reply)
}

func (v *VoiceLoverMain) IsUserCollectAudio(ctx context.Context, req *vl_pb.ReqCollect, reply *vl_pb.ResIsUserCollectAudio) error {
	return logic.MainLogic.IsUserCollectAudio(ctx, req, reply)
}

func (v *VoiceLoverMain) GetValidAudioUsers(ctx context.Context, req *vl_pb.ReqGetValidAudioUsers, reply *vl_pb.ResGetValidAudioUsers) error {
	return logic.MainLogic.GetValidAudioUsers(ctx, req, reply)
}

func (v *VoiceLoverMain) BatchGetAudioInfo(ctx context.Context, req *vl_pb.ReqBatchGetAudioInfo, reply *vl_pb.RespBatchGetAudioInfo) error {
	return logic.MainLogic.BatchGetAudioInfo(ctx, req, reply)
}

func (v *VoiceLoverMain) GenRecAlbum(ctx context.Context, req *vl_pb.ReqGenRecAlbum, reply *vl_pb.RespGenRecAlbum) error {
	return logic.MainLogic.GenRecAlbum(ctx, req, reply)
}

func (v *VoiceLoverMain) BatchCheckUserCollect(ctx context.Context, req *vl_pb.ReqBatchCheckUserCollect, reply *vl_pb.RespBatchCheckUserCollect) error {
	return logic.MainLogic.BatchCheckUserCollect(ctx, req, reply)
}

func (v *VoiceLoverMain) BatchGetCollectNum(ctx context.Context, req *vl_pb.ReqBatchGetCollectNum, reply *vl_pb.RespBatchGetCollectNum) error {
	return logic.MainLogic.BatchGetCollectNum(ctx, req, reply)
}

func (v *VoiceLoverMain) AudioCommentAuditCallback(ctx context.Context, req *vl_pb.ReqAudioCommentAuditCallback, reply *vl_pb.RespAudioCommentAuditCallback) error {
	return logic.MainLogic.AudioCommentAuditCallback(ctx, req, reply)
}

func (v *VoiceLoverMain) AlbumCommentAuditCallback(ctx context.Context, req *vl_pb.ReqAlbumCommentAuditCallback, reply *vl_pb.RespAlbumCommentAuditCallback) error {
	return logic.MainLogic.AlbumCommentAuditCallback(ctx, req, reply)
}
