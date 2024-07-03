package dao

import (
	"context"
	"strings"

	config2 "github.com/olaola-chat/rbp-proto/dao/config"
	"github.com/olaola-chat/rbp-proto/gen_pb/db/config"
)

var VoiceLoverActivityRankAwardDao = &voiceLoverActivityRankAwardDao{}

type voiceLoverActivityRankAwardDao struct{}

// Upsert 添加/更新排行奖励
func (v *voiceLoverActivityRankAwardDao) Upsert(ctx context.Context, data *config.EntityVoiceLoverActivityRankAward) (uint32, error) {
	if data.GetId() > 0 {
		_, err := config2.VoiceLoverActivityRankAward.Ctx(ctx).Where("id = ?", data.GetId()).Data(data).Update()
		if err != nil {
			return 0, err
		}
		return data.GetId(), nil
	}

	res, err := config2.VoiceLoverActivityRankAward.Ctx(ctx).Insert(data)
	if err != nil {
		return 0, err
	}
	lastId, _ := res.LastInsertId()
	return uint32(lastId), nil
}

func (v *voiceLoverActivityRankAwardDao) GetList(ctx context.Context, id uint32, name string, page, limit int) ([]*config.EntityVoiceLoverActivityRankAward, int, error) {
	dao := config2.VoiceLoverActivityRankAward.Ctx(ctx)
	if id > 0 {
		dao = dao.Where("id = ?", id)
	}
	if name = strings.TrimSpace(name); name != "" {
		dao = dao.Where("name = ?", name)
	}
	total, _ := dao.Count()
	list, err := dao.Order("id desc").Page(page, limit).FindAll()
	if err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func (v *voiceLoverActivityRankAwardDao) Delete(ctx context.Context, id uint32) error {
	_, err := config2.VoiceLoverActivityRankAward.Ctx(ctx).Where("id = ?", id).Delete()
	return err
}

func (v *voiceLoverActivityRankAwardDao) BatchGet(ctx context.Context, ids []uint32) (map[uint32]*config.EntityVoiceLoverActivityRankAward, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	data, err := config2.VoiceLoverActivityRankAward.Ctx(ctx).Where("id in (?)", ids).FindAll()
	if err != nil {
		return nil, err
	}
	res := make(map[uint32]*config.EntityVoiceLoverActivityRankAward)
	for _, v := range data {
		res[v.GetId()] = v
	}
	return res, nil
}
