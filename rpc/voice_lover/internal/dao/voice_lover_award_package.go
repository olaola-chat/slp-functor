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
		dao = dao.Where("name = ?", name)
	}
	list := ([]*functor.EntityVoiceLoverAwardPackage)(nil)
	err := dao.Order("id desc").Page(page, limit).Structs(&list)
	total, _ := dao.Count()
	if err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func (v *voiceLoverAwardPackageDao) Delete(ctx context.Context, id uint32) error {
	_, err := functor2.VoiceLoverAwardPackage.Ctx(ctx).Where("id = ?", id).Delete()
	return err
}

func (v *voiceLoverAwardPackageDao) BatchGet(ctx context.Context, ids []uint32) (map[uint32]*functor.EntityVoiceLoverAwardPackage, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	data, err := functor2.VoiceLoverAwardPackage.Ctx(ctx).Where("id in (?)", ids).FindAll()
	if err != nil {
		return nil, err
	}
	res := make(map[uint32]*functor.EntityVoiceLoverAwardPackage)
	for _, v := range data {
		res[v.GetId()] = v
	}
	return res, nil
}

func (v *voiceLoverAwardPackageDao) GetOne(ctx context.Context, id uint32) (*functor.EntityVoiceLoverAwardPackage, error) {
	return functor2.VoiceLoverAwardPackage.Ctx(ctx).Where("id = ?", id).FindOne()
}
