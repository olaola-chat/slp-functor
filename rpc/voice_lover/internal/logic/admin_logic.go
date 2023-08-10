package logic

import (
	"context"
	"strings"
	"time"

	"github.com/gogf/gf/frame/g"
	"github.com/olaola-chat/rbp-library/es"
	"github.com/olaola-chat/rbp-proto/gen_pb/rpc/voice_lover"

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
	if len(req.Labels) > 0 {
		data["labels"] = req.Labels
	}
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
