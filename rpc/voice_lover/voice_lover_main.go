package voice_lover

import (
	"context"
	"github.com/olaola-chat/rbp-functor/rpc/voice_lover/internal/logic"
	"github.com/olaola-chat/rbp-proto/gen_pb/rpc/voice_lover"
)

func NewVoiceLoverMain() interface{} {
	return &VoiceLoverMain{}
}

type VoiceLoverMain struct {
}

func (v *VoiceLoverMain) Post(ctx context.Context, req *voice_lover.ReqVoiceLoverPost, reply *voice_lover.ResVoiceLoverBase) error {
	return logic.MainLogic.Post(ctx, req, reply)
}
