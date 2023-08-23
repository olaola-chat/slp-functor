package dao

import (
	"context"
	"time"

	functor2 "github.com/olaola-chat/rbp-proto/dao/functor"
	"github.com/olaola-chat/rbp-proto/gen_pb/db/functor"
)

type voiceLoverBannerDao struct {
}

var VoiceLoverBannerDao = &voiceLoverBannerDao{}

func (v *voiceLoverBannerDao) GetBannerList(ctx context.Context, startTime, endTime uint64, title string, status int32, page, limit int) ([]*functor.EntityVoiceLoverBanner, int32, error) {

	d := functor2.VoiceLoverBanner.Ctx(ctx)
	if startTime > 0 {
		d = d.Where("start_time > ?", startTime)
	}
	if endTime > 0 {
		d = d.Where("end_time < ?", endTime)
	}
	if len(title) > 0 {
		d = d.Where("title", title)
	}
	now := time.Now().Unix()
	if status == 1 {
		d = d.Where("start_time < ?", now).Where("end_time > ?", now)
	}
	if status == 2 {
		d = d.Where("start_time > ? or end_time < ?", now, now)
	}
	count, _ := d.Count()
	list, err := d.Order("sort asc").Page(page, limit).FindAll()
	if err != nil {
		return nil, 0, err
	}
	return list, int32(count), nil
}
