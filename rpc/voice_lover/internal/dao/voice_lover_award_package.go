package dao

import (
	"context"
	"strings"

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

func (v *voiceLoverAwardPackageDao) GetList(ctx context.Context, id uint32, name string, page, limit int) ([]*functor.EntityVoiceLoverAwardPackage, int, error) {
	dao := functor2.VoiceLoverAwardPackage.Ctx(ctx)
	if id > 0 {
		dao = dao.Where("id = ?", id)
	}
	if name = strings.TrimSpace(name); name != "" {
		dao = dao.Where("name like %?%", name)
	}
	total, _ := dao.Count()
	data, err := dao.Order("id desc").Page(page, limit).FindAll()
	if err != nil {
		return nil, 0, err
	}
	return data, total, nil
}

func (v *voiceLoverAwardPackageDao) Delete(ctx context.Context, id uint32) error {
	_, err := functor2.VoiceLoverAwardPackage.Ctx(ctx).Where("id = ?", id).Delete()
	return err
}
