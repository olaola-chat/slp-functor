package service

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/util/gconv"
	"github.com/olaola-chat/rbp-library/es"
	xianshi2 "github.com/olaola-chat/rbp-proto/dao/xianshi"
	"github.com/olaola-chat/rbp-proto/gen_pb/db/xianshi"
	user2 "github.com/olaola-chat/rbp-proto/gen_pb/rpc/user"
	voice_lover3 "github.com/olaola-chat/rbp-proto/gen_pb/rpc/voice_lover"
	"github.com/olaola-chat/rbp-proto/rpcclient/user"
	voice_lover2 "github.com/olaola-chat/rbp-proto/rpcclient/voice_lover"

	vl_pb "github.com/olaola-chat/rbp-proto/gen_pb/rpc/voice_lover"

	"github.com/olaola-chat/rbp-functor/app/model/voice_lover"
	"github.com/olaola-chat/rbp-functor/app/utils"
)

type VoiceLoverAudioSearchQuery struct {
	StartTime   uint64
	EndTime     uint64
	Source      int32
	PubUIds     []uint64
	Label       string
	AuditStatus int32
	AudioStr    string
	AlbumIds    []uint64
	HasAlbum    int32
	PageNum     int32
	PageSize    int32
	Order       map[string]string
}

type voiceLoverService struct {
}

var VoiceLoverService = &voiceLoverService{}

func (s *voiceLoverService) BuildAudioCollectSearchQuery(ctx context.Context, req *vl_pb.ReqAdminAudioCollectList) *VoiceLoverAudioSearchQuery {
	g.Log().Infof("BuildAudioCollectSearchQuery req = %v", *req)
	pubUids := make([]uint64, 0)
	if len(req.UserStr) != 0 {
		uid, err := strconv.Atoi(req.UserStr)
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
				g.Log().Warningf("BuildAudioCollectSearchQuery search by name error, err = %v")
			}
			g.Log().Infof("BuildAudioCollectSearchQuery searchbyname name = %s data = %v", req.UserStr, res.Data)
			if err == nil {
				for _, r := range res.Data {
					pubUids = append(pubUids, uint64(r))
				}
				if len(pubUids) == 0 {
					pubUids = append(pubUids, 0)
				}
			}
		}
	}
	var albumIds []uint64
	if len(req.AlbumStr) > 0 {
		id, err := strconv.Atoi(req.AlbumStr)
		if err == nil {
			albumIds = append(albumIds, uint64(id))
		} else {
			albumReply, err := voice_lover2.VoiceLoverAdmin.GetAlbumDetail(ctx, &voice_lover3.ReqGetAlbumDetail{AlbumStr: []string{req.AlbumStr}})
			if err == nil && len(albumReply.Albums) > 0 {
				for albumId := range albumReply.Albums {
					albumIds = append(albumIds, albumId)
				}
			}
			if len(albumIds) == 0 {
				albumIds = append(albumIds, 0)
			}
		}
	}
	q := &VoiceLoverAudioSearchQuery{
		PubUIds:     pubUids,
		Label:       req.Label,
		AuditStatus: 1,
		PageNum:     int32(req.Page),
		PageSize:    int32(req.Limit),
		Order:       map[string]string{},
		AlbumIds:    albumIds,
		AudioStr:    req.AudioStr,
		HasAlbum:    req.CollectStatus,
	}
	_ = json.Unmarshal([]byte(req.Order), &q.Order)
	g.Log().Infof("BuildAudioCollectSearchQuery query = %v", *q)
	return q
}

func (s *voiceLoverService) BuildAudioSearchQuery(ctx context.Context, req *vl_pb.ReqAdminAudioList) *VoiceLoverAudioSearchQuery {
	g.Log().Infof("BuildAudioSearchQuery req = %v", *req)
	pubUids := make([]uint64, 0)
	if len(req.UserStr) != 0 {
		uid, err := strconv.Atoi(req.UserStr)
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
			g.Log().Infof("BuildAudioSearchQuery searchbyname name = %s data = %v", req.UserStr, res.Data)
			if err == nil {
				for _, r := range res.Data {
					pubUids = append(pubUids, uint64(r))
				}
				if len(pubUids) == 0 {
					pubUids = append(pubUids, 0)
				}
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
		HasAlbum:    -1,
	}
	_ = json.Unmarshal([]byte(req.Order), &q.Order)
	g.Log().Infof("BuildAudioSearchQuery query = %v", *q)
	return q
}

func (s *voiceLoverService) BuildVoiceLoverAudioPb(models []*voice_lover.VoiceLoverAudioEsModel) []*vl_pb.AdminAudio {
	data := make([]*vl_pb.AdminAudio, 0)
	uidMap := make(map[uint32]struct{})
	for _, model := range models {
		uidMap[model.PubUid] = struct{}{}
	}
	uids := make([]uint32, 0)
	for uid := range uidMap {
		uids = append(uids, uid)
	}
	userReply, _ := user.UserProfile.Mget(context.Background(), &user2.ReqUserProfiles{
		Uids:   uids,
		Fields: []string{"uid", "icon", "name"},
	})
	userMap := make(map[uint32]*xianshi.EntityXsUserProfile)
	for _, u := range userReply.Data {
		userMap[u.Uid] = u
	}
	userBrokers, _ := s.GetUserBroker(uids)
	for _, model := range models {
		covers := make([]string, 0)
		for _, c := range strings.Split(model.Cover, ",") {
			if len(c) == 0 {
				continue
			}
			covers = append(covers, c)
		}
		data = append(data, &vl_pb.AdminAudio{
			Id:          model.Id,
			CreateTime:  model.CreateTime,
			PubUid:      model.PubUid,
			PubUserName: userMap[model.PubUid].GetName(),
			Broker:      userBrokers[model.PubUid].GetBname(),
			Resource:    model.Resource,
			Covers:      covers,
			Source:      model.Source,
			Desc:        model.Desc,
			Labels:      model.Labels,
			AuditStatus: model.AuditStatus,
			OpUid:       model.OpUid,
			Title:       model.Title,
		})
	}
	return data
}

func (s *voiceLoverService) GetUserBroker(uids []uint32) (map[uint32]*xianshi.EntityXsBroker, error) {
	userBrokers, err := xianshi2.XsBrokerUser.Where("uid IN (?)", uids).Where("deleted = 0 and state = 1").FindAll()
	if err != nil {
		return map[uint32]*xianshi.EntityXsBroker{}, err
	}
	uidBidMap := make(map[int32]int32)
	bids := make([]int32, 0)
	for _, userBroker := range userBrokers {
		uidBidMap[userBroker.Uid] = userBroker.Bid
		bids = append(bids, userBroker.Bid)
	}
	bids = utils.DistinctInt32Slice(bids)
	brokerMap := make(map[int32]*xianshi.EntityXsBroker)
	brokers, err := xianshi2.XsBroker.Where("bid IN (?)", bids).Where("deleted = 0").FindAll()
	if err != nil {
		return map[uint32]*xianshi.EntityXsBroker{}, err
	}
	for _, broker := range brokers {
		brokerMap[broker.Bid] = broker
	}
	res := make(map[uint32]*xianshi.EntityXsBroker)
	for _, uid := range uids {
		res[uid] = brokerMap[uidBidMap[int32(uid)]]
	}
	return res, nil
}

func (s *voiceLoverService) BuildVoiceLoverAudioCollectPb(models []*voice_lover.VoiceLoverAudioEsModel) []*vl_pb.AdminAudioCollect {
	data := make([]*vl_pb.AdminAudioCollect, 0)
	uidMap := make(map[uint32]struct{})
	for _, model := range models {
		uidMap[model.PubUid] = struct{}{}
	}
	uids := make([]uint32, 0)
	for uid := range uidMap {
		uids = append(uids, uid)
	}
	userReply, _ := user.UserProfile.Mget(context.Background(), &user2.ReqUserProfiles{
		Uids:   uids,
		Fields: []string{"uid", "icon", "name"},
	})
	userMap := make(map[uint32]*xianshi.EntityXsUserProfile)
	for _, u := range userReply.Data {
		userMap[u.Uid] = u
	}
	albumIds := make([]uint64, 0)
	for _, model := range models {
		albumIds = append(albumIds, model.Albums...)
		collects := make([]*vl_pb.AdminAudioCollectAlbum, 0)
		for _, albumId := range model.Albums {
			collects = append(collects, &vl_pb.AdminAudioCollectAlbum{
				Id: albumId,
			})
		}
		data = append(data, &vl_pb.AdminAudioCollect{
			Id:          model.Id,
			CreateTime:  model.CreateTime,
			PubUid:      model.PubUid,
			PubUserName: userMap[model.PubUid].GetName(),
			Labels:      model.Labels,
			Title:       model.Title,
			Collects:    collects,
		})
	}
	albumStrs := utils.Uint64SliceToStrSlice(utils.DistinctUint64Slice(albumIds))
	if len(albumStrs) > 0 {
		albumReply, err := voice_lover2.VoiceLoverAdmin.GetAlbumDetail(context.Background(), &voice_lover3.ReqGetAlbumDetail{AlbumStr: albumStrs})
		if err != nil {
			g.Log().Errorf("[BuildVoiceLoverAudioCollectPb] GetAlbumDetail error, err = %v")
		}
		if err == nil {
			for _, d := range data {
				for _, c := range d.Collects {
					c.Name = albumReply.Albums[c.GetId()].GetName()
				}
			}
		}
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
			"match": map[string]interface{}{
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
	if len(query.AudioStr) > 0 {
		audioId, err := strconv.Atoi(query.AudioStr)
		if err == nil {
			must = append(must, map[string]interface{}{
				"term": map[string]interface{}{
					"id": audioId,
				},
			})
		} else {
			must = append(must, map[string]interface{}{
				"match": map[string]interface{}{
					"title": query.AudioStr,
				},
			})
		}
	}
	if len(query.AlbumIds) > 0 {
		must = append(must, map[string]interface{}{
			"terms": map[string]interface{}{
				"albums": query.AlbumIds,
			},
		})
	}
	if query.HasAlbum >= 0 {
		must = append(must, map[string]interface{}{
			"term": map[string]interface{}{
				"has_album": query.HasAlbum,
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
