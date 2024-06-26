package dao

import (
	"context"
	"time"

	functor2 "github.com/olaola-chat/rbp-proto/dao/functor"
	"github.com/olaola-chat/rbp-proto/gen_pb/db/functor"
)

type voiceLoverActivityDao struct{}

var VoiceLoverActivityDao = &voiceLoverActivityDao{}

// Add 添加活动
func (v *voiceLoverActivityDao) Add(ctx context.Context, title, intro, cover string, startTime, endTime uint32) (int64, error) {
	now := time.Now().Unix()
	data := &functor.EntityVoiceLoverActivity{
		Title:      title,
		Intro:      intro,
		Cover:      cover,
		StartTime:  startTime,
		EndTime:    endTime,
		CreateTime: uint32(now),
		UpdateTime: uint32(now),
	}
	res, err := functor2.VoiceLoverActivity.Ctx(ctx).Insert(data)
	if err != nil {
		return 0, err
	}
	lastId, _ := res.LastInsertId()
	return lastId, nil
}
