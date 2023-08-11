package dao

import (
	"context"
	"fmt"

	"github.com/olaola-chat/rbp-proto/dao/functor"
	functor2 "github.com/olaola-chat/rbp-proto/gen_pb/db/functor"
)

type voiceLoverAlbumSubjectDao struct {
}

var VoiceLoverAlbumSubjectDao = &voiceLoverAlbumSubjectDao{}

func (v *voiceLoverAlbumSubjectDao) GetListBySubjectIds(ctx context.Context, subjectIds []uint64) ([]*functor2.EntityVoiceLoverAlbumSubject, error) {
	list, err := functor.VoiceLoverAlbumSubject.Ctx(ctx).
		Where(fmt.Sprintf("%s IN (?)", functor.VoiceLoverAlbumSubject.Columns.SubjectID), subjectIds).
		Order(functor.VoiceLoverAlbumSubject.Columns.CreateTime, "desc").FindAll()
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (v *voiceLoverAlbumSubjectDao) GetListBySubjectId(ctx context.Context, subjectId uint64, page int, limit int) ([]*functor2.EntityVoiceLoverAlbumSubject, error) {
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * limit
	list, err := functor.VoiceLoverAlbumSubject.Ctx(ctx).
		Where(functor.VoiceLoverAlbumSubject.Columns.SubjectID, subjectId).
		Order(functor.VoiceLoverAlbum.Columns.CreateTime, "desc").
		Offset(offset).
		Limit(limit).FindAll()
	if err != nil {
		return nil, err
	}
	return list, nil
}
