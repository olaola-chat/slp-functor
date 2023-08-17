package dao

import (
	"context"

	"github.com/gogf/gf/frame/g"

	"github.com/olaola-chat/rbp-proto/dao/functor"
	dbpb "github.com/olaola-chat/rbp-proto/gen_pb/db/functor"
)

type voiceLoverAlbumCommentDao struct {
}

const (
	AlbumComment_StatusOk = 0
)

var VoiceLoverAlbumCommentDao = &voiceLoverAlbumCommentDao{}

func (v *voiceLoverAlbumCommentDao) GetValidCommentCountByAlbumId(ctx context.Context, albumId uint64) (int, error) {
	count, err := functor.VoiceLoverAlbumComment.Ctx(ctx).
		Where(functor.VoiceLoverAlbumComment.Columns.AlbumID, albumId).
		Where(functor.VoiceLoverAlbumComment.Columns.Status, AlbumComment_StatusOk).Count()
	if err != nil {
		g.Log().Errorf("voiceLoverAlbumCommentDao GetValidCommentCountByAlbumId error=%v", err)
		return 0, err
	}
	return count, nil
}

func (v *voiceLoverAlbumCommentDao) GetList(ctx context.Context, albumId uint64, offset int32, limit uint32) ([]*dbpb.EntityVoiceLoverAlbumComment, error) {

	if offset <= 1 {
		offset = 1
	}
	if limit <= 0 {
		limit = 10
	}
	res, err := functor.VoiceLoverAlbumComment.Ctx(ctx).Where(functor.VoiceLoverAlbumComment.Columns.AlbumID,
		albumId).Offset(int(offset)).Limit(int(limit)).FindAll()
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (v *voiceLoverAlbumCommentDao) Insert(ctx context.Context, data g.Map) (bool, error) {
	sqlRes, err := functor.VoiceLoverAlbumComment.Ctx(ctx).Insert(data)
	if err != nil {
		return false, err
	}
	affect, _ := sqlRes.RowsAffected()
	return affect > 0, nil
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
