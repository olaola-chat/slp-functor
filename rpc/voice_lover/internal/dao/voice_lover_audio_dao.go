package dao

import (
	"context"

	"github.com/gogf/gf/frame/g"
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

func (v *voiceLoverAudioDao) UpdateAudioById(ctx context.Context, id uint64, data g.Map) (int64, error) {
	sqlRes, err := functor2.VoiceLoverAudio.Ctx(ctx).Where("id", id).Update(data)
	if err != nil {
		return 0, err
	}
	affect, _ := sqlRes.RowsAffected()
	return affect, nil
}
