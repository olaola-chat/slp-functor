package dao

import (
	"context"
	"time"

	functor2 "github.com/olaola-chat/rbp-proto/dao/functor"
	"github.com/olaola-chat/rbp-proto/gen_pb/db/functor"
)

var VoiceLoverActivityRankAwardDao = &voiceLoverActivityRankAwardDao{}

type voiceLoverActivityRankAwardDao struct{}

func (v *voiceLoverActivityRankAwardDao) Create(ctx context.Context, name string, pkgId uint32, content string) (int64, error) {
	now := time.Now().Unix()
	data := &functor.EntityVoiceLoverActivityRankAward{
		Name:       name,
		PackageId:  pkgId,
		Content:    content,
		CreateTime: uint32(now),
		UpdateTime: uint32(now),
	}
	res, err := functor2.VoiceLoverActivityRankAward.Ctx(ctx).Insert(data)
	if err != nil {
		return 0, err
	}
	lastId, _ := res.LastInsertId()
	return lastId, nil
}
