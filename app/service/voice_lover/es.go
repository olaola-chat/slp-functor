package voice_lover

import (
	"encoding/json"
	"time"

	"github.com/gogf/gf/util/gconv"
	"github.com/olaola-chat/rbp-library/es"

	"github.com/olaola-chat/rbp-functor/app/model/voice_lover"
)

type VoiceLoverAudioSearchQuery struct {
	StartTime   uint64
	EndTime     uint64
	Source      int32
	PubUIds     []uint64
	Label       string
	AuditStatus int32
	AudioId     uint64
	AlbumId     uint64
	HasAlbum    int32
	PageNum     int32
	PageSize    int32
	Order       map[string]string
}

func (s *voiceLoverService) Search(query *VoiceLoverAudioSearchQuery) ([]*voice_lover.VoiceLoverEsModel, int32, error) {
	client := es.EsClient(es.EsVpc)
	if query.EndTime == 0 {
		query.EndTime = uint64(time.Now().Unix())
	}
	from := (query.PageNum - 1) * query.PageSize
	must := make([]map[string]interface{}, 0)
	must = append(must, map[string]interface{}{
		"range": map[string]interface{}{
			"create_time": map[string]interface{}{
				"gte": query.StartTime,
				"lte": query.EndTime,
			},
		},
	})
	if query.Source > 0 {
		must = append(must, map[string]interface{}{
			"term": map[string]interface{}{
				"source": query.Source,
			},
		})
	}
	if len(query.PubUIds) > 0 {
		must = append(must, map[string]interface{}{
			"terms": map[string]interface{}{
				"pub_uid": query.PubUIds,
			},
		})
	}
	if len(query.Label) > 0 {
		must = append(must, map[string]interface{}{
			"term": map[string]interface{}{
				"labels": query.Label,
			},
		})
	}
	if query.AuditStatus >= 0 {
		must = append(must, map[string]interface{}{
			"term": map[string]interface{}{
				"audit_status": query.AuditStatus,
			},
		})
	}
	if query.AudioId > 0 {
		must = append(must, map[string]interface{}{
			"term": map[string]interface{}{
				"id": query.AudioId,
			},
		})
	}
	if query.AlbumId > 0 {
		must = append(must, map[string]interface{}{
			"term": map[string]interface{}{
				"albums": query.AlbumId,
			},
		})
	}
	sort := make(map[string]interface{})
	if len(query.Order) == 0 {
		query.Order["id"] = "desc"
	}
	for field, s := range query.Order {
		sort[field] = map[string]string{
			"order": s,
		}
	}
	body := map[string]interface{}{
		"from": from,
		"size": query.PageSize,
		//"_source": []string{"id", "name", "rcmd_name", "icon", "rcmd_icon", "user_num"},
		"sort": sort,
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
	res := make([]*voice_lover.VoiceLoverEsModel, 0)
	for _, hit := range resp.Hits.Hits {
		source := hit.Source
		m := &voice_lover.VoiceLoverEsModel{}
		_ = gconv.Struct(source, m)
		res = append(res, m)
	}
	return res, int32(resp.Hits.Total), nil
}
