package dao

import (
	"context"

	"github.com/gogf/gf/frame/g"

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
	//dao := functor2.VoiceLoverAwardPackage.Ctx(ctx)
	//if id > 0 {
	//	dao = dao.Where("id = ?", id)
	//}
	//if name = strings.TrimSpace(name); name != "" {
	//	dao = dao.Where("name = ?", name)
	//}
	//data, err := dao.Order("id desc").Page(page, limit).FindAll() // TODO(tanlian)
	data, err := functor2.VoiceLoverAwardPackage.FindArray()
	total := 1
	//total, _ := dao.Count()
	g.Log().Infof("data: %+v, len: %d", data, len(data))
	if err != nil {
		return nil, 0, err
	}
	return nil, total, nil
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
