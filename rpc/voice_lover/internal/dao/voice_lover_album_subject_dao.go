package dao

import (
	"context"
	"fmt"
	"time"

	functor2 "github.com/olaola-chat/rbp-proto/dao/functor"
	"github.com/olaola-chat/rbp-proto/gen_pb/db/functor"
)

type voiceLoverAlbumSubjectDao struct {
}

var VoiceLoverAlbumSubjectDao = &voiceLoverAlbumSubjectDao{}

func (v *voiceLoverAlbumSubjectDao) GetListBySubjectId(ctx context.Context, subjectId uint64, page int, limit int) ([]*functor.EntityVoiceLoverAlbumSubject, error) {
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * limit
	list, err := functor2.VoiceLoverAlbumSubject.Ctx(ctx).
		Where(functor2.VoiceLoverAlbumSubject.Columns.SubjectID, subjectId).
		Order(functor2.VoiceLoverAlbum.Columns.CreateTime, "desc").
		Offset(offset).
		Limit(limit).FindAll()
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (v *voiceLoverAlbumSubjectDao) GetListBySubjectIds(ctx context.Context, subjectIds []uint64) ([]*functor.EntityVoiceLoverAlbumSubject, error) {
	list, err := functor2.VoiceLoverAlbumSubject.Ctx(ctx).
		Where(fmt.Sprintf("%s IN (?)", functor2.VoiceLoverAlbumSubject.Columns.SubjectID), subjectIds).
		Order(functor2.VoiceLoverAlbumSubject.Columns.CreateTime, "desc").FindAll()
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (v *voiceLoverAlbumSubjectDao) GetCountBySubjectId(ctx context.Context, id uint64) (int, error) {
	count, err := functor2.VoiceLoverAlbumSubject.Ctx(ctx).Where("subject_id", id).Count()
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (v *voiceLoverAlbumSubjectDao) GetCountByAlbumId(ctx context.Context, id uint64) (int, error) {
	count, err := functor2.VoiceLoverAlbumSubject.Ctx(ctx).Where("album_id", id).Count()
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (v *voiceLoverAlbumSubjectDao) GetAlbumSubjectByAIdAndSId(ctx context.Context, albumId uint64, subjectId uint64) (*functor.EntityVoiceLoverAlbumSubject, error) {
	data, err := functor2.VoiceLoverAlbumSubject.Ctx(ctx).Where("album_id", albumId).Where("subject_id", subjectId).One()
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (v *voiceLoverAlbumSubjectDao) Create(ctx context.Context, albumId uint64, subjectId uint64) error {
	data := &functor.EntityVoiceLoverAlbumSubject{
		AlbumId:    albumId,
		SubjectId:  subjectId,
		CreateTime: uint64(time.Now().Unix()),
		UpdateTime: uint64(time.Now().Unix()),
	}
	_, err := functor2.VoiceLoverAlbumSubject.Ctx(ctx).Insert(data)
	return err
}

func (v *voiceLoverAlbumSubjectDao) GetAlbumCollect(ctx context.Context, albumId uint64, subjectId uint64, page int32, limit int32) ([]*functor.EntityVoiceLoverAlbumSubject, int, error) {
	d := functor2.VoiceLoverAlbumSubject.Ctx(ctx)
	if albumId > 0 {
		d = d.Where("album_id", albumId)
	}
	if subjectId > 0 {
		d = d.Where("subject_id", subjectId)
	}
	total, _ := d.Count()
	list, err := d.Order("id desc").Page(int(page), int(limit)).FindAll()
	if err != nil {
		return nil, 0, err
	}
	return list, total, nil
}
