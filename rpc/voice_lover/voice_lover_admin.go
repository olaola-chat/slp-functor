package voice_lover

import (
	"context"
	"github.com/olaola-chat/rbp-functor/rpc/voice_lover/internal/logic"
	"github.com/olaola-chat/rbp-proto/gen_pb/rpc/voice_lover"
)

func NewVoiceLoverAdmin() interface{} {
	return &VoiceLoverAdmin{}
}

type VoiceLoverAdmin struct {
}

func (v *VoiceLoverAdmin) Post(ctx context.Context, req *voice_lover.ReqVoiceLoverPost, reply *voice_lover.ResVoiceLoverBase) error {
	return logic.MainLogic.Post(ctx, req, reply)
}
