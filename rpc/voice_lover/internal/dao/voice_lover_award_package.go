package dao

import (
	"context"
	functor2 "github.com/olaola-chat/rbp-proto/dao/functor"
	"github.com/olaola-chat/rbp-proto/gen_pb/db/functor"
)

var VoiceLoverAwardPackageDao = &voiceLoverAwardPackageDao{}

type voiceLoverAwardPackageDao struct{}

// Upsert 创建/更新奖励包
func (v *voiceLoverAwardPackageDao) Upsert(ctx context.Context, data *functor.EntityVoiceLoverAwardPackage) (uint32, error) {
	if data.GetId() > 0 {
		_, err := functor2.VoiceLoverAwardPackage.Ctx(ctx).Where("id = ?", data.GetId()).Data(data).Update()
		if err != nil {
			return 0, err
		}
		return data.GetId(), nil
	}

	res, err := functor2.VoiceLoverAwardPackage.Ctx(ctx).Insert(data)
	if err != nil {
		return 0, err
	}
	lastId, _ := res.LastInsertId()
	return uint32(lastId), nil
}
