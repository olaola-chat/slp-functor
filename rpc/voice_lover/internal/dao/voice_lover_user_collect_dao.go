package dao

import (
	"context"

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
