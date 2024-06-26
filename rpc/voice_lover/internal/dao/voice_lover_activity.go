package dao

import (
	"context"

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
