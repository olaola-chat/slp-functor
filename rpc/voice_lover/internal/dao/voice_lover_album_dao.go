package dao

import (
	"context"
	"fmt"

	"github.com/olaola-chat/rbp-proto/dao/functor"
	functor2 "github.com/olaola-chat/rbp-proto/gen_pb/db/functor"
)

type voiceLoverAlbumDao struct {
}

var VoiceLoverAlbumDao = &voiceLoverAlbumDao{}

const (
	ChoiceDefault = 0 // 类型-默认
	ChoiceRec     = 1 // 类型-精选

	Album_IsDeletedDefault = 0 // 是否已删除-未删除
	Album_IsDeletedTrue    = 1 // 是否已删除-已删除
)

func (v *voiceLoverAlbumDao) GetValidAlbumListByChoice(ctx context.Context, choice uint32, page int, limit int) ([]*functor2.EntityVoiceLoverAlbum, error) {
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * limit
	list, err := functor.VoiceLoverAlbum.Ctx(ctx).
		Where(functor.VoiceLoverAlbum.Columns.Choice, choice).
		Where(functor.VoiceLoverAlbum.Columns.IsDeleted, Album_IsDeletedDefault).
		Order(functor.VoiceLoverAlbum.Columns.CreateTime, "desc").
		Offset(offset).
		Limit(limit).FindAll()
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (v *voiceLoverAlbumDao) GetValidAlbumListByIds(ctx context.Context, ids []uint64) ([]*functor2.EntityVoiceLoverAlbum, error) {
	list, err := functor.VoiceLoverAlbum.Ctx(ctx).
		Where(fmt.Sprintf("%s IN (?)", functor.VoiceLoverAlbum.Columns.ID), ids).
		Where(functor.VoiceLoverAlbum.Columns.IsDeleted, Album_IsDeletedDefault).
		Order(functor.VoiceLoverAlbum.Columns.CreateTime, "desc").FindAll()
	if err != nil {
		return nil, err
	}
	return list, nil
}
