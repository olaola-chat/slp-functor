package logic

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/gogf/gf/frame/g"
	"github.com/olaola-chat/rbp-library/es"
	"github.com/olaola-chat/rbp-proto/gen_pb/db/functor"
	"github.com/olaola-chat/rbp-proto/gen_pb/rpc/voice_lover"

	"github.com/olaola-chat/rbp-functor/app/consts"
	"github.com/olaola-chat/rbp-functor/rpc/voice_lover/internal/dao"
)

type adminLogic struct {
}

var AdminLogic = &adminLogic{}

const (
	Collect = iota
	CollectRemove
)

func (a *adminLogic) GetAudioDetail(ctx context.Context, id uint64) (*voice_lover.AudioData, error) {
	res, err := dao.VoiceLoverAudioDao.GetAudioDetailByAudioId(ctx, id)
	if err != nil {
		return nil, err
	}
	if res == nil {
		return nil, nil
	}
	covers := make([]string, 0)
	for _, s := range strings.Split(res.Cover, ",") {
		if len(s) == 0 {
			continue
		}
		covers = append(covers, s)
	}
	labels := make([]string, 0)
	for _, s := range strings.Split(res.Labels, ",") {
		if len(s) == 0 {
			continue
		}
		labels = append(labels, s)
	}
	audio := &voice_lover.AudioData{
		Id:          res.Id,
		Uid:         uint32(res.PubUid),
		Resource:    res.Resource,
		Covers:      covers,
		Title:       res.Title,
		Desc:        res.Title,
		Labels:      labels,
		AuditStatus: int32(res.AuditStatus),
		CreateTime:  res.CreateTime,
		OpUid:       res.OpUid,
		Seconds:     res.Seconds,
	}
	edit, err := dao.VoiceLoverAudioPartnerDao.GetAudioPartnerByAudioId(ctx, id)
	if err != nil {
		return nil, err
	}
	for _, e := range edit {
		if e.Type == Dub {
			audio.EditDubs = append(audio.EditDubs, &voice_lover.AudioEditData{
				Uid:  uint32(e.Uid),
				Type: e.Type,
			})
		}
		if e.Type == Content {
			audio.EditContents = append(audio.EditContents, &voice_lover.AudioEditData{
				Uid:  uint32(e.Uid),
				Type: e.Type,
			})
		}
		if e.Type == Post {
			audio.EditPosts = append(audio.EditPosts, &voice_lover.AudioEditData{
				Uid:  uint32(e.Uid),
				Type: e.Type,
			})
		}
		if e.Type == Cover {
			audio.EditCovers = append(audio.EditCovers, &voice_lover.AudioEditData{
				Uid:  uint32(e.Uid),
				Type: e.Type,
			})
		}
	}
	return audio, nil
}

func (a *adminLogic) UpdateAudio(ctx context.Context, req *voice_lover.ReqUpdateAudio) error {
	data := g.Map{}
	if len(req.Title) > 0 {
		data["title"] = req.Title
	}
	if len(req.Desc) > 0 {
		data["desc"] = req.Desc
	}
	data["labels"] = req.Labels
	data["update_time"] = time.Now().Unix()
	data["op_uid"] = req.OpUid
	affect, err := dao.VoiceLoverAudioDao.UpdateAudioById(ctx, req.Id, data)
	if err != nil {
		return err
	}
	if affect > 0 {
		delete(data, "update_time")
		delete(data, "labels")
		labelsSlice := make([]string, 0)
		for _, l := range strings.Split(req.Labels, ",") {
			if len(l) == 0 {
				continue
			}
			labelsSlice = append(labelsSlice, l)
		}
		data["labels"] = labelsSlice
		_ = es.EsClient(es.EsVpc).Update("voice_lover_audio", req.Id, data)
	}
	return nil
}

func (a *adminLogic) AuditAudio(ctx context.Context, req *voice_lover.ReqAuditAudio) error {
	if req.AuditStatus != dao.AuditNoPass && req.AuditStatus != dao.AuditPass {
		return consts.ERROR_AUDIT_AUDIO_PARAM
	}
	res, err := dao.VoiceLoverAudioDao.GetAudioDetailByAudioId(ctx, req.GetId())
	if err != nil {
		return err
	}
	if res == nil {
		return consts.ERROR_AUDIO_NOT_EXIST
	}
	if req.AuditStatus == int32(res.AuditStatus) {
		return nil
	}
	if res.AuditStatus == dao.AuditPass {
		count, err := dao.VoiceLoverAudioAlbumDao.GetCountByAudioId(ctx, req.GetId())
		if err != nil {
			return err
		}
		if count > 0 {
			return consts.ERROR_AUDIO_COLLECT
		}
	}
	data := g.Map{}
	data["update_time"] = time.Now().Unix()
	data["audit_status"] = req.AuditStatus
	data["audit_reason"] = req.AuditReason
	data["op_uid"] = req.OpUid
	affect, err := dao.VoiceLoverAudioDao.UpdateAudioById(ctx, req.Id, data)
	if err != nil {
		return err
	}
	if affect > 0 {
		delete(data, "audit_reason")
		delete(data, "update_time")
		_ = es.EsClient(es.EsVpc).Update("voice_lover_audio", req.Id, data)
	}
	return nil
}

func (a *adminLogic) CreateAlbum(ctx context.Context, req *voice_lover.ReqCreateAlbum) (uint64, error) {
	id, err := dao.VoiceLoverAlbumDao.CreateAlbum(ctx, req.Name, req.Intro, req.Cover, req.OpUid)
	if err != nil {
		return 0, err
	}
	return uint64(id), nil
}

func (a *adminLogic) DelAlbum(ctx context.Context, req *voice_lover.ReqDelAlbum) error {
	info, err := dao.VoiceLoverAlbumDao.GetValidAlbumById(ctx, req.GetId())
	if err != nil {
		return err
	}
	if info == nil {
		return consts.ERROR_ALBUM_NOT_EXIST
	}
	count, err := dao.VoiceLoverAudioAlbumDao.GetCountByAlbumId(ctx, req.GetId())
	if err != nil {
		return err
	}
	if count > 0 {
		return consts.ERROR_ALBUM_HAS_AUDIO
	}
	count, err = dao.VoiceLoverAlbumSubjectDao.GetCountByAlbumId(ctx, req.GetId())
	if err != nil {
		return err
	}
	if count > 0 {
		return consts.ERROR_ALBUM_COLLECT
	}
	err = dao.VoiceLoverAlbumDao.DelAlbum(ctx, req.GetId(), req.GetOpUid())
	if err != nil {
		return err
	}
	return nil
}

func (a *adminLogic) UpdateAlbum(ctx context.Context, req *voice_lover.ReqUpdateAlbum) error {
	info, err := dao.VoiceLoverAlbumDao.GetValidAlbumById(ctx, req.GetId())
	if err != nil {
		return err
	}
	if info == nil {
		return consts.ERROR_ALBUM_NOT_EXIST
	}
	err = dao.VoiceLoverAlbumDao.UpdateAlbum(ctx, req.GetId(), req.GetName(), req.GetIntro(), req.GetCover(), req.GetOpUid())
	if err != nil {
		return err
	}
	return nil
}

func (a *adminLogic) GetAlbumDetail(ctx context.Context, req *voice_lover.ReqGetAlbumDetail) (map[uint64]*voice_lover.AlbumData, error) {
	albumIds := make([]uint64, 0)
	albumNames := make([]string, 0)
	for _, albumStr := range req.GetAlbumStr() {
		id, err := strconv.Atoi(albumStr)
		if err == nil {
			albumIds = append(albumIds, uint64(id))
		} else {
			albumNames = append(albumNames, albumStr)
		}
	}
	res := make(map[uint64]*voice_lover.AlbumData)
	if len(albumIds) > 0 {
		m, err := dao.VoiceLoverAlbumDao.GetValidAlbumByIds(ctx, albumIds)
		if err != nil {
			return nil, err
		}
		for id, info := range m {
			count, _ := dao.VoiceLoverAudioAlbumDao.GetCountByAlbumId(ctx, info.GetId())
			res[id] = &voice_lover.AlbumData{
				Id:         info.Id,
				Name:       info.Name,
				Intro:      info.Intro,
				Cover:      info.Cover,
				CreateTime: info.CreateTime,
				AudioCount: uint32(count),
				OpUid:      info.OpUid,
			}
		}
	}
	if len(albumNames) > 0 {
		m, err := dao.VoiceLoverAlbumDao.GetValidAlbumByNames(ctx, albumNames)
		if err != nil {
			return nil, err
		}
		for id, info := range m {
			count, _ := dao.VoiceLoverAudioAlbumDao.GetCountByAlbumId(ctx, info.GetId())
			res[id] = &voice_lover.AlbumData{
				Id:         info.Id,
				Name:       info.Name,
				Intro:      info.Intro,
				Cover:      info.Cover,
				CreateTime: info.CreateTime,
				AudioCount: uint32(count),
				OpUid:      info.OpUid,
			}
		}
	}
	return res, nil
}

func (a *adminLogic) GetAlbumList(ctx context.Context, req *voice_lover.ReqGetAlbumList) ([]*voice_lover.AlbumData, int32, error) {
	list, total, err := dao.VoiceLoverAlbumDao.GetValidAlbumList(ctx, req.StartTime, req.EndTime, req.Name, int(req.Page), int(req.Limit))
	if err != nil {
		return nil, 0, err
	}
	res := make([]*voice_lover.AlbumData, 0)
	for _, l := range list {
		count, _ := dao.VoiceLoverAudioAlbumDao.GetCountByAlbumId(ctx, l.GetId())
		res = append(res, &voice_lover.AlbumData{
			Id:         l.Id,
			Name:       l.Name,
			Intro:      l.Intro,
			Cover:      l.Cover,
			CreateTime: l.CreateTime,
			AudioCount: uint32(count),
			OpUid:      l.OpUid,
		})
	}
	return res, int32(total), nil
}

func (a *adminLogic) AudioCollect(ctx context.Context, req *voice_lover.ReqAudioCollect) error {
	audio, err := dao.VoiceLoverAudioDao.GetAudioDetailByAudioId(ctx, req.AudioId)
	if err != nil {
		return err
	}
	if audio == nil {
		return consts.ERROR_AUDIO_NOT_EXIST
	}
	if audio.AuditStatus != dao.AuditPass {
		return consts.ERROR_AUDIO_STATUS_INVALID
	}
	album, err := dao.VoiceLoverAlbumDao.GetValidAlbumById(ctx, req.GetAlbumId())
	if err != nil {
		return err
	}
	if album == nil {
		return consts.ERROR_ALBUM_NOT_EXIST
	}
	audioAlbum, err := dao.VoiceLoverAudioAlbumDao.GetAudioAlbumByAudioIdAlbumId(ctx, req.AudioId, req.AlbumId)
	if err != nil {
		return err
	}
	if req.Type == Collect {
		if audioAlbum != nil {
			return consts.ERROR_AUDIO_ALBUM_COLLECT
		}
		err := dao.VoiceLoverAudioAlbumDao.Create(ctx, req.AudioId, req.AlbumId)
		if err != nil {
			return err
		}
	}
	if req.Type == CollectRemove {
		if audioAlbum == nil {
			return consts.ERROR_AUDIO_ALBUM_COLLECT_REMOVE
		}
		err := dao.VoiceLoverAudioAlbumDao.Del(ctx, req.AudioId, req.AlbumId)
		if err != nil {
			return err
		}
	}
	albumIds, err := dao.VoiceLoverAudioAlbumDao.GetAlbumIdsByAudioId(ctx, req.AudioId)
	if err != nil {
		return err
	}
	data := g.Map{}
	if len(albumIds) > 0 {
		data["has_album"] = 1
	}
	data["albums"] = albumIds
	_ = es.EsClient(es.EsVpc).Update("voice_lover_audio", req.AudioId, data)
	return nil
}

func (a *adminLogic) CreateSubject(ctx context.Context, req *voice_lover.ReqCreateSubject) (uint64, error) {
	id, err := dao.VoiceLoverSubjectDao.CreateSubject(ctx, req.Name, req.OpUid)
	if err != nil {
		return 0, err
	}
	return uint64(id), err
}

func (a *adminLogic) UpdateSubject(ctx context.Context, req *voice_lover.ReqUpdateSubject) error {
	return dao.VoiceLoverSubjectDao.UpdateSubject(ctx, req.Id, req.Name, req.OpUid)
}

func (a *adminLogic) DelSubject(ctx context.Context, req *voice_lover.ReqDelSubject) error {
	count, err := dao.VoiceLoverAlbumSubjectDao.GetCountBySubjectId(ctx, req.Id)
	if err != nil {
		return err
	}
	if count > 0 {
		return consts.ERROR_SUBJECT_HAS_ALBUM
	}
	return dao.VoiceLoverSubjectDao.DelSubject(ctx, req.Id, req.OpUid)
}

func (a *adminLogic) GetSubjectDetail(ctx context.Context, req *voice_lover.ReqGetSubjectDetail) (map[uint64]*voice_lover.SubjectData, error) {
	data, err := dao.VoiceLoverSubjectDao.GetValidSubjectByIds(ctx, req.Ids)
	if err != nil {
		return nil, err
	}
	res := make(map[uint64]*voice_lover.SubjectData)
	for k, v := range data {
		res[k] = &voice_lover.SubjectData{
			Id:         v.Id,
			Name:       v.Name,
			CreateTime: v.CreateTime,
		}
	}
	return res, nil
}

func (a *adminLogic) GetSubjectList(ctx context.Context, req *voice_lover.ReqGetSubjectList) ([]*voice_lover.SubjectData, int32, error) {
	list, total, err := dao.VoiceLoverSubjectDao.GetValidSubjectListByName(ctx, req.StartTime, req.EndTime, req.Name, int(req.Page), int(req.Limit))
	if err != nil {
		return nil, 0, err
	}
	res := make([]*voice_lover.SubjectData, 0)
	for _, r := range list {
		count, _ := dao.VoiceLoverAlbumSubjectDao.GetCountBySubjectId(ctx, r.GetId())
		res = append(res, &voice_lover.SubjectData{
			Id:         r.Id,
			Name:       r.Name,
			CreateTime: r.CreateTime,
			AlbumCount: uint64(count),
		})
	}
	return res, int32(total), nil
}

func (a *adminLogic) AlbumCollect(ctx context.Context, req *voice_lover.ReqAlbumCollect) error {
	albumSubject, err := dao.VoiceLoverAlbumSubjectDao.GetAlbumSubjectByAIdAndSId(ctx, req.AlbumId, req.SubjectId)
	if err != nil {
		return err
	}
	if req.CollectType == Collect {
		if albumSubject != nil {
			return consts.ERROR_ALBUM_SUBJECT_COLLECT
		}
	}
	if req.CollectType == CollectRemove {
		if albumSubject == nil {
			return consts.ERROR_ALBUM_SUBJECT_COLLECT_REMOVE
		}
	}
	err = dao.VoiceLoverAlbumSubjectDao.Create(ctx, req.AlbumId, req.SubjectId)
	return err
}

func (a *adminLogic) GetAlbumCollect(ctx context.Context, req *voice_lover.ReqGetAlbumCollect, reply *voice_lover.ResGetAlbumCollect) error {
	albumStr := req.AlbumStr
	subjectStr := req.SubjectStr
	var err error
	var albumId = -1
	if len(albumStr) > 0 {
		albumId, err = strconv.Atoi(albumStr)
		if err != nil {
			album, err := dao.VoiceLoverAlbumDao.GetValidAlbumByName(ctx, albumStr)
			if err != nil {
				return err
			}
			albumId = int(album.GetId())
		}
	}
	var subjectId = -1
	if len(subjectStr) > 0 {
		subjectId, err = strconv.Atoi(subjectStr)
		if err != nil {
			subject, err := dao.VoiceLoverSubjectDao.GetValidSubjectByName(ctx, subjectStr)
			if err != nil {
				return err
			}
			subjectId = int(subject.GetId())
		}
	}
	list, total, err := dao.VoiceLoverAlbumSubjectDao.GetAlbumCollect(ctx, uint64(albumId), uint64(subjectId), req.Page, req.Limit)
	if err != nil {
		return err
	}
	albumIds := make([]uint64, 0)
	subjectIds := make([]uint64, 0)
	for _, l := range list {
		albumIds = append(albumIds, l.AlbumId)
		subjectIds = append(subjectIds, l.SubjectId)
	}
	albums, err := dao.VoiceLoverAlbumDao.GetValidAlbumByIds(ctx, albumIds)
	if err != nil {
		return err
	}
	subjects, err := dao.VoiceLoverSubjectDao.GetValidSubjectByIds(ctx, subjectIds)
	if err != nil {
		return err
	}
	res := make([]*voice_lover.AlbumCollectData, 0)
	for _, l := range list {
		res = append(res, &voice_lover.AlbumCollectData{
			Id:          l.Id,
			AlbumName:   albums[l.AlbumId].GetName(),
			SubjectName: subjects[l.SubjectId].GetName(),
			AlbumId:     l.AlbumId,
			SubjectId:   l.SubjectId,
		})
	}
	reply.AlbumCollects = res
	reply.Total = int32(total)
	return nil
}

func (a *adminLogic) AlbumChoice(ctx context.Context, req *voice_lover.ReqAlbumChoice) error {
	album, err := dao.VoiceLoverAlbumDao.GetValidAlbumById(ctx, req.Id)
	if err != nil {
		return err
	}
	if album == nil {
		return consts.ERROR_ALBUM_NOT_EXIST
	}
	return dao.VoiceLoverAlbumDao.AlbumChoice(ctx, req.Id, req.Type)
}

func (a *adminLogic) GetAlbumChoice(ctx context.Context, req *voice_lover.ReqGetAlbumChoice) ([]*functor.EntityVoiceLoverAlbum, error) {
	return dao.VoiceLoverAlbumDao.GetAlbumChoice(ctx)
}
