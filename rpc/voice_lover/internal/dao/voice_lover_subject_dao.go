package dao

import (
	"context"
	"time"

	"github.com/gogf/gf/frame/g"
	"github.com/olaola-chat/rbp-proto/dao/functor"
	functor2 "github.com/olaola-chat/rbp-proto/gen_pb/db/functor"
)

type voiceLoverSubjectDao struct {
}

var VoiceLoverSubjectDao = &voiceLoverSubjectDao{}

const (
	Subject_IsDeletedDefault = 0 // 是否已删除-未删除
	Subject_IsDeletedTrue    = 1 // 是否已删除-已删除
)

func (v *voiceLoverSubjectDao) GetValidSubjectList(ctx context.Context, page int, limit int) ([]*functor2.EntityVoiceLoverSubject, error) {
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * limit
	list, err := functor.VoiceLoverSubject.Ctx(ctx).
		Where(functor.VoiceLoverSubject.Columns.IsDeleted, Subject_IsDeletedDefault).
		Order(functor.VoiceLoverSubject.Columns.CreateTime, "desc").
		Offset(offset).
		Limit(limit).FindAll()
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (v *voiceLoverSubjectDao) CreateSubject(ctx context.Context, name string, opUid uint64) (int64, error) {
	data := &functor2.EntityVoiceLoverSubject{
		Name:       name,
		OpUid:      opUid,
		CreateTime: uint64(time.Now().Unix()),
		UpdateTime: uint64(time.Now().Unix()),
	}
	sqlRes, err := functor.VoiceLoverSubject.Ctx(ctx).Insert(data)
	if err != nil {
		return 0, err
	}
	lastId, _ := sqlRes.LastInsertId()
	return lastId, nil
}

func (v *voiceLoverSubjectDao) UpdateSubject(ctx context.Context, id uint64, name string, opUid uint64) error {
	if len(name) == 0 {
		return nil
	}
	data := g.Map{
		"update_time": time.Now().Unix(),
		"name":        name,
		"op_uid":      opUid,
	}
	_, err := functor.VoiceLoverSubject.Ctx(ctx).Where("id", id).Update(data)
	return err
}

func (v *voiceLoverSubjectDao) DelSubject(ctx context.Context, id uint64, opUid uint64) error {
	data := g.Map{
		"update_time": time.Now().Unix(),
		"op_uid":      opUid,
		"is_deleted":  Subject_IsDeletedTrue,
	}
	_, err := functor.VoiceLoverSubject.Ctx(ctx).Where("id", id).Update(data)
	return err
}

func (v *voiceLoverSubjectDao) GetValidSubjectByIds(ctx context.Context, ids []uint64) (map[uint64]*functor2.EntityVoiceLoverSubject, error) {
	datas, err := functor.VoiceLoverSubject.Ctx(ctx).Where("id in (?)", ids).Where("is_deleted", Subject_IsDeletedDefault).FindAll()
	if err != nil {
		return nil, err
	}
	r := make(map[uint64]*functor2.EntityVoiceLoverSubject)
	for _, data := range datas {
		r[data.Id] = data
	}
	return r, nil
}

func (v *voiceLoverSubjectDao) GetValidSubjectListByName(ctx context.Context, startTime uint64, endTime uint64, name string, page, limit int) ([]*functor2.EntityVoiceLoverSubject, int, error) {
	if endTime == 0 {
		endTime = uint64(time.Now().Unix())
	}
	d := functor.VoiceLoverSubject.Ctx(ctx).Where("create_time > ?", startTime).Where("create_time < ?", endTime).
		Where(functor.VoiceLoverSubject.Columns.IsDeleted, Subject_IsDeletedDefault)
	if len(name) > 0 {
		d = d.Where("name", name)
	}
	count, _ := d.Count()
	list, err := d.Order("id desc").Page(page, limit).FindAll()
	if err != nil {
		return nil, 0, err
	}
	return list, count, nil
}

func (v *voiceLoverSubjectDao) GetValidSubjectByName(ctx context.Context, name string) (*functor2.EntityVoiceLoverSubject, error) {
	subject, err := functor.VoiceLoverSubject.Ctx(ctx).Where("name", name).Where(functor.VoiceLoverSubject.Columns.IsDeleted, Subject_IsDeletedDefault).One()
	if err != nil {
		return nil, err
	}
	return subject, nil
}

func (v *voiceLoverSubjectDao) GetValidSubjectById(ctx context.Context, id uint64) (*functor2.EntityVoiceLoverSubject, error) {
	subject, err := functor.VoiceLoverSubject.Ctx(ctx).Where("id", id).Where(functor.VoiceLoverSubject.Columns.IsDeleted, Subject_IsDeletedDefault).One()
	if err != nil {
		return nil, err
	}
	return subject, nil
}
