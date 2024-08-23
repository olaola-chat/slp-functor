package logic

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/gogf/gf/util/gconv"

	"github.com/gogf/gf/database/gdb"
	"github.com/gogf/gf/frame/g"
	"github.com/olaola-chat/rbp-library/es"
	functor2 "github.com/olaola-chat/rbp-proto/dao/functor"
	"github.com/olaola-chat/rbp-proto/gen_pb/db/config"
	"github.com/olaola-chat/rbp-proto/gen_pb/db/functor"
	"github.com/olaola-chat/rbp-proto/gen_pb/rpc/voice_lover"

	"github.com/olaola-chat/rbp-functor/app/consts"
	"github.com/olaola-chat/rbp-functor/library/tool"
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
	g.Log().Infof("UpdateAudio recv req: %+v", req)
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

		/**if req.AuditStatus == dao.AuditPass {
			// 审核通过 发送mq消息
			if mqData, pErr := proto.Marshal(&pb.MQVoiceLoverAudioAuditPassData{AudioId: req.Id}); pErr == nil {
				_ = rocketmq.NewClient("default").Produce(consts.AudioPassTopic, mqData)
			}
		}**/
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
	list, total, err := dao.VoiceLoverAlbumDao.GetValidAlbumList(ctx, req.StartTime, req.EndTime, req.Name, req.CollectStatus, int(req.Page), int(req.Limit))
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
			HasSubject: int32(l.HasSubject),
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
	var (
		entity *functor.EntityVoiceLoverAudioAlbum
	)
	if req.Type == Collect {
		if audioAlbum != nil {
			return consts.ERROR_AUDIO_ALBUM_COLLECT
		}
		entity, err = dao.VoiceLoverAudioAlbumDao.Create(ctx, req.AudioId, req.AlbumId)
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
	if entity != nil {
		albumIds = append(albumIds, entity.AlbumId)
	}
	albumIds = tool.Slice.UniqueUint64Array(albumIds)
	data := g.Map{}
	data["has_album"] = 0
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
	album, err := dao.VoiceLoverAlbumDao.GetValidAlbumById(ctx, req.AlbumId)
	if err != nil {
		return err
	}
	if album == nil {
		return consts.ERROR_ALBUM_NOT_EXIST
	}
	subject, err := dao.VoiceLoverSubjectDao.GetValidSubjectById(ctx, req.SubjectId)
	if err != nil {
		return err
	}
	if subject == nil {
		return consts.ERROR_SUBJECT_NOT_EXIST
	}
	albumSubject, err := dao.VoiceLoverAlbumSubjectDao.GetAlbumSubjectByAIdAndSId(ctx, req.AlbumId, req.SubjectId)
	if err != nil {
		return err
	}
	count, err := dao.VoiceLoverAlbumSubjectDao.GetCountByAlbumId(ctx, req.GetAlbumId())
	if err != nil {
		return err
	}
	err = functor2.VoiceLoverAlbumSubject.DB.Transaction(func(tx *gdb.TX) error {
		if req.CollectType == Collect {
			if albumSubject != nil {
				return consts.ERROR_ALBUM_SUBJECT_COLLECT
			}
			err = dao.VoiceLoverAlbumSubjectDao.Create(tx, req.AlbumId, req.SubjectId)
			if err != nil {
				return err
			}
		}
		if req.CollectType == CollectRemove {
			if albumSubject == nil {
				return consts.ERROR_ALBUM_SUBJECT_COLLECT_REMOVE
			}
			err = dao.VoiceLoverAlbumSubjectDao.Delete(tx, req.AlbumId, req.SubjectId)
			if err != nil {
				return err
			}
		}
		if req.CollectType == Collect && album.HasSubject == 0 {
			err = dao.VoiceLoverAlbumDao.UpdateAlbumHasSubject(tx, req.AlbumId, int32(1))
		}
		if req.CollectType == CollectRemove && count == 1 {
			err = dao.VoiceLoverAlbumDao.UpdateAlbumHasSubject(tx, req.AlbumId, int32(0))
		}
		if err != nil {
			return err
		}
		return nil
	})
	return err
}

func (a *adminLogic) GetAlbumCollect(ctx context.Context, req *voice_lover.ReqGetAlbumCollect, reply *voice_lover.ResGetAlbumCollect) error {
	albumStr := req.AlbumStr
	subjectStr := req.SubjectStr
	var err error
	var albumId = dao.AlbumSubjectPlaceHolderID
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
	var subjectId = dao.AlbumSubjectPlaceHolderID
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

func (a *adminLogic) GetBannerList(ctx context.Context, req *voice_lover.ReqGetBannerList) ([]*functor.EntityVoiceLoverBanner, int32, error) {
	return dao.VoiceLoverBannerDao.GetBannerList(ctx, req.StartTime, req.EndTime, req.Title, req.Status, int(req.Page), int(req.Limit))
}

func (a *adminLogic) CreateBanner(ctx context.Context, req *voice_lover.ReqCreateBanner) (uint64, error) {
	return dao.VoiceLoverBannerDao.CreateBanner(ctx, req.StartTime, req.EndTime, req.Title, req.Cover, req.Schema, req.OpUid, req.Sort)
}

func (a *adminLogic) UpdateBanner(ctx context.Context, req *voice_lover.ReqUpdateBanner) error {
	return dao.VoiceLoverBannerDao.UpdateBanner(ctx, req.Id, req.StartTime, req.EndTime, req.Title, req.Cover, req.Schema, req.OpUid, req.Sort)
}

func (a *adminLogic) GetBannerDetail(ctx context.Context, req *voice_lover.ReqGetBannerDetail) (*functor.EntityVoiceLoverBanner, error) {
	return dao.VoiceLoverBannerDao.GetBannerById(ctx, req.Id)
}

// AddActivity 添加挑战/活动
func (a *adminLogic) AddActivity(ctx context.Context, req *voice_lover.ReqAdminAddActivity) (uint32, error) {
	if req.GetStartTime() > req.GetEndTime() {
		return 0, errors.New("开始时间不能晚于结束时间")
	}
	if req.GetStartTime() < time.Now().Unix() {
		return 0, errors.New("开始时间不可小于当前时间")
	}
	data := &config.EntityVoiceLoverActivity{
		Title:       req.GetTitle(),
		Intro:       req.GetIntro(),
		Cover:       req.GetCover(),
		StartTime:   uint32(req.GetStartTime()),
		EndTime:     uint32(req.GetEndTime()),
		RankAwardId: req.GetRankAwardId(),
		Id:          req.GetId(),
		RuleUrl:     req.GetJumpUrl(),
		CreateTime:  uint32(time.Now().Unix()),
		UpdateTime:  uint32(time.Now().Unix()),
	}
	id, err := dao.VoiceLoverActivityDao.Upsert(ctx, data)
	if err != nil {
		g.Log().Errorf("adminLogic AddActivity err: %v, req: %+v", err, req)
		return 0, err
	}
	return id, nil
}

// AddActivityAwardPackage 添加挑战奖励包配置
func (a *adminLogic) AddActivityAwardPackage(ctx context.Context, req *voice_lover.ReqAdminAddAwardPackage) (uint32, error) {
	if req.GetName() == "" {
		return 0, errors.New("名称不能为空")
	}
	if len(req.GetPretendIds()) == 0 {
		return 0, errors.New("奖励内容不能为空")
	}
	awardsMap := map[string]string{"pretend": strings.Join(gconv.Strings(req.PretendIds), ",")}
	awards, err := json.Marshal(awardsMap)
	if err != nil {
		g.Log().Errorf("adminLogic AddActivityAwardPackage marshal pretend id err: %v, req: %+v", err, req)
		return 0, err
	}

	data := &config.EntityVoiceLoverAwardPackage{
		Id:         req.GetId(),
		Name:       req.GetName(),
		Awards:     string(awards),
		CreateTime: uint32(time.Now().Unix()),
		UpdateTime: uint32(time.Now().Unix()),
	}
	id, err := dao.VoiceLoverAwardPackageDao.Upsert(ctx, data)
	if err != nil {
		g.Log().Errorf("adminLogic AddActivityAwardPackage err: %v, req: %+v", err, req)
		return 0, err
	}
	return id, nil
}

// AddActivityRankAward 添加挑战排行奖励配置
func (a *adminLogic) AddActivityRankAward(ctx context.Context, req *voice_lover.ReqAdminAddRankAward) (uint32, error) {
	if req.GetName() == "" {
		return 0, errors.New("名称不能为空")
	}
	if req.GetPackageId() == 0 {
		return 0, errors.New("奖励包不能为空")
	}
	if len(req.GetInfo()) == 0 {
		return 0, errors.New("名次奖励不能为空")
	}

	content, err := json.Marshal(req.GetInfo())
	if err != nil {
		g.Log().Errorf("adminLogic AddActivityRankAward marshal content err: %v, req: %+v", err, req)
		return 0, err
	}

	data := &config.EntityVoiceLoverActivityRankAward{
		Name:       req.GetName(),
		PackageId:  uint32(req.GetPackageId()),
		Content:    string(content),
		Id:         req.GetId(),
		CreateTime: uint32(time.Now().Unix()),
		UpdateTime: uint32(time.Now().Unix()),
	}
	id, err := dao.VoiceLoverActivityRankAwardDao.Upsert(ctx, data)
	if err != nil {
		g.Log().Errorf("adminLogic AddActivityRankAward err: %v, req: %+v", err, req)
		return 0, err
	}
	return id, nil
}

// AdminActivityList 获取挑战列表
func (a *adminLogic) AdminActivityList(ctx context.Context, req *voice_lover.ReqAdminActivityList) ([]*voice_lover.RespAdminActivityList_Item, int, error) {
	if req.Page == 0 {
		req.Page = 1
	}
	if req.Limit == 0 {
		req.Limit = 10
	}
	data, total, err := dao.VoiceLoverActivityDao.GetList(ctx, req.Id, req.Title, int(req.Page), int(req.Limit))
	if err != nil {
		g.Log().Errorf("adminLogic AdminActivityList err: %v, req: %+v", err, req)
		return nil, 0, err
	}

	// 批量获取排行奖励名称
	var rankAwardIds []uint32
	for _, v := range data {
		rankAwardIds = append(rankAwardIds, v.GetRankAwardId())
	}
	rankAwardMap, err := dao.VoiceLoverActivityRankAwardDao.BatchGet(ctx, rankAwardIds)
	if err != nil {
		g.Log().Errorf("adminLogic AdminActivityList err: %v, rankAwardIds: %v", err, rankAwardIds)
		return nil, 0, err
	}

	var items []*voice_lover.RespAdminActivityList_Item
	for _, v := range data {
		item := &voice_lover.RespAdminActivityList_Item{
			Id:            v.GetId(),
			Title:         v.GetTitle(),
			Intro:         v.GetIntro(),
			Cover:         v.GetCover(),
			StartTime:     int64(v.GetStartTime()),
			EndTime:       int64(v.GetEndTime()),
			RankAwardId:   v.GetRankAwardId(),
			RankAwardName: rankAwardMap[v.GetRankAwardId()].GetName(),
			JumpUrl:       v.GetRuleUrl(),
			CreateTime:    int64(v.GetCreateTime()),
			UpdateTime:    int64(v.GetUpdateTime()),
		}
		items = append(items, item)
	}
	return items, total, nil
}

// AdminAwardPackageList 获取奖励包列表
func (a *adminLogic) AdminAwardPackageList(ctx context.Context, req *voice_lover.ReqAdminAwardPackageList) ([]*voice_lover.RespAdminAwardPackageList_Item, int, error) {
	if req.Page == 0 {
		req.Page = 1
	}
	if req.Limit == 0 {
		req.Limit = 10
	}
	g.Log().Infof("AdminAwardPackageList recv req: %+v", req)
	data, total, err := dao.VoiceLoverAwardPackageDao.GetList(ctx, req.Id, req.Name, int(req.Page), int(req.Limit))
	if err != nil {
		g.Log().Errorf("adminLogic AdminAwardPackageList err: %v, req: %+v", err, req)
		return nil, 0, err
	}
	g.Log().Infof("AdminAwardPackageList data: %+v, total: %d", data, total)

	var items []*voice_lover.RespAdminAwardPackageList_Item
	for _, v := range data {
		if v.GetAwards() == "" {
			continue
		}
		awards := make(map[string]string)
		if err := json.Unmarshal([]byte(v.GetAwards()), &awards); err != nil {
			g.Log().Errorf("AdminAwardPackageList unmarshal awards err: %v, awards: %s", err, v.GetAwards())
			continue
		}
		pretends, ok := awards["pretend"]
		if !ok {
			g.Log().Errorf("AdminAwardPackageList invalid awards: %s", v.GetAwards())
			continue
		}
		item := &voice_lover.RespAdminAwardPackageList_Item{
			Id:         v.GetId(),
			Name:       v.GetName(),
			PretendIds: gconv.Uint32s(strings.Split(pretends, ",")),
			CreateTime: int64(v.GetCreateTime()),
			UpdateTime: int64(v.GetUpdateTime()),
		}
		items = append(items, item)
	}
	return items, total, nil
}

// AdminRankAwardList 获取排行奖励列表
func (a *adminLogic) AdminRankAwardList(ctx context.Context, req *voice_lover.ReqAdminRankAwardList) ([]*voice_lover.RespAdminRankAwardList_Item, int, error) {
	if req.Page == 0 {
		req.Page = 1
	}
	if req.Limit == 0 {
		req.Limit = 10
	}
	data, total, err := dao.VoiceLoverActivityRankAwardDao.GetList(ctx, req.Id, req.Name, int(req.Page), int(req.Limit))
	if err != nil {
		g.Log().Errorf("adminLogic AdminRankAwardList err: %v, req: %+v", err, req)
		return nil, 0, err
	}

	// 批量获取奖励包名称
	var pkgIds []uint32
	for _, v := range data {
		pkgIds = append(pkgIds, v.GetPackageId())
	}
	pkgMap, err := dao.VoiceLoverAwardPackageDao.BatchGet(ctx, pkgIds)
	if err != nil {
		g.Log().Errorf("adminLogic AdminRankAwardList err: %v, pkg ids: %v", err, pkgIds)
		return nil, 0, err
	}

	var items []*voice_lover.RespAdminRankAwardList_Item
	for _, v := range data {
		var ranks []*voice_lover.RankInfo
		if err := json.Unmarshal([]byte(v.GetContent()), &ranks); err != nil {
			g.Log().Errorf("adminLogic AdminRankAwardList unmarshal content err: %v, content: %s", err, v.GetContent())
			continue
		}
		item := &voice_lover.RespAdminRankAwardList_Item{
			Id:          v.GetId(),
			Name:        v.GetName(),
			PackageId:   v.GetPackageId(),
			PackageName: pkgMap[v.GetPackageId()].GetName(),
			Info:        ranks,
			CreateTime:  int64(v.GetCreateTime()),
			UpdateTime:  int64(v.GetUpdateTime()),
		}
		items = append(items, item)
	}
	return items, total, nil
}

// AdminActivityDelete 删除活动
func (a *adminLogic) AdminActivityDelete(ctx context.Context, id uint32) error {
	if err := dao.VoiceLoverActivityDao.Delete(ctx, id); err != nil {
		g.Log().Errorf("adminLogic AdminActivityDelete err: %v, id: %d", err, id)
		return err
	}
	return nil
}

// AdminAwardPackageDelete 删除奖励包
func (a *adminLogic) AdminAwardPackageDelete(ctx context.Context, id uint32) error {
	if rankAward, _ := dao.VoiceLoverActivityRankAwardDao.GetByPkgId(ctx, id); rankAward.GetId() > 0 {
		return errors.New("该奖励包已被应用，无法删除")
	}
	if err := dao.VoiceLoverAwardPackageDao.Delete(ctx, id); err != nil {
		g.Log().Errorf("adminLogic AdminAwardPackageDelete err: %v, id: %d", err, id)
		return err
	}
	return nil
}

// AdminRankAwardDelete 删除奖励排行
func (a *adminLogic) AdminRankAwardDelete(ctx context.Context, id uint32) error {
	if activity, _ := dao.VoiceLoverActivityDao.GetByRankAwardId(ctx, id); activity.GetId() > 0 {
		return errors.New("该奖励排行已被应用，无法删除")
	}
	if err := dao.VoiceLoverActivityRankAwardDao.Delete(ctx, id); err != nil {
		g.Log().Errorf("adminLogic AdminRankAwardDelete err: %v, id: %d", err, id)
		return err
	}
	return nil
}
