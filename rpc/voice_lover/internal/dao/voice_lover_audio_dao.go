package dao

import (
	"context"

	vl_pb "github.com/olaola-chat/rbp-proto/gen_pb/rpc/voice_lover"
)

type voiceLoverDao struct {
}

const (
	None = iota
	Dub
	Content
	Post
	Cover
)

var VoiceLoverDao = &voiceLoverDao{}

func (v *voiceLoverDao) Post(ctx context.Context, req *vl_pb.ReqVoiceLoverPost) error {
	return nil
}
