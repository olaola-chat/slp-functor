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

func (v *VoiceLoverMain) GetRecAlbums(ctx context.Context, req *vl_pb.ReqGetRecAlbums, reply *vl_pb.ResGetRecAlbums) error {
	return logic.MainLogic.GetRecAlbums(ctx, req, reply)
}

func (v *VoiceLoverMain) BatchGetAlbumAudioCount(ctx context.Context, req *vl_pb.ReqBatchGetAlbumAudioCount, reply *vl_pb.ResBatchGetAlbumAudioCount) error {
	return logic.MainLogic.BatchGetAlbumAudioCount(ctx, req, reply)
}
