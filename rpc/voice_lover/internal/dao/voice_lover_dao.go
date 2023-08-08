package dao

import (
	"context"
	"github.com/gogf/gf/database/gdb"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/util/gconv"
	"github.com/olaola-chat/rbp-proto/dao/functor"
	functor2 "github.com/olaola-chat/rbp-proto/gen_pb/db/functor"
	"github.com/olaola-chat/rbp-proto/gen_pb/rpc/voice_lover"
	"strings"
	"time"
)

type voiceLoverDao struct {
}

const (
	None = iota
	Dub
	Content
	Post
	Cover
)

var VoiceLoverDao = &voiceLoverDao{}

func (v *voiceLoverDao) Post(ctx context.Context, req *voice_lover.ReqVoiceLoverPost) error {
	now := uint64(time.Now().Unix())
	err := functor.VoiceLoverAudio.DB.Transaction(func(tx *gdb.TX) error {
		last, err := functor.VoiceLoverAudio.TX(tx).Insert(&functor2.EntityVoiceLoverAudio{
			Desc:       req.Desc,
			Resource:   req.Resource,
			Cover:      req.Cover,
			From:       uint32(req.Source),
			PubUid:     req.Uid,
			CreateTime: now,
			UpdateTime: now,
			Title:      req.Title,
		})
		if err != nil {
			return err
		}
		lastId, _ := last.LastInsertId()
		audioId := uint64(lastId)
		editDatas := make([]*functor2.EntityVoiceLoverAudioPartner, 0)
		editDubs := strings.Split(req.EditDub, ",")
		for _, editDub := range editDubs {
			uid := gconv.Uint64(editDub)
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
		audioLabels := strings.Split(req.Labels, ",")
		labelDatas := make([]*functor2.EntityVoiceLoverAudioLabel, 0)
		if len(audioLabels) > 0 {
			for _, label := range audioLabels {
				labelDatas = append(labelDatas, &functor2.EntityVoiceLoverAudioLabel{
					AudioId:    audioId,
					Label:      label,
					CreateTime: now,
					UpdateTime: now,
				})
			}
		}
		if len(labelDatas) > 0 {
			_, err := functor.VoiceLoverAudioLabel.TX(tx).Insert(labelDatas)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		g.Log().Errorf("VoiceLover post error, err = %v", err)
	}
	return err
}
