package dao

import (
	"context"
	"fmt"
	"time"

	"github.com/gogf/gf/frame/g"
	"github.com/olaola-chat/rbp-proto/dao/functor"
	functor2 "github.com/olaola-chat/rbp-proto/gen_pb/db/functor"
)

type voiceLoverAlbumDao struct {
}

var VoiceLoverAlbumDao = &voiceLoverAlbumDao{}

const (
	ChoiceDefault = 0 // 类型-默认
	ChoiceRec     = 1 // 类型-精选

	Album_IsDeletedDefault = 0 // 是否已删除-未删除
	Album_IsDeletedTrue    = 1 // 是否已删除-已删除
)

func (v *voiceLoverAlbumDao) GetValidAlbumListByChoice(ctx context.Context, choice uint32, page int, limit int) ([]*functor2.EntityVoiceLoverAlbum, error) {
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * limit
	list, err := functor.VoiceLoverAlbum.Ctx(ctx).
		Where(functor.VoiceLoverAlbum.Columns.Choice, choice).
		Where(functor.VoiceLoverAlbum.Columns.IsDeleted, Album_IsDeletedDefault).
		Order(functor.VoiceLoverAlbum.Columns.CreateTime, "desc").
		Offset(offset).
		Limit(limit).FindAll()
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (v *voiceLoverAlbumDao) GetValidAlbumListByIds(ctx context.Context, ids []uint64) ([]*functor2.EntityVoiceLoverAlbum, error) {
	list, err := functor.VoiceLoverAlbum.Ctx(ctx).
		Where(fmt.Sprintf("%s IN (?)", functor.VoiceLoverAlbum.Columns.ID), ids).
		Where(functor.VoiceLoverAlbum.Columns.IsDeleted, Album_IsDeletedDefault).
		Order(functor.VoiceLoverAlbum.Columns.CreateTime, "desc").FindAll()
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (v *voiceLoverAlbumDao) CreateAlbum(ctx context.Context, name string, intro string, cover string, opUid uint64) (int64, error) {
	data := &functor2.EntityVoiceLoverAlbum{
		Name:       name,
		Intro:      intro,
		Cover:      cover,
		OpUid:      opUid,
		Choice:     0,
		ChoiceTime: 0,
		CreateTime: uint64(time.Now().Unix()),
		UpdateTime: uint64(time.Now().Unix()),
	}
	sqlRes, err := functor.VoiceLoverAlbum.Ctx(ctx).Insert(data)
	if err != nil {
		return 0, err
	}
	lastId, _ := sqlRes.LastInsertId()
	return lastId, nil
}

func (v *voiceLoverAlbumDao) GetValidAlbumById(ctx context.Context, id uint64) (*functor2.EntityVoiceLoverAlbum, error) {
	data, err := functor.VoiceLoverAlbum.Ctx(ctx).Where("id", id).Where("is_deleted", 0).One()
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (v *voiceLoverAlbumDao) DelAlbum(ctx context.Context, id uint64, opUid uint64) error {
	_, err := functor.VoiceLoverAlbum.Ctx(ctx).Where("id", id).Update(
		g.Map{
			"update_time": time.Now().Unix(),
			"op_uid":      opUid,
			"is_deleted":  1,
		})
	if err != nil {
		return err
	}
	return nil
}

func (v *voiceLoverAlbumDao) UpdateAlbum(ctx context.Context, id uint64, name string, intro string, cover string, opUid uint64) error {
	data := g.Map{
		"update_time": time.Now().Unix(),
		"op_uid":      opUid,
	}
	if len(name) > 0 {
		data["name"] = name
	}
	if len(intro) > 0 {
		data["intro"] = intro
	}
	if len(cover) > 0 {
		data["cover"] = cover
	}
	_, err := functor.VoiceLoverAlbum.Ctx(ctx).Where("id", id).Update(data)
	if err != nil {
		return err
	}
	return nil
}

func (v *voiceLoverAlbumDao) GetValidAlbumByName(ctx context.Context, name string) (*functor2.EntityVoiceLoverAlbum, error) {
	data, err := functor.VoiceLoverAlbum.Ctx(ctx).Where("name", name).Where("is_deleted", 0).One()
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (v *voiceLoverAlbumDao) GetValidAlbumList(ctx context.Context, startTime uint64, endTime uint64, name string, page int, limit int) ([]*functor2.EntityVoiceLoverAlbum, uint32, error) {
	if endTime == 0 {
		endTime = uint64(time.Now().Unix())
	}
	d := functor.VoiceLoverAlbum.Ctx(ctx).Where("create_time > ?", startTime).
		Where("create_time < ?", endTime).
		Where("is_deleted", 0)
	if len(name) > 0 {
		d = d.Where("name", name)
	}
	total, _ := d.Count()
	list, err := d.Page(page, limit).Order("id desc").FindAll()
	if err != nil {
		return nil, 0, err
	}
	return list, uint32(total), nil
}
