package voice_lover

import (
	"context"

	"github.com/gogf/gf/frame/g"
	vl_pb "github.com/olaola-chat/rbp-proto/gen_pb/rpc/voice_lover"

	"github.com/olaola-chat/rbp-functor/rpc/voice_lover/internal/logic"
)

func NewVoiceLoverAdmin() interface{} {
	return &VoiceLoverAdmin{}
}

type VoiceLoverAdmin struct {
}

func (v *VoiceLoverAdmin) GetAudioDetail(ctx context.Context, req *vl_pb.ReqGetAudioDetail, reply *vl_pb.ResGetAudioDetail) error {
	data, err := logic.AdminLogic.GetAudioDetail(ctx, req.GetId())
	if err != nil {
		return err
	}
	reply.Audio = data
	g.Log().Infof("GetAudioDetail data = %v", data)
	return nil
}
