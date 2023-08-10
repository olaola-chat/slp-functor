package dao

import (
	"context"
	"fmt"

	"github.com/gogf/gf/frame/g"
	"github.com/olaola-chat/rbp-proto/dao/functor"
	functor2 "github.com/olaola-chat/rbp-proto/gen_pb/db/functor"
)

type voiceLoverAlbumSubjectDao struct {
}

var VoiceLoverAlbumSubjectDao = &voiceLoverAlbumSubjectDao{}

func (v *voiceLoverAlbumSubjectDao) GetListBySubjectIds(ctx context.Context, subjectIds []uint64) ([]*functor2.EntityVoiceLoverAlbumSubject, error) {
	list, err := functor.VoiceLoverAlbumSubject.Ctx(ctx).
		Where("status IN (?)", g.Slice{1, 2, 3}).
		Where(fmt.Sprintf("%s IN (?)", functor.VoiceLoverAlbumSubject.Columns.SubjectID), subjectIds).
		Order(functor.VoiceLoverAlbumSubject.Columns.CreateTime, "desc").FindAll()
	if err != nil {
		return nil, err
	}
	return list, nil
}
