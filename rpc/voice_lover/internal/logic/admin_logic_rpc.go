package logic

import (
	"context"
	"fmt"
	"sort"

	"github.com/gogf/gf/frame/g"
	"github.com/olaola-chat/rbp-proto/gen_pb/db/xianshi"
	user2 "github.com/olaola-chat/rbp-proto/gen_pb/rpc/user"
	vl_pb "github.com/olaola-chat/rbp-proto/gen_pb/rpc/voice_lover"
	"github.com/olaola-chat/rbp-proto/rpcclient/user"

	"github.com/olaola-chat/rbp-functor/app/consts"
	"github.com/olaola-chat/rbp-functor/app/utils"
	"github.com/olaola-chat/rbp-functor/rpc/voice_lover/internal/service"
)

func (v *adminLogic) AdminAudioList(ctx context.Context, req *vl_pb.ReqAdminAudioList, reply *vl_pb.ResAdminAudioList) error {
	q := service.VoiceLoverService.BuildAudioSearchQuery(ctx, req)
	res, total, err := service.VoiceLoverService.SearchAudio(ctx, q)
	if err != nil {
		g.Log().Errorf("AdminAudioList error, err = %v", err)
		return err
	}
	reply.Audios = service.VoiceLoverService.BuildVoiceLoverAudioPb(res)
	reply.Total = total
	return nil
}

func (v *adminLogic) AdminAudioDetail(ctx context.Context, req *vl_pb.ReqAdminAudioDetail, reply *vl_pb.ResAdminAudioDetail) error {
	audio, err := v.GetAudioDetail(ctx, req.GetId())
	if err != nil {
		return err
	}
	if audio == nil {
		return nil
	}
	uids := make([]uint32, 0)
	for _, r := range audio.EditContents {
		uids = append(uids, r.Uid)
	}
	for _, r := range audio.EditDubs {
		uids = append(uids, r.Uid)
	}
	for _, r := range audio.EditContents {
		uids = append(uids, r.Uid)
	}
	for _, r := range audio.EditPosts {
		uids = append(uids, r.Uid)
	}
	uids = utils.DistinctUint32Slice(uids)
	userMap := make(map[uint32]*xianshi.EntityXsUserProfile)
	userReply, err := user.UserProfile.Mget(ctx, &user2.ReqUserProfiles{
		Uids:   uids,
		Fields: []string{"name", "uid", "icon"},
	})
	for _, u := range userReply.Data {
		userMap[u.Uid] = u
	}
	info := &vl_pb.AdminAudio{
		Id:           audio.Id,
		CreateTime:   audio.CreateTime,
		PubUid:       audio.Uid,
		Resource:     audio.Resource,
		Covers:       audio.Covers,
		Source:       audio.Source,
		Desc:         audio.Desc,
		Labels:       audio.Labels,
		AuditStatus:  audio.AuditStatus,
		OpUid:        audio.OpUid,
		EditDubs:     buildAudioEdit(audio.EditDubs, userMap),
		EditContents: buildAudioEdit(audio.EditContents, userMap),
		EditPosts:    buildAudioEdit(audio.EditPosts, userMap),
		EditCovers:   buildAudioEdit(audio.EditCovers, userMap),
	}
	reply.Audio = info
	return nil
}

func buildAudioEdit(edits []*vl_pb.AudioEditData, userMap map[uint32]*xianshi.EntityXsUserProfile) []*vl_pb.AdminAudioEdit {
	editDubs := make([]*vl_pb.AdminAudioEdit, 0)
	for _, e := range edits {
		editDubs = append(editDubs, &vl_pb.AdminAudioEdit{
			Uid:    e.Uid,
			Name:   userMap[e.Uid].GetName(),
			Avatar: userMap[e.Uid].GetIcon(),
		})
	}
	return editDubs
}

func (v *adminLogic) AdminAudioUpdate(ctx context.Context, req *vl_pb.ReqAdminAudioUpdate, reply *vl_pb.ResAdminAudioUpdate) error {
	return v.UpdateAudio(ctx, &vl_pb.ReqUpdateAudio{
		Id:     req.Id,
		Title:  req.Title,
		Desc:   req.Title,
		Labels: req.Labels,
		OpUid:  req.OpUid,
	})
}

func (v *adminLogic) AdminAudioAudit(ctx context.Context, req *vl_pb.ReqAdminAudioAudit, reply *vl_pb.ResAdminAudioAudit) error {
	return v.AuditAudio(ctx, &vl_pb.ReqAuditAudio{
		Id:          req.Id,
		AuditStatus: req.AuditStatus,
		AuditReason: req.AuditReason,
		OpUid:       req.OpUid,
	})
}

func (v *adminLogic) AdminAudioAuditReason(ctx context.Context, reply *vl_pb.ResAdminAudioAuditReason) error {
	for id, reason := range consts.AuditAudioReasonMap {
		reply.Reasons = append(reply.Reasons, &vl_pb.AdminAudioAuditReason{
			Id:     id,
			Reason: reason,
		})
	}
	sort.Slice(reply.Reasons, func(i, j int) bool {
		return reply.Reasons[i].Id < reply.Reasons[j].Id
	})
	return nil
}

func (v *adminLogic) AdminAlbumCreate(ctx context.Context, req *vl_pb.ReqAdminAlbumCreate, reply *vl_pb.ResAdminAlbumCreate) error {
	id, err := v.CreateAlbum(ctx, &vl_pb.ReqCreateAlbum{
		Name:  req.Name,
		Intro: req.Intro,
		Cover: req.Cover,
		OpUid: req.OpUid,
	})
	if err != nil {
		return err
	}
	reply.Id = id
	return nil
}

func (v *adminLogic) AdminAlbumUpdate(ctx context.Context, req *vl_pb.ReqAdminAlbumUpdate, reply *vl_pb.ResAdminAlbumUpdate) error {
	return v.UpdateAlbum(ctx, &vl_pb.ReqUpdateAlbum{
		Id:    req.Id,
		Name:  req.Name,
		Intro: req.Intro,
		Cover: req.Cover,
		OpUid: req.OpUid,
	})
}

func (v *adminLogic) AdminAlbumDel(ctx context.Context, req *vl_pb.ReqAdminAlbumDel, reply *vl_pb.ResAdminAlbumDel) error {
	return v.DelAlbum(ctx, &vl_pb.ReqDelAlbum{
		Id:    req.Id,
		OpUid: req.OpUid,
	})
}

func (v *adminLogic) AdminAlbumDetail(ctx context.Context, req *vl_pb.ReqAdminAlbumDetail, reply *vl_pb.ResAdminAlbumDetail) error {
	albums, err := v.GetAlbumDetail(ctx, &vl_pb.ReqGetAlbumDetail{AlbumStr: []string{fmt.Sprintf("%d", req.Id)}})
	if err != nil {
		return err
	}
	if albums == nil || albums[req.Id] == nil {
		return nil
	}
	album := &vl_pb.AdminAlbum{
		Id:         albums[req.Id].Id,
		Name:       albums[req.Id].Name,
		Intro:      albums[req.Id].Intro,
		Cover:      albums[req.Id].Cover,
		AudioCount: int32(albums[req.Id].AudioCount),
		CreateTime: albums[req.Id].CreateTime,
		OpUid:      albums[req.Id].OpUid,
	}
	reply.Album = album
	return nil
}

func (v *adminLogic) AdminAlbumList(ctx context.Context, req *vl_pb.ReqAdminAlbumList, reply *vl_pb.ResAdminAlbumList) error {
	list, total, err := v.GetAlbumList(ctx, &vl_pb.ReqGetAlbumList{
		StartTime:     req.StartTime,
		EndTime:       req.EndTime,
		Name:          req.Name,
		CollectStatus: req.CollectStatus,
		Page:          req.Page,
		Limit:         req.Limit,
	})
	if err != nil {
		return err
	}
	reply.Total = total
	albums := make([]*vl_pb.AdminAlbum, 0)
	for _, l := range list {
		albums = append(albums, &vl_pb.AdminAlbum{
			Id:         l.Id,
			Name:       l.Name,
			Intro:      l.Intro,
			Cover:      l.Cover,
			OpUid:      l.OpUid,
			AudioCount: int32(l.AudioCount),
			CreateTime: l.CreateTime,
			HasSubject: l.HasSubject,
		})
	}
	reply.Albums = albums
	return nil
}

func (v *adminLogic) AdminAudioCollectList(ctx context.Context, req *vl_pb.ReqAdminAudioCollectList, reply *vl_pb.ResAdminAudioCollectList) error {
	q := service.VoiceLoverService.BuildAudioCollectSearchQuery(ctx, req)
	res, total, err := service.VoiceLoverService.SearchAudio(ctx, q)
	if err != nil {
		return err
	}
	reply.Audios = service.VoiceLoverService.BuildVoiceLoverAudioCollectPb(res)
	reply.Total = total
	return nil
}

func (v *adminLogic) AdminAudioCollect(ctx context.Context, req *vl_pb.ReqAdminAudioCollect, reply *vl_pb.ResAdminAudioCollect) error {
	return v.AudioCollect(ctx, &vl_pb.ReqAudioCollect{
		AudioId: req.AudioId,
		AlbumId: req.AlbumId,
		Type:    req.Type,
	})
}

func (v *adminLogic) AdminSubjectCreate(ctx context.Context, req *vl_pb.ReqAdminSubjectCreate, reply *vl_pb.ResAdminSubjectCreate) error {
	id, err := v.CreateSubject(ctx, &vl_pb.ReqCreateSubject{
		Name:  req.Name,
		OpUid: req.OpUid,
	})
	if err != nil {
		return err
	}
	reply.Id = id
	return nil
}

func (v *adminLogic) AdminSubjectUpdate(ctx context.Context, req *vl_pb.ReqAdminSubjectUpdate, reply *vl_pb.ResAdminSubjectUpdate) error {
	return v.UpdateSubject(ctx, &vl_pb.ReqUpdateSubject{
		Id:    req.Id,
		Name:  req.Name,
		OpUid: req.OpUid,
	})
}

func (v *adminLogic) AdminSubjectDel(ctx context.Context, req *vl_pb.ReqAdminSubjectDel, reply *vl_pb.ResAdminSubjectDel) error {
	return v.DelSubject(ctx, &vl_pb.ReqDelSubject{
		Id:    req.Id,
		OpUid: req.OpUid,
	})
}

func (v *adminLogic) AdminSubjectList(ctx context.Context, req *vl_pb.ReqAdminSubjectList, reply *vl_pb.ResAdminSubjectList) error {
	list, total, err := v.GetSubjectList(ctx, &vl_pb.ReqGetSubjectList{
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		Name:      req.Name,
		Page:      int32(req.Page),
		Limit:     int32(req.Limit),
	})
	if err != nil {
		return err
	}
	res := make([]*vl_pb.AdminSubjectData, 0)
	for _, subject := range list {
		res = append(res, &vl_pb.AdminSubjectData{
			Id:         subject.Id,
			Title:      subject.Name,
			AlbumTotal: uint32(subject.AlbumCount),
		})
	}
	reply.List = res
	reply.Total = total
	return nil
}

func (v *adminLogic) AdminAlbumCollect(ctx context.Context, req *vl_pb.ReqAdminAlbumCollect, reply *vl_pb.ResAdminAlbumCollect) error {
	return v.AlbumCollect(ctx, &vl_pb.ReqAlbumCollect{
		AlbumId:     req.AlbumId,
		SubjectId:   req.SubjectId,
		CollectType: req.Type,
	})
}

func (v *adminLogic) AdminAlbumCollectList(ctx context.Context, req *vl_pb.ReqAdminAlbumCollectList, reply *vl_pb.ResAdminAlbumCollectList) error {
	rep := &vl_pb.ResGetAlbumCollect{}
	err := v.GetAlbumCollect(ctx, &vl_pb.ReqGetAlbumCollect{
		AlbumStr:   req.AlbumStr,
		SubjectStr: req.SubjectStr,
		Page:       int32(req.Page),
		Limit:      int32(req.Limit),
	}, rep)
	if err != nil {
		return err
	}
	list := make([]*vl_pb.AdminAlbumSubject, 0)
	for _, a := range rep.AlbumCollects {
		list = append(list, &vl_pb.AdminAlbumSubject{
			Id:          a.Id,
			AlbumName:   a.AlbumName,
			SubjectName: a.SubjectName,
			AlbumId:     a.AlbumId,
			SubjectId:   a.SubjectId,
		})
	}
	reply.List = list
	reply.Total = rep.Total
	return nil
}

func (v *adminLogic) AdminSubjectDetail(ctx context.Context, req *vl_pb.ReqAdminSubjectDetail, reply *vl_pb.ResAdminSubjectDetail) error {
	subjects, err := v.GetSubjectDetail(ctx, &vl_pb.ReqGetSubjectDetail{Ids: []uint64{req.GetId()}})
	if err != nil {
		return err
	}
	if subjects == nil || subjects[req.Id] == nil {
		return nil
	}
	subject := subjects[req.Id]
	reply.Subject = &vl_pb.AdminSubjectData{
		Id:         subject.Id,
		Title:      subject.Name,
		AlbumTotal: uint32(subject.AlbumCount),
	}
	return nil
}

func (v *adminLogic) AdminAlbumChoice(ctx context.Context, req *vl_pb.ReqAdminAlbumChoice, reply *vl_pb.ResAdminAlbumChoice) error {
	return v.AlbumChoice(ctx, &vl_pb.ReqAlbumChoice{
		Id:   req.Id,
		Type: req.Choice,
	})
}

func (v *adminLogic) AdminAlbumChoiceList(ctx context.Context, req *vl_pb.ReqAdminAlbumChoiceList, reply *vl_pb.ResAdminAlbumChoiceList) error {
	choices, err := v.GetAlbumChoice(ctx, &vl_pb.ReqGetAlbumChoice{})
	if err != nil {
		return err
	}
	res := make([]*vl_pb.AdminAlbumData, 0)
	level := 1
	for _, r := range choices {
		res = append(res, &vl_pb.AdminAlbumData{
			Id:         r.Id,
			Name:       r.Name,
			CreateTime: r.CreateTime,
			Level:      int32(level),
		})
		level = level + 1
	}
	reply.Albums = res
	return nil
}

func (v *adminLogic) AdminBannerList(ctx context.Context, req *vl_pb.ReqAdminBannerList, reply *vl_pb.ResAdminBannerList) error {
	list, total, err := v.GetBannerList(ctx, &vl_pb.ReqGetBannerList{
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		Title:     req.Title,
		Status:    req.Status,
		Page:      req.Pate,
		Limit:     req.Limit,
	})
	if err != nil {
		return err
	}
	reply.Total = total
	for _, l := range list {
		reply.List = append(reply.List, &vl_pb.AdminBannerData{
			Id:         l.Id,
			Title:      l.Title,
			Cover:      l.Cover,
			Schema:     l.Schema,
			OpUid:      l.OpUid,
			Sort:       l.Sort,
			StartTime:  l.StartTime,
			EndTime:    l.EndTime,
			CreateTime: l.CreateTime,
		})
	}
	return nil
}

func (v *adminLogic) AdminBannerCreate(ctx context.Context, req *vl_pb.ReqAdminBannerCreate, reply *vl_pb.ResAdminBannerCreate) error {
	id, err := v.CreateBanner(ctx, &vl_pb.ReqCreateBanner{
		Title:     req.Title,
		Cover:     req.Cover,
		Schema:    req.Schema,
		OpUid:     req.OpUid,
		Sort:      req.Sort,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
	})
	if err != nil {
		return err
	}
	reply.Id = id
	return nil
}

func (v *adminLogic) AdminBannerUpdate(ctx context.Context, req *vl_pb.ReqAdminBannerUpdate, reply *vl_pb.ResAdminBannerUpdate) error {
	return v.UpdateBanner(ctx, &vl_pb.ReqUpdateBanner{
		Id:        req.Id,
		Title:     req.Title,
		Cover:     req.Cover,
		Schema:    req.Schema,
		OpUid:     req.OpUid,
		Sort:      req.Sort,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
	})
}

func (v *adminLogic) AdminBannerDetail(ctx context.Context, req *vl_pb.ReqAdminBannerDetail, reply *vl_pb.ResAdminBannerDetail) error {
	banner, err := v.GetBannerDetail(ctx, &vl_pb.ReqGetBannerDetail{
		Id: req.Id,
	})
	if err != nil {
		return err
	}
	if banner == nil {
		return nil
	}
	reply.Banner = &vl_pb.AdminBannerData{
		Id:         banner.Id,
		Title:      banner.Title,
		Cover:      banner.Cover,
		Schema:     banner.Schema,
		OpUid:      banner.OpUid,
		Sort:       banner.Sort,
		StartTime:  banner.StartTime,
		EndTime:    banner.EndTime,
		CreateTime: banner.CreateTime,
	}
	return nil
}
