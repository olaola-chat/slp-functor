package dao

import (
	"context"

	"github.com/gogf/gf/frame/g"

	functor2 "github.com/olaola-chat/rbp-proto/dao/functor"
	"github.com/olaola-chat/rbp-proto/gen_pb/db/functor"
)

type voiceLoverAudioCommentDao struct {
}

var VoiceLoverAudioCommentDao = &voiceLoverAudioCommentDao{}

func (v *voiceLoverAudioCommentDao) GetList(ctx context.Context, id uint64) ([]*functor.EntityVoiceLoverAudioComment, error) {
	res, err := functor2.VoiceLoverAudioComment.Ctx(ctx).Where("id=?", id).Limit(20).FindAll()
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (v *voiceLoverAudioCommentDao) Insert(ctx context.Context, data g.Map) (bool, error) {
	sqlRes, err := functor2.VoiceLoverAudioComment.Ctx(ctx).Insert(data)
	if err != nil {
		return false, err
	}
	affect, _ := sqlRes.RowsAffected()
	return affect > 0, nil
}
