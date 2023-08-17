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

func (v *voiceLoverAudioCommentDao) GetList(ctx context.Context, audioId uint64, offset int32, limit uint32) (
	[]*functor.EntityVoiceLoverAudioComment, error) {

	if offset <= 1 {
		offset = 1
	}
	if limit <= 0 {
		limit = 10
	}
	res, err := functor2.VoiceLoverAudioComment.Ctx(ctx).Where(
		functor2.VoiceLoverAudioComment.Columns.AudioID, audioId).Offset(int(offset)).Limit(int(limit)).FindAll()
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

func (v *voiceLoverAudioCommentDao) UpdateStatus(ctx context.Context, id uint64, status uint32) (bool, error) {
	_, err := functor2.VoiceLoverAudioComment.Ctx(ctx).Where("id=?", id).Data(
		g.Map{
			functor2.VoiceLoverAudioComment.Columns.Status: status,
		},
	).Update()
	if err != nil {
		return false, err
	}
	return true, nil
}
