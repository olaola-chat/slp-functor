package dao

import (
	"context"

	functor2 "github.com/olaola-chat/rbp-proto/dao/functor"
	"github.com/olaola-chat/rbp-proto/gen_pb/db/functor"
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

func (v *voiceLoverAudioDao) GetAudioDetailByAudioId(ctx context.Context, id uint64) (*functor.EntityVoiceLoverAudio, error) {
	res, err := functor2.VoiceLoverAudio.Ctx(ctx).Where("id", id).One()
	if err != nil {
		return nil, err
	}
	return res, nil
}
