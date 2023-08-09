package voice_lover

import (
	"context"

	"github.com/olaola-chat/rbp-functor/app/pb"
)

var VoiceLoverService = &voiceLoverService{}

type voiceLoverService struct{}

func (serv *voiceLoverService) GetMainData(ctx context.Context, uid uint32) (*pb.RespVoiceLoverMain, error) {
	return nil, nil
}
