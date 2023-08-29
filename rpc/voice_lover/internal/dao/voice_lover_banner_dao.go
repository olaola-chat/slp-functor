package dao

import (
	"context"
	"time"

	"github.com/gogf/gf/frame/g"
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

func (v *voiceLoverBannerDao) GetValidListByLimit(ctx context.Context, limit int) ([]*functor.EntityVoiceLoverBanner, error) {
	now := time.Now().Unix()
	list, err := functor2.VoiceLoverBanner.Ctx(ctx).
		Where("start_time < ?", now).
		Where("end_time > ?", now).
		Limit(limit).
		Order(functor2.VoiceLoverBanner.Columns.Sort, "desc").
		FindAll()
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (v *voiceLoverBannerDao) CreateBanner(ctx context.Context, startTime, endTime uint64, title, cover, schema string, opUid uint64, sort uint32) (uint64, error) {
	d := functor2.VoiceLoverBanner.Ctx(ctx)
	b := &functor.EntityVoiceLoverBanner{
		Title:      title,
		Cover:      cover,
		Schema:     schema,
		OpUid:      opUid,
		StartTime:  startTime,
		EndTime:    endTime,
		Sort:       sort,
		CreateTime: uint64(time.Now().Unix()),
		UpdateTime: uint64(time.Now().Unix()),
	}
	sqlRes, err := d.Insert(b)
	if err != nil {
		return 0, err
	}
	id, _ := sqlRes.LastInsertId()
	return uint64(id), nil
}

func (v *voiceLoverBannerDao) UpdateBanner(ctx context.Context, id, startTime, endTime uint64, title, cover, schema string, opUid uint64, sort uint32) error {
	d := functor2.VoiceLoverBanner.Ctx(ctx)
	data := g.Map{}
	data["start_time"] = startTime
	data["end_time"] = endTime
	data["title"] = title
	data["cover"] = cover
	data["schema"] = schema
	data["op_uid"] = opUid
	data["sort"] = sort
	data["update_time"] = time.Now().Unix()
	_, err := d.Where("id", id).Update(data)
	return err
}

func (v *voiceLoverBannerDao) GetBannerById(ctx context.Context, id uint64) (*functor.EntityVoiceLoverBanner, error) {
	data, err := functor2.VoiceLoverBanner.Ctx(ctx).Where("id", id).One()
	if err != nil {
		return nil, err
	}
	return data, nil
}
