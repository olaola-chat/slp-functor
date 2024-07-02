package dao

import (
	"context"
	"fmt"

	"github.com/gogf/gf/frame/g"
	functor2 "github.com/olaola-chat/rbp-proto/dao/functor"
	"github.com/olaola-chat/rbp-proto/gen_pb/db/functor"
)

type voiceLoverAudioDao struct {
}

const (
	None = iota
	Dub
	Content
	Post
	Cover
)

const (
	AuditDefault = iota
	AuditPass
	AuditNoPass
)

var VoiceLoverAudioDao = &voiceLoverAudioDao{}

func (v *voiceLoverAudioDao) GetAudioDetailByAudioId(ctx context.Context, id uint64) (*functor.EntityVoiceLoverAudio, error) {
	res, err := functor2.VoiceLoverAudio.Ctx(ctx).Where(functor2.VoiceLoverAudio.Columns.ID, id).One()
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (v *voiceLoverAudioDao) GetAudioDetailsByAudioIds(ctx context.Context, ids []uint64) ([]*functor.EntityVoiceLoverAudio, error) {
	list, err := functor2.VoiceLoverAudio.Ctx(ctx).
		Where(fmt.Sprintf("%s IN (?)", functor2.VoiceLoverAudio.Columns.ID), ids).
		FindAll()
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (v *voiceLoverAudioDao) GetValidAudioListByIds(ctx context.Context, ids []uint64) ([]*functor.EntityVoiceLoverAudio, error) {
	list, err := functor2.VoiceLoverAudio.Ctx(ctx).
		Where(fmt.Sprintf("%s IN (?)", functor2.VoiceLoverAudio.Columns.ID), ids).
		Where(functor2.VoiceLoverAudio.Columns.AuditStatus, AuditPass).
		Order(functor2.VoiceLoverAudio.Columns.CreateTime, "desc").FindAll()
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (v *voiceLoverAudioDao) UpdateAudioById(ctx context.Context, id uint64, data g.Map) (int64, error) {
	sqlRes, err := functor2.VoiceLoverAudio.Ctx(ctx).Where(functor2.VoiceLoverAudio.Columns.ID, id).Update(data)
	if err != nil {
		return 0, err
	}
	affect, _ := sqlRes.RowsAffected()
	return affect, nil
}

func (v *voiceLoverAudioDao) GetValidUidsByUid(ctx context.Context, uid uint32) ([]*functor.EntityVoiceLoverAudio, error) {
	list, err := functor2.VoiceLoverAudio.Ctx(ctx).Fields(functor2.VoiceLoverAudio.Columns.PubUID).
		Where(fmt.Sprintf("%s != ?", functor2.VoiceLoverAudio.Columns.PubUID), uid).
		Where(functor2.VoiceLoverAudio.Columns.AuditStatus, AuditPass).Limit(1000).FindAll()
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (v *voiceLoverAudioDao) GetValidAudios(ctx context.Context) ([]*functor.EntityVoiceLoverAudio, error) {
	return functor2.VoiceLoverAudio.Ctx(ctx).
		Where(functor2.VoiceLoverAudio.Columns.AuditStatus, AuditPass).
		Order(functor2.VoiceLoverAudio.Columns.CreateTime, "desc").
		Limit(1000).
		FindAll()
}
