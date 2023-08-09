package voice_lover

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/util/gconv"
	"github.com/olaola-chat/rbp-library/es"
	"github.com/olaola-chat/rbp-proto/gen_pb/db/xianshi"
	user2 "github.com/olaola-chat/rbp-proto/gen_pb/rpc/user"
	"github.com/olaola-chat/rbp-proto/rpcclient/user"

	"github.com/olaola-chat/rbp-functor/app/model/voice_lover"
	"github.com/olaola-chat/rbp-functor/app/pb"
	"github.com/olaola-chat/rbp-functor/app/query"
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

func (s *voiceLoverService) GetAudioList(ctx context.Context, req *query.ReqAdminVoiceLoverAudioList) ([]*pb.AdminVoiceLoverAudio, int32, error) {
	q := s.BuildAudioSearchQuery(ctx, req)
	res, total, err := s.SearchAudio(ctx, q)
	if err != nil {
		return nil, 0, err
	}
	g.Log().Infof("GetAudioList total = %d audios = %v", total, res)
	return s.BuildVoiceLoverAudioPb(res), total, nil
}

func (s *voiceLoverService) BuildAudioSearchQuery(ctx context.Context, req *query.ReqAdminVoiceLoverAudioList) *VoiceLoverAudioSearchQuery {
	g.Log().Infof("BuildAudioSearchQuery req = %v", *req)
	uid, err := strconv.Atoi(req.UserStr)
	pubUids := make([]uint64, 0)
	if err == nil {
		pubUids = append(pubUids, uint64(uid))
	}
	if err != nil {
		res, err := user.UserProfile.SearchByName(ctx, &user2.ReqUserSearchName{
			Keyword:       req.UserStr,
			Limit:         10,
			SearcherLevel: 1,
		})
		if err != nil {
			g.Log().Warningf("BuildAudioSearchQuery search by name error, err = %v")
		}
		if err == nil {
			for _, r := range res.Data {
				pubUids = append(pubUids, uint64(r))
			}
		}
	}
	q := &VoiceLoverAudioSearchQuery{
		StartTime:   req.StartTime,
		EndTime:     req.EndTime,
		Source:      req.Source,
		PubUIds:     pubUids,
		Label:       req.Label,
		AuditStatus: req.AuditStatus,
		PageNum:     int32(req.Page),
		PageSize:    int32(req.Limit),
		Order:       map[string]string{},
	}
	_ = json.Unmarshal([]byte(req.Order), &q.Order)
	g.Log().Infof("BuildAudioSearchQuery query = %v", *q)
	return q
}

func (s *voiceLoverService) BuildVoiceLoverAudioPb(models []*voice_lover.VoiceLoverAudioEsModel) []*pb.AdminVoiceLoverAudio {
	data := make([]*pb.AdminVoiceLoverAudio, 0)
	uidMap := make(map[uint64]struct{})
	for _, model := range models {
		uidMap[model.PubUid] = struct{}{}
	}
	uids := make([]uint32, 0)
	for uid := range uidMap {
		uids = append(uids, uint32(uid))
	}
	userReply, _ := user.UserProfile.Mget(context.Background(), &user2.ReqUserProfiles{
		Uids:   uids,
		Fields: []string{"uid", "icon", "name"},
	})
	userMap := make(map[uint64]*xianshi.EntityXsUserProfile)
	for _, u := range userReply.Data {
		userMap[uint64(u.Uid)] = u
	}
	for _, model := range models {
		covers := make([]string, 0)
		for _, c := range strings.Split(model.Cover, ",") {
			if len(c) == 0 {
				continue
			}
			covers = append(covers, c)
		}
		data = append(data, &pb.AdminVoiceLoverAudio{
			Id:          model.Id,
			CreateTime:  model.CreateTime,
			PubUid:      model.PubUid,
			PubUserName: userMap[model.PubUid].GetName(),
			Broker:      "",
			Resource:    model.Resource,
			Covers:      covers,
			Source:      model.Source,
			Desc:        model.Desc,
			Labels:      model.Labels,
			AuditStatus: model.AuditStatus,
			OpUid:       model.OpUid,
		})
	}

	return data
}

func (s *voiceLoverService) SearchAudio(ctx context.Context, query *VoiceLoverAudioSearchQuery) ([]*voice_lover.VoiceLoverAudioEsModel, int32, error) {
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
	g.Log().Infof("VoiceLoverSearch req = %v", string(bodyJson))
	resp, err := client.Search("voice_lover_audio", bodyJson)
	if err != nil {
		return nil, 0, err
	}
	res := make([]*voice_lover.VoiceLoverAudioEsModel, 0)
	for _, hit := range resp.Hits.Hits {
		source := hit.Source
		m := &voice_lover.VoiceLoverAudioEsModel{}
		_ = gconv.Struct(source, m)
		res = append(res, m)
	}
	return res, int32(resp.Hits.Total), nil
}
