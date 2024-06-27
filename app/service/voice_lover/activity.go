package voice_lover

import (
	"context"

	"github.com/olaola-chat/rbp-functor/app/pb"
)

var ActivitySrv = &activitySrv{}

type activitySrv struct{}

func (a *activitySrv) GetInfo(ctx context.Context, id uint32) (*pb.VoiceLoverActivityMain, error) {
	return nil, nil
}

func (a *activitySrv) getInfo(ctx context.Context, id uint32) (*pb.VoiceLoverActivity, *pb.VoiceLoverActivityAward, error) {
	//vl_rpc.VoiceLoverAdmin
	return nil, nil, nil
}
