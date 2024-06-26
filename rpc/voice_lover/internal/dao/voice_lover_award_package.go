package dao

import (
	"context"
	"time"

	functor2 "github.com/olaola-chat/rbp-proto/dao/functor"
	"github.com/olaola-chat/rbp-proto/gen_pb/db/functor"
)

var VoiceLoverAwardPackageDao = &voiceLoverAwardPackageDao{}

type voiceLoverAwardPackageDao struct{}

func (v *voiceLoverAwardPackageDao) Create(ctx context.Context, name string, awards string) (int64, error) {
	now := time.Now().Unix()
	data := &functor.EntityVoiceLoverAwardPackage{
		Name:       name,
		Awards:     awards,
		CreateTime: uint32(now),
		UpdateTime: uint32(now),
	}
	res, err := functor2.VoiceLoverAwardPackage.Ctx(ctx).Insert(data)
	if err != nil {
		return 0, err
	}
	lastId, _ := res.LastInsertId()
	return lastId, nil
}
