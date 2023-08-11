package voice_lover

import (
	"context"

	vl_pb "github.com/olaola-chat/rbp-proto/gen_pb/rpc/voice_lover"

	"github.com/olaola-chat/rbp-functor/rpc/voice_lover/internal/logic"
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

func (v *VoiceLoverMain) GetAudioListByAlbumId(ctx context.Context, req *vl_pb.ReqGetAudioListByAlbumId, reply *vl_pb.ResGetAudioListByAlbumId) error {
	return logic.MainLogic.GetAudioListByAlbumId(ctx, req, reply)
}
