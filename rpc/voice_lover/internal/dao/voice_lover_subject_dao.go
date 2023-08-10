package dao

import (
	"context"

	"github.com/olaola-chat/rbp-proto/dao/functor"
	functor2 "github.com/olaola-chat/rbp-proto/gen_pb/db/functor"
)

type voiceLoverSubjectDao struct {
}

var VoiceLoverSubjectDao = &voiceLoverSubjectDao{}

const (
	Subject_IsDeletedDefault = 0 // 是否已删除-未删除
	Subject_IsDeletedTrue    = 1 // 是否已删除-已删除
)

func (v *voiceLoverSubjectDao) GetValidSubjectList(ctx context.Context, page int, limit int) ([]*functor2.EntityVoiceLoverSubject, error) {
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * limit
	list, err := functor.VoiceLoverSubject.Ctx(ctx).
		Where(functor.VoiceLoverSubject.Columns.IsDeleted, Subject_IsDeletedDefault).
		Order(functor.VoiceLoverSubject.Columns.CreateTime, "desc").
		Offset(offset).
		Limit(limit).FindAll()
	if err != nil {
		return nil, err
	}
	return list, nil
}
