package dao

import (
	"context"

	"github.com/gogf/gf/frame/g"

	"github.com/olaola-chat/slp-proto/dao/functor"
	dbpb "github.com/olaola-chat/slp-proto/gen_pb/db/functor"
)

type voiceLoverAlbumCommentDao struct {
}

const (
	AlbumComment_StatusOk = 0
)

const (
	AlbumCommentStatusWait     = 1 // 待审核
	AlbumCommentStatusPass     = 2 // 审核通过
	AlbumCommentStatusRejected = 3 // 审核拒绝
)

var VoiceLoverAlbumCommentDao = &voiceLoverAlbumCommentDao{}

func (v *voiceLoverAlbumCommentDao) GetValidCommentCountByAlbumId(ctx context.Context, albumId uint64) (int, error) {
	count, err := functor.VoiceLoverAlbumComment.Ctx(ctx).
		Where(functor.VoiceLoverAlbumComment.Columns.AlbumID, albumId).
		Where(functor.VoiceLoverAlbumComment.Columns.AuditStatus, AlbumCommentStatusPass).
		Where(functor.VoiceLoverAlbumComment.Columns.Status, AlbumComment_StatusOk).Count()
	if err != nil {
		g.Log().Errorf("voiceLoverAlbumCommentDao GetValidCommentCountByAlbumId error=%v", err)
		return 0, err
	}
	return count, nil
}

func (v *voiceLoverAlbumCommentDao) GetList(ctx context.Context, albumId uint64, offset int32, limit uint32) ([]*dbpb.EntityVoiceLoverAlbumComment, error) {
	res, err := functor.VoiceLoverAlbumComment.Ctx(ctx).
		Where(functor.VoiceLoverAlbumComment.Columns.AlbumID, albumId).
		Where(functor.VoiceLoverAlbumComment.Columns.AuditStatus, AlbumCommentStatusPass).
		Offset(int(offset)).
		Limit(int(limit)).
		Order(functor.VoiceLoverAlbumComment.Columns.CreateTime, "desc").FindAll()
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (v *voiceLoverAlbumCommentDao) Insert(ctx context.Context, data g.Map) (int64, error) {
	sqlRes, err := functor.VoiceLoverAlbumComment.Ctx(ctx).Insert(data)
	if err != nil {
		return 0, err
	}
	return sqlRes.LastInsertId()
}

func (v *voiceLoverAlbumCommentDao) UpdateStatus(ctx context.Context, id uint64, status uint32) (bool, error) {
	_, err := functor.VoiceLoverAlbumComment.Ctx(ctx).Where("id=?", id).Data(
		g.Map{
			functor.VoiceLoverAlbumComment.Columns.Status: status,
		},
	).Update()
	if err != nil {
		return false, err
	}
	return true, nil
}

func (v *voiceLoverAlbumCommentDao) UpdateAuditStatus(ctx context.Context, id uint64, status uint32) error {
	_, err := functor.VoiceLoverAlbumComment.Ctx(ctx).Where("id=?", id).Data(
		g.Map{
			functor.VoiceLoverAlbumComment.Columns.AuditStatus: status,
		},
	).Update()
	if err != nil {
		return err
	}
	return nil
}
