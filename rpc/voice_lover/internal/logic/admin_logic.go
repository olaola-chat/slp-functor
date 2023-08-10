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
		return consts.ERROR_PARAM
	}
	data := g.Map{
		"update_time": time.Now().Unix(),
	}
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
		return nil
	}
	count, err := dao.VoiceLoverAudioAlbumDao.GetCountByAlbumId(ctx, req.GetId())
	if err != nil {
		return err
	}
	if count > 0 {
		return consts.ERROR_PARAM
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
		return nil
	}
	count, err := dao.VoiceLoverAudioAlbumDao.GetCountByAlbumId(ctx, req.GetId())
	if err != nil {
		return err
	}
	if count > 0 {
		return consts.ERROR_PARAM
	}
	err = dao.VoiceLoverAlbumDao.UpdateAlbum(ctx, req.GetId(), req.GetName(), req.GetIntro(), req.GetCover(), req.GetOpUid())
	if err != nil {
		return err
	}
	return nil
}

func (a *adminLogic) GetAlbumDetail(ctx context.Context, req *voice_lover.ReqGetAlbumDetail) (*voice_lover.AlbumData, error) {
	albumStr := req.GetAlbumStr()
	albumId, err := strconv.Atoi(albumStr)
	var info *functor.EntityVoiceLoverAlbum
	if err == nil {
		info, err = dao.VoiceLoverAlbumDao.GetValidAlbumById(ctx, uint64(albumId))
	} else {
		info, err = dao.VoiceLoverAlbumDao.GetValidAlbumByName(ctx, albumStr)
	}
	if err != nil {
		return nil, err
	}
	if info == nil {
		return nil, nil
	}
	count, err := dao.VoiceLoverAudioAlbumDao.GetCountByAlbumId(ctx, info.GetId())
	if err != nil {
		return nil, err
	}
	album := &voice_lover.AlbumData{
		Id:         info.Id,
		Name:       info.Name,
		Intro:      info.Intro,
		Cover:      info.Cover,
		CreateTime: info.CreateTime,
		AudioCount: uint32(count),
		OpUid:      info.OpUid,
	}
	return album, nil
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
