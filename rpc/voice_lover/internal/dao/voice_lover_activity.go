package dao

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/olaola-chat/rbp-proto/gen_pb/db/config"

	config2 "github.com/olaola-chat/rbp-proto/dao/config"
)

type voiceLoverActivityDao struct{}

var VoiceLoverActivityDao = &voiceLoverActivityDao{}

// Upsert 创建/更新活动
func (v *voiceLoverActivityDao) Upsert(ctx context.Context, data *config.EntityVoiceLoverActivity) (uint32, error) {
	if data.GetId() > 0 {
		_, err := config2.VoiceLoverActivity.Ctx(ctx).Where("id = ?", data.GetId()).Data(data).Update()
		if err != nil {
			return 0, err
		}
		return data.GetId(), nil
	}

	res, err := config2.VoiceLoverActivity.Ctx(ctx).Insert(data)
	if err != nil {
		return 0, err
	}
	lastId, _ := res.LastInsertId()
	return uint32(lastId), nil
}

func (v *voiceLoverActivityDao) GetList(ctx context.Context, id uint32, title string, page, limit int) ([]*config.EntityVoiceLoverActivity, int, error) {
	dao := config2.VoiceLoverActivity.Ctx(ctx)
	if id > 0 {
		dao = dao.Where("id = ?", id)
	}
	if title = strings.TrimSpace(title); title != "" {
		dao = dao.Where("title = ?", title)
	}
	total, _ := dao.Count()
	list, err := dao.Order("id desc").Page(page, limit).FindAll()
	if err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func (v *voiceLoverActivityDao) Delete(ctx context.Context, id uint32) error {
	_, err := config2.VoiceLoverActivity.Ctx(ctx).Where("id = ?", id).Delete()
	return err
}

func (v *voiceLoverActivityDao) GetOne(ctx context.Context, id uint32) (*config.EntityVoiceLoverActivity, error) {
	return config2.VoiceLoverActivity.Ctx(ctx).
		Cache(time.Minute, fmt.Sprintf("voice.lover.activity.%d", id)).
		Where("id = ?", id).
		FindOne()
}
