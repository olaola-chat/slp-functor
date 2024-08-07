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

const (
	AudioCommentStatusWait     = 1 // 待审核
	AudioCommentStatusPass     = 2 // 审核通过
	AudioCommentStatusRejected = 3 // 审核拒绝
)

func (v *voiceLoverAudioCommentDao) GetList(ctx context.Context, audioId uint64, offset int32, limit uint32) (
	[]*functor.EntityVoiceLoverAudioComment, error) {
	res, err := functor2.VoiceLoverAudioComment.Ctx(ctx).
		Where(functor2.VoiceLoverAudioComment.Columns.AudioID, audioId).
		Where("audit_status IN (?)", g.Slice{0, AudioCommentStatusPass}).
		Offset(int(offset)).
		Limit(int(limit)).
		Order(functor2.VoiceLoverAudioComment.Columns.CreateTime, "desc").
		FindAll()
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (v *voiceLoverAudioCommentDao) Insert(ctx context.Context, data g.Map) (int64, error) {
	sqlRes, err := functor2.VoiceLoverAudioComment.Ctx(ctx).Insert(data)
	if err != nil {
		return 0, err
	}
	return sqlRes.LastInsertId()
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

func (v *voiceLoverAudioCommentDao) UpdateAuditStatus(ctx context.Context, id uint64, status int) error {
	_, err := functor2.VoiceLoverAudioComment.Ctx(ctx).Where("id=?", id).Data(
		g.Map{
			functor2.VoiceLoverAudioComment.Columns.AuditStatus: status,
		},
	).Update()
	if err != nil {
		return err
	}
	return nil
}
