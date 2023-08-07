package logic

import (
	"context"
	"github.com/gogf/gf/frame/g"
	"github.com/olaola-chat/rbp-proto/gen_pb/rpc/voice_lover"
)

type mainLogic struct {
}

var MainLogic = &mainLogic{}

func (m *mainLogic) Post(ctx context.Context, req *voice_lover.ReqVoiceLoverPost, reply *voice_lover.ResVoiceLoverBase) error {
	g.Log().Infof("VoiceLoverPost req = %v", req)
	return nil
}
