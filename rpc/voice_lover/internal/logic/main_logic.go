package logic

import (
	"context"

	"github.com/gogf/gf/frame/g"
	vl_pb "github.com/olaola-chat/rbp-proto/gen_pb/rpc/voice_lover"

	"github.com/olaola-chat/rbp-functor/rpc/voice_lover/internal/dao"
)

type mainLogic struct {
}

var MainLogic = &mainLogic{}

func (m *mainLogic) Post(ctx context.Context, req *vl_pb.ReqVoiceLoverPost, reply *vl_pb.ResVoiceLoverBase) error {
	g.Log().Infof("VoiceLoverPost req = %v", req)
	return dao.VoiceLoverAudioDao.Post(ctx, req)
}

func (m *mainLogic) GetRecAlbums(ctx context.Context, req *vl_pb.ReqGetRecAlbums, reply *vl_pb.ResGetRecAlbums) error {
	return nil
}
