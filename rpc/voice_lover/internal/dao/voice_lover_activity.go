package dao

import (
	"context"
	"strings"

	functor2 "github.com/olaola-chat/rbp-proto/dao/functor"
	"github.com/olaola-chat/rbp-proto/gen_pb/db/functor"
)

type voiceLoverActivityDao struct{}

var VoiceLoverActivityDao = &voiceLoverActivityDao{}

// Upsert 创建/更新活动
func (v *voiceLoverActivityDao) Upsert(ctx context.Context, data *functor.EntityVoiceLoverActivity) (uint32, error) {
	if data.GetId() > 0 {
		_, err := functor2.VoiceLoverActivity.Ctx(ctx).Where("id = ?", data.GetId()).Data(data).Update()
		if err != nil {
			return 0, err
		}
		return data.GetId(), nil
	}

	res, err := functor2.VoiceLoverActivity.Ctx(ctx).Insert(data)
	if err != nil {
		return 0, err
	}
	lastId, _ := res.LastInsertId()
	return uint32(lastId), nil
}

func (v *voiceLoverActivityDao) GetList(ctx context.Context, id uint32, title string, page, limit int) ([]*functor.EntityVoiceLoverActivity, int, error) {
	dao := functor2.VoiceLoverActivity.Ctx(ctx)
	if id > 0 {
		dao = dao.Where("id = ?", id)
	}
	if title = strings.TrimSpace(title); title != "" {
		dao = dao.Where("title like %?%", title)
	}
	total, _ := dao.Count()
	data, err := dao.Order("id desc").Page(page, limit).FindAll()
	if err != nil {
		return nil, 0, err
	}
	return data, total, nil
}

func (v *voiceLoverActivityDao) Delete(ctx context.Context, id uint32) error {
	_, err := functor2.VoiceLoverActivity.Ctx(ctx).Where("id = ?", id).Delete()
	return err
}
