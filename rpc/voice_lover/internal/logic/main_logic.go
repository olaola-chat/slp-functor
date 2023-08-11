package logic

import (
	"context"
	"strings"
	"sync"
	"time"

	"github.com/gogf/gf/database/gdb"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/util/gconv"
	"github.com/olaola-chat/rbp-library/es"
	"github.com/olaola-chat/rbp-proto/dao/functor"
	functor2 "github.com/olaola-chat/rbp-proto/gen_pb/db/functor"
	vl_pb "github.com/olaola-chat/rbp-proto/gen_pb/rpc/voice_lover"

	voice_lover2 "github.com/olaola-chat/rbp-functor/app/model/voice_lover"
	"github.com/olaola-chat/rbp-functor/rpc/voice_lover/internal/dao"
)

type mainLogic struct {
}

var MainLogic = &mainLogic{}

const (
	None = iota
	Dub
	Content
	Post
	Cover
)

func (m *mainLogic) Post(ctx context.Context, req *vl_pb.ReqPost, reply *vl_pb.ResBase) error {
	g.Log().Infof("VoiceLoverPost req = %v", req)
	now := uint64(time.Now().Unix())
	err := functor.VoiceLoverAudio.DB.Transaction(func(tx *gdb.TX) error {
		data := &functor2.EntityVoiceLoverAudio{
			Desc:       req.Desc,
			Resource:   req.Resource,
			Cover:      req.Cover,
			From:       uint32(req.Source),
			PubUid:     uint64(req.Uid),
			CreateTime: now,
			UpdateTime: now,
			Title:      req.Title,
			Labels:     req.Labels,
			Seconds:    req.Seconds,
		}
		last, err := functor.VoiceLoverAudio.TX(tx).Insert(data)
		if err != nil {
			return err
		}
		lastId, _ := last.LastInsertId()
		audioId := uint64(lastId)
		data.Id = audioId
		editDatas := make([]*functor2.EntityVoiceLoverAudioPartner, 0)
		editDubs := strings.Split(req.EditDub, ",")
		for _, editDub := range editDubs {
			uid := gconv.Uint64(editDub)
			if uid == 0 {
				continue
			}
			editDatas = append(editDatas, &functor2.EntityVoiceLoverAudioPartner{
				AudioId:    audioId,
				Type:       Dub,
				Uid:        uid,
				CreateTime: now,
				UpdateTime: now,
			})
		}
		editContents := strings.Split(req.EditContent, ",")
		for _, editContent := range editContents {
			uid := gconv.Uint64(editContent)
			if uid == 0 {
				continue
			}
			editDatas = append(editDatas, &functor2.EntityVoiceLoverAudioPartner{
				AudioId:    audioId,
				Type:       Content,
				Uid:        uid,
				CreateTime: now,
				UpdateTime: now,
			})
		}
		editPosts := strings.Split(req.EditPost, ",")
		for _, editPost := range editPosts {
			uid := gconv.Uint64(editPost)
			if uid == 0 {
				continue
			}
			editDatas = append(editDatas, &functor2.EntityVoiceLoverAudioPartner{
				AudioId:    audioId,
				Type:       Post,
				Uid:        uid,
				CreateTime: now,
				UpdateTime: now,
			})
		}
		editCovers := strings.Split(req.EditCover, ",")
		for _, editCover := range editCovers {
			uid := gconv.Uint64(editCover)
			if uid == 0 {
				continue
			}
			editDatas = append(editDatas, &functor2.EntityVoiceLoverAudioPartner{
				AudioId:    audioId,
				Type:       Cover,
				Uid:        uid,
				CreateTime: now,
				UpdateTime: now,
			})
		}
		if len(editDatas) > 0 {
			_, err := functor.VoiceLoverAudioPartner.TX(tx).Insert(editDatas)
			if err != nil {
				return err
			}
		}
		_ = es.EsClient(es.EsVpc).Put("voice_lover_audio", gconv.String(data.Id), gconv.Map(buildAudioEsModel(data)))
		return nil
	})
	if err != nil {
		g.Log().Errorf("VoiceLover post error, err = %v", err)
	}
	return err
}

func buildAudioEsModel(data *functor2.EntityVoiceLoverAudio) *voice_lover2.VoiceLoverAudioEsModel {
	labelsSlice := make([]string, 0)
	for _, l := range strings.Split(data.Labels, ",") {
		if len(l) == 0 {
			continue
		}
		labelsSlice = append(labelsSlice, l)
	}
	esModel := &voice_lover2.VoiceLoverAudioEsModel{
		Id:          data.Id,
		PubUid:      uint32(data.PubUid),
		Title:       data.Title,
		Cover:       data.Cover,
		Desc:        data.Desc,
		CreateTime:  data.CreateTime,
		Labels:      labelsSlice,
		Source:      int32(data.From),
		AuditStatus: int32(data.AuditStatus),
		Albums:      []uint64{},
		HasAlbum:    0,
		OpUid:       data.OpUid,
		Resource:    data.Resource,
		Seconds:     data.Seconds,
	}
	return esModel
}

func (m *mainLogic) BuildRecAlbumsExtendInfo(ctx context.Context, infos []*vl_pb.AlbumData) {
	countMap := make(map[uint64]uint32)
	wg := sync.WaitGroup{}
	for _, v := range infos {
		if _, ok := countMap[v.Id]; ok {
			continue
		}
		countMap[v.Id] = 0
		wg.Add(1)
		go func(albumId uint64) {
			defer wg.Done()
			total, _ := dao.VoiceLoverAudioAlbumDao.GetCountByAlbumId(ctx, albumId)
			countMap[albumId] = uint32(total)
		}(v.Id)
	}
	wg.Wait()
	for _, v := range infos {
		if count, ok := countMap[v.Id]; ok {
			v.AudioCount = count
		}
	}
}

func (m *mainLogic) GetRecAlbums(ctx context.Context, req *vl_pb.ReqGetRecAlbums, reply *vl_pb.ResGetRecAlbums) error {
	reply.Albums = make([]*vl_pb.AlbumData, 0)
	list, err := dao.VoiceLoverAlbumDao.GetValidAlbumListByChoice(ctx, dao.ChoiceRec, 0, 3)
	if err != nil {
		return err
	}
	for _, v := range list {
		reply.Albums = append(reply.Albums, &vl_pb.AlbumData{
			Id:         v.Id,
			Name:       v.Name,
			Intro:      v.Intro,
			Cover:      v.Cover,
			CreateTime: v.CreateTime,
		})
	}
	m.BuildRecAlbumsExtendInfo(ctx, reply.Albums)
	return nil
}

func (m *mainLogic) GetAlbumsByPage(ctx context.Context, req *vl_pb.ReqGetAlbumsByPage, reply *vl_pb.ResGetAlbumsByPage) error {
	reply.Albums = make([]*vl_pb.AlbumData, 0)
	list, err := dao.VoiceLoverAlbumDao.GetValidAlbumListByChoice(ctx, req.Choice, int(req.Page), int(req.Limit)+1)
	if err != nil {
		return err
	}
	if len(list) > int(req.Limit) {
		list = list[:req.Limit]
		reply.HasMore = true
	}
	for _, v := range list {
		reply.Albums = append(reply.Albums, &vl_pb.AlbumData{
			Id:         v.Id,
			Name:       v.Name,
			Intro:      v.Intro,
			Cover:      v.Cover,
			CreateTime: v.CreateTime,
		})
	}
	m.BuildRecAlbumsExtendInfo(ctx, reply.Albums)
	return nil
}

func (m *mainLogic) GetSubjectAlbumsByPage(ctx context.Context, req *vl_pb.ReqGetSubjectAlbumsByPage, reply *vl_pb.ResGetAlbumsByPage) error {
	reply.Albums = make([]*vl_pb.AlbumData, 0)
	list, err := dao.VoiceLoverAlbumSubjectDao.GetListBySubjectId(ctx, req.SubjectId, int(req.Page), int(req.Limit)+1)
	if err != nil {
		return err
	}
	if len(list) > int(req.Limit) {
		list = list[:req.Limit]
		reply.HasMore = true
	}
	albumIds := make([]uint64, 0)

	for _, v := range list {
		albumIds = append(albumIds, v.AlbumId)
	}
	albumList, err := dao.VoiceLoverAlbumDao.GetValidAlbumListByIds(ctx, albumIds)
	for _, v := range albumList {
		reply.Albums = append(reply.Albums, &vl_pb.AlbumData{
			Id:         v.Id,
			Name:       v.Name,
			Intro:      v.Intro,
			Cover:      v.Cover,
			CreateTime: v.CreateTime,
		})
	}
	m.BuildRecAlbumsExtendInfo(ctx, reply.Albums)
	return nil
}

func (m *mainLogic) GetRecSubjects(ctx context.Context, req *vl_pb.ReqGetRecSubjects, reply *vl_pb.ResGetRecSubjects) error {
	reply.Subjects = make([]*vl_pb.SubjectData, 0)
	list, err := dao.VoiceLoverSubjectDao.GetValidSubjectList(ctx, 0, 3)
	if err != nil {
		return err
	}
	subjectIds := make([]uint64, 0)
	for _, v := range list {
		subjectIds = append(subjectIds, v.Id)
		reply.Subjects = append(reply.Subjects, &vl_pb.SubjectData{
			Id:         v.Id,
			Name:       v.Name,
			CreateTime: v.CreateTime,
			Albums:     make([]*vl_pb.AlbumData, 0),
		})
	}

	albumSubjectRelList, err := dao.VoiceLoverAlbumSubjectDao.GetListBySubjectIds(ctx, subjectIds)
	if err != nil {
		return err
	}
	subjectAlbumsMap := make(map[uint64][]uint64)
	albumIds := make([]uint64, 0)
	for _, v := range albumSubjectRelList {
		subjectAlbumsMap[v.SubjectId] = append(subjectAlbumsMap[v.SubjectId], v.AlbumId)
		albumIds = append(albumIds, v.AlbumId)
	}
	albumList, err := dao.VoiceLoverAlbumDao.GetValidAlbumListByIds(ctx, albumIds)
	if err != nil {
		return err
	}
	albums := make([]*vl_pb.AlbumData, 0)
	albumsMap := make(map[uint64]*vl_pb.AlbumData)
	for _, v := range albumList {
		albums = append(albums, &vl_pb.AlbumData{
			Id:         v.Id,
			Name:       v.Name,
			Intro:      v.Intro,
			Cover:      v.Cover,
			CreateTime: v.CreateTime,
		})
	}
	m.BuildRecAlbumsExtendInfo(ctx, albums)
	for _, v := range albums {
		albumsMap[v.Id] = v
	}

	for _, v := range reply.Subjects {
		if _, ok := subjectAlbumsMap[v.Id]; !ok {
			continue
		}
		for _, albumId := range subjectAlbumsMap[v.Id] {
			if _, ok := albumsMap[albumId]; !ok {
				continue
			}
			v.Albums = append(v.Albums, albumsMap[albumId])
		}
	}
	return nil
}

func (m *mainLogic) BatchGetAlbumAudioCount(ctx context.Context, req *vl_pb.ReqBatchGetAlbumAudioCount, reply *vl_pb.ResBatchGetAlbumAudioCount) error {
	reply.AlbumCounts = make(map[uint64]uint32)
	wg := sync.WaitGroup{}
	for _, v := range req.AlbumIds {
		if _, ok := reply.AlbumCounts[v]; ok {
			continue
		}
		reply.AlbumCounts[v] = 0
		wg.Add(1)
		go func(albumId uint64) {
			defer wg.Done()
			total, _ := dao.VoiceLoverAudioAlbumDao.GetCountByAlbumId(ctx, albumId)
			reply.AlbumCounts[albumId] = uint32(total)
		}(v)
	}
	wg.Wait()
	return nil
}
