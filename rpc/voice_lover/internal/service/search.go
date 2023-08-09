package service

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/util/gconv"
	"github.com/olaola-chat/rbp-library/es"
	"github.com/olaola-chat/rbp-proto/gen_pb/db/functor"
)

type VoiceLoverService struct {

}

func (s *VoiceLoverService) Put(data *functor.EntityVoiceLoverAudio) error {
	esClient := es.EsClient(es.EsVpc)
	labelsSlice := make([]string, 0)
	for _, l := range strings.Split(data.Labels, ",") {
		if len(l) == 0 {
			continue
		}
		labelsSlice = append(labelsSlice, l)
	}
	esModel := &VoiceLoverEsModel{
		Id:          data.Id,
		PubUid:      data.PubUid,
		Title:       data.Title,
		Cover:       data.Cover,
		Desc:        data.Desc,
		CreateTime:  data.CreateTime,
		Labels:      labelsSlice,
		Source:      int32(data.From),
		AuditStatus: int32(data.AuditStatus),
		Albums:      []uint64{},
		HasAlbum:    0,
		OpUid: data.OpUid,
	}
	err := esClient.Put("voice_lover_audio", gconv.String(data.Id), gconv.Map(esModel))
	if err != nil {
		g.Log().Errorf("VoiceLoverAudio Put es error, err = %v", err)
	}
	return err
}

func (s *VoiceLoverService) Search(startTime, endTime uint64, source int32, pubUids []uint64, label string, auditStatus int32, audioId uint64, albumId uint64, hasAlbum bool, pageNum int32, pageSize int32) ([]*VoiceLoverEsModel, int32, error) {
	client := es.EsClient(es.EsVpc)
	if endTime == 0 {
		endTime = uint64(time.Now().Unix())
	}
	from := (pageNum - 1) * pageSize
	must := make([]map[string]interface{}, 0)
	must = append(must, map[string]interface{}{
		"range": map[string]interface{}{
			"create_time": map[string]interface{}{
				"gte": startTime,
				"lte": endTime,
			},
		},
	})
	if source > 0 {
		must = append(must, map[string]interface{}{
			"term": map[string]interface{}{
				"source": source,
			},
		})
	}
	if len(pubUids) > 0 {
		must = append(must, map[string]interface{}{
			"terms": map[string]interface{}{
				"pub_uid": pubUids,
			},
		})
	}
	if len(label) > 0 {
		must = append(must, map[string]interface{}{
			"term": map[string]interface{}{
				"labels": label,
			},
		})
	}
	if auditStatus >= 0 {
		must = append(must, map[string]interface{}{
			"term": map[string]interface{}{
				"audit_status": auditStatus,
			},
		})
	}
	if audioId > 0 {
		must = append(must, map[string]interface{}{
			"term": map[string]interface{}{
				"id": audioId,
			},
		})
	}
	if albumId > 0 {
		must = append(must, map[string]interface{}{
			"term": map[string]interface{}{
				"Albums": albumId,
			},
		})
	}
	body := map[string]interface{}{
		"from":    from,
		"size":    pageSize,
		//"_source": []string{"id", "name", "rcmd_name", "icon", "rcmd_icon", "user_num"},
		"sort": map[string]interface{}{
			"id":  map[string]string{"order": "desc"},
		},
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": must,
			},
		},
	}
	bodyJson, _ := json.Marshal(body)
	resp, err := client.Search("voice_lover_audio", bodyJson)
	if err != nil {
		return nil, 0, err
	}
	res := make([]*VoiceLoverEsModel, 0)
	for _, hit := range resp.Hits.Hits {
		source := hit.Source
		m := &VoiceLoverEsModel{}
		_ = gconv.Struct(source, m)
		res = append(res, m)
	}
	return res, int32(resp.Hits.Total), nil
}

type VoiceLoverEsModel struct {
	Id uint64 `json:"id"`
	PubUid uint64 `json:"pub_uid"`
	Title string `json:"title"`
	Cover string `json:"cover"`
	Desc string `json:"desc"`
	CreateTime uint64 `json:"create_time"`
	Labels []string `json:"labels"`
	Source int32 `json:"source"`
	AuditStatus int32 `json:"audit_status"`
	Albums []uint64 `json:"albums"`
	HasAlbum int32 `json:"has_album"`
	OpUid  uint64 `json:"op_uid"`
}