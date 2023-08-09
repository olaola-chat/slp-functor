package dao

import (
	"context"

	"github.com/gogf/gf/frame/g"
	"github.com/olaola-chat/rbp-proto/dao/functor"
	functor2 "github.com/olaola-chat/rbp-proto/gen_pb/db/functor"
)

type voiceLoverAlbumDao struct {
}

var VoiceLoverAlbumDao = &voiceLoverAlbumDao{}

const (
	ChoiceDefault = 0 // 类型-默认
	ChoiceRec     = 1 // 类型-精选
)

func (v *voiceLoverAlbumDao) GetAlbumListByChoice(ctx context.Context, choice uint32, page int, limit int) ([]*functor2.EntityVoiceLoverAlbum, error) {
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * limit
	list, err := functor.VoiceLoverAlbum.Ctx(ctx).
		Where(functor.VoiceLoverAlbum.Columns.Choice, choice).
		Order(functor.VoiceLoverAlbum.Columns.CreateTime, "desc").
		Offset(offset).
		Limit(limit).FindAll()
	if err != nil {
		g.Log().Errorf("voiceLoverAlbumDao GetAlbumListByChoice FindAll error=%v", err)
		return nil, err
	}
	return list, nil
}
