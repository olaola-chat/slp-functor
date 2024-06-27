package dao

import (
	"context"
	"strings"

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

func (v *voiceLoverActivityRankAwardDao) GetList(ctx context.Context, id uint32, name string, page, limit int) ([]*functor.EntityVoiceLoverActivityRankAward, int, error) {
	dao := functor2.VoiceLoverActivityRankAward.Ctx(ctx)
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

func (v *voiceLoverActivityRankAwardDao) Delete(ctx context.Context, id uint32) error {
	_, err := functor2.VoiceLoverActivityRankAward.Ctx(ctx).Where("id = ?", id).Delete()
	return err
}

func (v *voiceLoverActivityRankAwardDao) BatchGet(ctx context.Context, ids []uint32) (map[uint32]*functor.EntityVoiceLoverActivityRankAward, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	data, err := functor2.VoiceLoverActivityRankAward.Ctx(ctx).Where("id in (?)", ids).FindAll()
	if err != nil {
		return nil, err
	}
	res := make(map[uint32]*functor.EntityVoiceLoverActivityRankAward)
	for _, v := range data {
		res[v.GetId()] = v
	}
	return res, nil
}
