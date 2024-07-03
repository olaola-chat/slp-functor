package dao

import (
	"context"
	"time"

	functor2 "github.com/olaola-chat/rbp-proto/dao/functor"
	"github.com/olaola-chat/rbp-proto/gen_pb/db/functor"
)

type voiceLoverUserCollectDao struct {
}

var VoiceLoverUserCollectDao = &voiceLoverUserCollectDao{}

const (
	CollectTypeAlbum = 0 // 专辑收藏
	CollectTypeAudio = 1 // 音频收藏
)

func (v *voiceLoverUserCollectDao) Add(ctx context.Context, uid uint32, collectId uint64, collectType int) (int64, error) {
	now := uint64(time.Now().Unix())
	data := &functor.EntityVoiceLoverUserCollect{
		Uid:         uint64(uid),
		CollectId:   collectId,
		CollectType: int32(collectType),
		CreateTime:  now,
		UpdateTime:  now,
	}
	res, err := functor2.VoiceLoverUserCollect.Ctx(ctx).Insert(data)
	if err != nil {
		return 0, err
	}
	lastId, _ := res.LastInsertId()
	return lastId, nil
}

func (v *voiceLoverUserCollectDao) Delete(ctx context.Context, uid uint32, collectId uint64, collectType int) error {
	_, err := functor2.VoiceLoverUserCollect.Ctx(ctx).
		Where(functor2.VoiceLoverUserCollect.Columns.UID, uid).
		Where(functor2.VoiceLoverUserCollect.Columns.CollectID, collectId).
		Where(functor2.VoiceLoverUserCollect.Columns.CollectType, collectType).Delete()
	if err != nil {
		return err
	}
	return nil
}

func (v *voiceLoverUserCollectDao) GetListByUidAndType(ctx context.Context, uid uint32, collectType int, page int, limit int) ([]*functor.EntityVoiceLoverUserCollect, error) {
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * limit
	list, err := functor2.VoiceLoverUserCollect.Ctx(ctx).
		Where(functor2.VoiceLoverUserCollect.Columns.UID, uid).
		Where(functor2.VoiceLoverUserCollect.Columns.CollectType, collectType).
		Offset(offset).
		Limit(limit).FindAll()
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (v *voiceLoverUserCollectDao) GetInfoByUidAndTypeAndId(ctx context.Context, uid uint32, id uint64, collectType int) (*functor.EntityVoiceLoverUserCollect, error) {
	data, err := functor2.VoiceLoverUserCollect.Ctx(ctx).
		Where(functor2.VoiceLoverUserCollect.Columns.UID, uid).
		Where(functor2.VoiceLoverUserCollect.Columns.CollectID, id).
		Where(functor2.VoiceLoverUserCollect.Columns.CollectType, collectType).
		FindOne()
	if err != nil {
		return nil, err
	}
	return data, nil
}
func (v *voiceLoverUserCollectDao) BatchCheckUserCollected(ctx context.Context, uid uint32, typ int, ids []uint32) (map[uint32]bool, error) {
	data, err := functor2.VoiceLoverUserCollect.Ctx(ctx).
		Where(functor2.VoiceLoverUserCollect.Columns.UID, uid).
		Where(functor2.VoiceLoverUserCollect.Columns.CollectType, typ).
		Where("collect_id in (?)", ids).
		FindAll()
	if err != nil {
		return nil, err
	}
	res := make(map[uint32]bool)
	for _, v := range data {
		res[uint32(v.GetCollectId())] = true
	}
	return res, nil
}

func (v *voiceLoverUserCollectDao) BatchGetCollectNum(ctx context.Context, ids []uint32) (map[uint32]uint32, error) {
	data, err := functor2.VoiceLoverUserCollect.Ctx(ctx).Where("collect_id in (?)", ids).FindAll()
	if err != nil {
		return nil, err
	}
	res := make(map[uint32]uint32)
	for _, v := range data {
		res[uint32(v.GetCollectId())]++
	}
	return res, nil
}
