package logic

import (
	"context"
	"strings"
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

func (m *mainLogic) Post(ctx context.Context, req *vl_pb.ReqVoiceLoverPost, reply *vl_pb.ResVoiceLoverBase) error {
	g.Log().Infof("VoiceLoverPost req = %v", req)
	now := uint64(time.Now().Unix())
	err := functor.VoiceLoverAudio.DB.Transaction(func(tx *gdb.TX) error {
		data := &functor2.EntityVoiceLoverAudio{
			Desc:       req.Desc,
			Resource:   req.Resource,
			Cover:      req.Cover,
			From:       uint32(req.Source),
			PubUid:     req.Uid,
			CreateTime: now,
			UpdateTime: now,
			Title:      req.Title,
			Labels:     req.Labels,
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

func buildAudioEsModel(data *functor2.EntityVoiceLoverAudio) *voice_lover2.VoiceLoverEsModel {
	labelsSlice := make([]string, 0)
	for _, l := range strings.Split(data.Labels, ",") {
		if len(l) == 0 {
			continue
		}
		labelsSlice = append(labelsSlice, l)
	}
	esModel := &voice_lover2.VoiceLoverEsModel{
		Id:          data.Id,
		PubUid:      data.PubUid,
		Title:       data.Title,
		Cover:       data.Cover,
		Desc:        data.Desc,
		CreateTime:  data.CreateTime,
		Labels:      labelsSlice,
		Source:      data.From,
		AuditStatus: data.AuditStatus,
		Albums:      []uint64{},
		HasAlbum:    0,
		OpUid:       data.OpUid,
	}
	return esModel
}

func (m *mainLogic) GetRecAlbums(ctx context.Context, req *vl_pb.ReqGetRecAlbums, reply *vl_pb.ResGetRecAlbums) error {
	reply.Albums = make([]*vl_pb.AlbumData, 0)
	list, _ := dao.VoiceLoverAlbumDao.GetAlbumListByChoice(ctx, dao.ChoiceRec, 0, 3)
	for _, v := range list {
		reply.Albums = append(reply.Albums, &vl_pb.AlbumData{
			Id:         v.Id,
			Name:       v.Name,
			Intro:      v.Intro,
			Cover:      v.Cover,
			CreateTime: v.CreateTime,
		})
	}
	return nil
}
