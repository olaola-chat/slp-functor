package dao

import (
	"context"
	"strings"

	config2 "github.com/olaola-chat/rbp-proto/dao/config"
	"github.com/olaola-chat/rbp-proto/gen_pb/db/config"
)

var VoiceLoverAwardPackageDao = &voiceLoverAwardPackageDao{}

type voiceLoverAwardPackageDao struct{}

// Upsert 创建/更新奖励包
func (v *voiceLoverAwardPackageDao) Upsert(ctx context.Context, data *config.EntityVoiceLoverAwardPackage) (uint32, error) {
	if data.GetId() > 0 {
		_, err := config2.VoiceLoverAwardPackage.Ctx(ctx).Where("id = ?", data.GetId()).Data(data).Update()
		if err != nil {
			return 0, err
		}
		return data.GetId(), nil
	}

	res, err := config2.VoiceLoverAwardPackage.Ctx(ctx).Insert(data)
	if err != nil {
		return 0, err
	}
	lastId, _ := res.LastInsertId()
	return uint32(lastId), nil
}

func (v *voiceLoverAwardPackageDao) GetList(ctx context.Context, id uint32, name string, page, limit int) ([]*config.EntityVoiceLoverAwardPackage, int, error) {
	dao := config2.VoiceLoverAwardPackage.Ctx(ctx)
	if id > 0 {
		dao = dao.Where("id = ?", id)
	}
	if name = strings.TrimSpace(name); name != "" {
		dao = dao.Where("name = ?", name)
	}
	list, err := dao.Order("id desc").Page(page, limit).FindAll()
	total, _ := dao.Count()
	if err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func (v *voiceLoverAwardPackageDao) Delete(ctx context.Context, id uint32) error {
	_, err := config2.VoiceLoverAwardPackage.Ctx(ctx).Where("id = ?", id).Delete()
	return err
}

func (v *voiceLoverAwardPackageDao) BatchGet(ctx context.Context, ids []uint32) (map[uint32]*config.EntityVoiceLoverAwardPackage, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	list, err := config2.VoiceLoverAwardPackage.Ctx(ctx).Where("id in (?)", ids).FindAll()
	if err != nil {
		return nil, err
	}
	res := make(map[uint32]*config.EntityVoiceLoverAwardPackage)
	for _, v := range list {
		res[v.GetId()] = v
	}
	return res, nil
}

func (v *voiceLoverAwardPackageDao) GetOne(ctx context.Context, id uint32) (*config.EntityVoiceLoverAwardPackage, error) {
	return config2.VoiceLoverAwardPackage.Ctx(ctx).Where("id = ?", id).FindOne()
}
