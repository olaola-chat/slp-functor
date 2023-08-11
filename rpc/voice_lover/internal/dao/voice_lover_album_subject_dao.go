package dao

import (
	"context"

	functor2 "github.com/olaola-chat/rbp-proto/dao/functor"
)

type voiceLoverAlbumSubjectDao struct {
}

var VoiceLoverAlbumSubjectDao = &voiceLoverAlbumSubjectDao{}

func (v *voiceLoverAlbumSubjectDao) GetCountBySubjectId(ctx context.Context, id uint64) (int, error) {
	count, err := functor2.VoiceLoverAlbumSubject.Ctx(ctx).Where("subject_id", id).Count()
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (v *voiceLoverAlbumSubjectDao) GetCountByAlbumId(ctx context.Context, id uint64) (int, error) {
	count, err := functor2.VoiceLoverAlbumSubject.Ctx(ctx).Where("album_id", id).Count()
	if err != nil {
		return 0, err
	}
	return count, nil
}
