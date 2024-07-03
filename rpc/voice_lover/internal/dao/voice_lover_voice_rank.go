package dao

import (
	"context"
	"time"

	"github.com/gogf/gf/database/gdb"
	"github.com/gogf/gf/frame/g"
	"github.com/olaola-chat/rbp-proto/dao/xianshi"
)

var VoiceLoverVoiceRank = &voiceLoverVoiceRank{}

type voiceLoverVoiceRank struct{}

func (v *voiceLoverVoiceRank) IncrLikeNum(ctx context.Context, activityId uint32, audioId uint64) error {
	now := time.Now().Unix()
	s := "INSERT INTO voice_lover_voice_rank (`activity_id`,`audio_id`,`like_num`,`create_time`,`update_time`) VALUES (?,?,?,?,?) ON DUPLICATE KEY UPDATE like_num=like_num+?,update_time=?;"
	_, err := xianshi.VoiceLoverVoiceRank.DB.Ctx(ctx).Exec(s, activityId, audioId, 1, now, now, 1, now)
	return err
}

func (v *voiceLoverVoiceRank) DecLikeNum(ctx context.Context, activityId uint32, audioId uint64) error {
	now := time.Now().Unix()
	_, err := xianshi.VoiceLoverVoiceRank.Ctx(ctx).
		Where("activity_id = ? and audio_id = ?", activityId, audioId).
		Data(g.Map{"like_num": gdb.Raw("like_num-1"), "update_time": now}).
		Update()
	return err
}
