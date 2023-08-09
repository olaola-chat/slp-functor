package dao

import (
	"context"

	vl_pb "github.com/olaola-chat/rbp-proto/gen_pb/rpc/voice_lover"
)

type voiceLoverAudioDao struct {
}

const (
	None = iota
	Dub
	Content
	Post
	Cover
)

var VoiceLoverAudioDao = &voiceLoverAudioDao{}

func (v *voiceLoverAudioDao) Post(ctx context.Context, req *vl_pb.ReqVoiceLoverPost) error {
	return nil
}
