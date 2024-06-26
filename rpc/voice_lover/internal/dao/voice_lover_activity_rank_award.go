package dao

import (
	"context"

	functor2 "github.com/olaola-chat/rbp-proto/dao/functor"
	"github.com/olaola-chat/rbp-proto/gen_pb/db/functor"
)

var VoiceLoverActivityRankAwardDao = &voiceLoverActivityRankAwardDao{}

type voiceLoverActivityRankAwardDao struct{}

// Upsert 添加/更新排行奖励
func (v *voiceLoverActivityRankAwardDao) Upsert(ctx context.Context, data *functor.EntityVoiceLoverActivityRankAward) (uint32, error) {
	if data.GetId() > 0 {
		_, err := functor2.VoiceLoverActivityRankAward.Ctx(ctx).Where("id = ?", data.GetId()).Data(data).Update()
		if err != nil {
			return 0, err
		}
		return data.GetId(), nil
	}

	res, err := functor2.VoiceLoverActivityRankAward.Ctx(ctx).Insert(data)
	if err != nil {
		return 0, err
	}
	lastId, _ := res.LastInsertId()
	return uint32(lastId), nil
}
