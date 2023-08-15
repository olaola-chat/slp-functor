package dao

import (
	"context"

	"github.com/gogf/gf/frame/g"

	"github.com/olaola-chat/rbp-proto/dao/functor"
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

func (v *voiceLoverAlbumCommentDao) Insert(ctx context.Context, data g.Map) error {
	return nil
}
