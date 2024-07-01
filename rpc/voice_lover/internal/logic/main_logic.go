package logic

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	v8 "github.com/go-redis/redis/v8"
	"github.com/gogf/gf/database/gdb"
	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/util/gconv"
	"github.com/olaola-chat/rbp-functor/rpc/client"
	"github.com/olaola-chat/rbp-library/es"
	"github.com/olaola-chat/rbp-library/redis"
	"github.com/olaola-chat/rbp-proto/dao/functor"
	functor2 "github.com/olaola-chat/rbp-proto/gen_pb/db/functor"
	"github.com/olaola-chat/rbp-proto/gen_pb/db/xianshi"
	vl_pb "github.com/olaola-chat/rbp-proto/gen_pb/rpc/voice_lover"
	"github.com/olaola-chat/rbp-proto/rpcclient/user"

	userpb "github.com/olaola-chat/rbp-proto/gen_pb/rpc/user"

	voice_lover2 "github.com/olaola-chat/rbp-functor/app/model/voice_lover"
	"github.com/olaola-chat/rbp-functor/app/utils"
	"github.com/olaola-chat/rbp-functor/rpc/consts"
	"github.com/olaola-chat/rbp-functor/rpc/voice_lover/internal/dao"
)

type mainLogic struct {
	rds *v8.Client
}

var MainLogic = &mainLogic{
	rds: redis.RedisClient(consts.RedisDefault),
}

const (
	None = iota
	Dub
	Content
	Post
	Cover
)

func (m *mainLogic) Post(ctx context.Context, req *vl_pb.ReqPost, reply *vl_pb.ResBase) error {
	g.Log().Infof("VoiceLoverPost req = %v", req)
	now := uint64(time.Now().Unix())
	err := functor.VoiceLoverAudio.DB.Transaction(func(tx *gdb.TX) error {
		data := &functor2.EntityVoiceLoverAudio{
			Desc:       req.Desc,
			Resource:   req.Resource,
			Cover:      req.Cover,
			From:       uint32(req.Source),
			PubUid:     uint64(req.Uid),
			CreateTime: now,
			UpdateTime: now,
			Title:      req.Title,
			Labels:     req.Labels,
			Seconds:    req.Seconds,
		}
		if req.Uid == 200000169 {
			data.AuditStatus = dao.AuditPass
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

func buildAudioEsModel(data *functor2.EntityVoiceLoverAudio) *voice_lover2.VoiceLoverAudioEsModel {
	labelsSlice := make([]string, 0)
	for _, l := range strings.Split(data.Labels, ",") {
		if len(l) == 0 {
			continue
		}
		labelsSlice = append(labelsSlice, l)
	}
	esModel := &voice_lover2.VoiceLoverAudioEsModel{
		Id:          data.Id,
		PubUid:      uint32(data.PubUid),
		Title:       data.Title,
		Cover:       data.Cover,
		Desc:        data.Desc,
		CreateTime:  data.CreateTime,
		Labels:      labelsSlice,
		Source:      int32(data.From),
		AuditStatus: int32(data.AuditStatus),
		Albums:      []uint64{},
		HasAlbum:    0,
		OpUid:       data.OpUid,
		Resource:    data.Resource,
		Seconds:     data.Seconds,
	}
	return esModel
}

func (m *mainLogic) BuildRecAlbumsExtendInfo(ctx context.Context, infos []*vl_pb.AlbumData) {
	countMap := make(map[uint64]uint32)
	playCountsMap := make(map[uint64]uint64)
	wg := sync.WaitGroup{}
	mutex := &sync.Mutex{}
	for _, v := range infos {
		if _, ok := countMap[v.Id]; ok {
			continue
		}
		countMap[v.Id] = 0
		wg.Add(1)
		go func(albumId uint64) {
			defer wg.Done()
			total, _ := dao.VoiceLoverAudioAlbumDao.GetCountByAlbumId(ctx, albumId)
			countMap[albumId] = uint32(total)
			playCount := gconv.Uint64(m.rds.Get(ctx, consts.VoiceLoverAlbumPlayCount.Key(albumId)).Val())
			mutex.Lock()
			playCountsMap[albumId] = playCount
			mutex.Unlock()
		}(v.Id)
	}
	wg.Wait()
	for _, v := range infos {
		if count, ok := countMap[v.Id]; ok {
			v.AudioCount = count
		}
		if playCount, ok := playCountsMap[v.Id]; ok {
			v.PlayCount = playCount
		}
		if v.PlayCount < 10000 {
			v.PlayCountDesc = fmt.Sprintf("%d", v.PlayCount)
		} else {
			v.PlayCountDesc = fmt.Sprintf("%.1fw", float64(v.PlayCount)/10000.0)
		}
	}
}

func convertCoversToArray(coversStr string) []string {
	covers := make([]string, 0)
	for _, s := range strings.Split(coversStr, ",") {
		if len(s) == 0 {
			continue
		}
		covers = append(covers, s)
	}
	return covers
}

func convertLabelsToArray(labelsStr string) []string {
	labels := make([]string, 0)
	for _, s := range strings.Split(labelsStr, ",") {
		if len(s) == 0 {
			continue
		}
		labels = append(labels, s)
	}
	return labels
}

func (m *mainLogic) GetAlbumInfoById(ctx context.Context, req *vl_pb.ReqGetAlbumInfoById, reply *vl_pb.ResGetAlbumInfoById) error {
	albumInfo, err := dao.VoiceLoverAlbumDao.GetValidAlbumById(ctx, req.Id)
	if err != nil {
		return err
	}
	if albumInfo.GetId() == 0 {
		return gerror.New(fmt.Sprintf("album id=%d empty", req.Id))
	}
	reply.Album = &vl_pb.AlbumData{
		Id:         albumInfo.Id,
		Name:       albumInfo.Name,
		Intro:      albumInfo.Intro,
		Cover:      albumInfo.Cover,
		CreateTime: albumInfo.CreateTime,
	}
	m.BuildRecAlbumsExtendInfo(ctx, []*vl_pb.AlbumData{reply.Album})
	return nil
}

func (m *mainLogic) GetAlbumCommentCount(ctx context.Context, req *vl_pb.ReqGetAlbumCommentCount, reply *vl_pb.ResGetAlbumCommentCount) error {
	total, err := dao.VoiceLoverAlbumCommentDao.GetValidCommentCountByAlbumId(ctx, req.AlbumId)
	if err != nil {
		return err
	}
	reply.Total = uint32(total)
	return nil
}

func (m *mainLogic) GetRecAlbums(ctx context.Context, req *vl_pb.ReqGetRecAlbums, reply *vl_pb.ResGetRecAlbums) error {
	reply.Albums = make([]*vl_pb.AlbumData, 0)
	list, err := dao.VoiceLoverAlbumDao.GetValidAlbumListByChoice(ctx, dao.ChoiceRec, 0, 3)
	if err != nil {
		return err
	}
	for _, v := range list {
		reply.Albums = append(reply.Albums, &vl_pb.AlbumData{
			Id:         v.Id,
			Name:       v.Name,
			Intro:      v.Intro,
			Cover:      v.Cover,
			CreateTime: v.CreateTime,
		})
	}
	m.BuildRecAlbumsExtendInfo(ctx, reply.Albums)
	return nil
}

func (m *mainLogic) GetRecBanners(ctx context.Context, req *vl_pb.ReqGetRecBanners, reply *vl_pb.ResGetRecBanners) error {
	reply.Banners = make([]*vl_pb.BannerData, 0)
	list, err := dao.VoiceLoverBannerDao.GetValidListByLimit(ctx, 10)
	if err != nil {
		return err
	}
	for _, v := range list {
		reply.Banners = append(reply.Banners, &vl_pb.BannerData{
			Id:     v.Id,
			Cover:  v.Cover,
			Schema: v.Schema,
		})
	}
	return nil
}

func (m *mainLogic) GetRecCommonAlbums(ctx context.Context, req *vl_pb.ReqGetRecCommonAlbums, reply *vl_pb.ResGetRecAlbums) error {
	reply.Albums = make([]*vl_pb.AlbumData, 0)
	list, err := dao.VoiceLoverAlbumDao.GetValidAlbumListByChoice(ctx, dao.ChoiceDefault, 0, 3)
	if err != nil {
		return err
	}
	for _, v := range list {
		reply.Albums = append(reply.Albums, &vl_pb.AlbumData{
			Id:         v.Id,
			Name:       v.Name,
			Intro:      v.Intro,
			Cover:      v.Cover,
			CreateTime: v.CreateTime,
		})
	}
	m.BuildRecAlbumsExtendInfo(ctx, reply.Albums)
	return nil
}

func (m *mainLogic) GetAlbumsByPage(ctx context.Context, req *vl_pb.ReqGetAlbumsByPage, reply *vl_pb.ResGetAlbumsByPage) error {
	reply.Albums = make([]*vl_pb.AlbumData, 0)
	list, err := dao.VoiceLoverAlbumDao.GetValidAlbumListByChoice(ctx, req.Choice, int(req.Page), int(req.Limit)+1)
	if err != nil {
		return err
	}
	if len(list) > int(req.Limit) {
		list = list[:req.Limit]
		reply.HasMore = true
	}
	for _, v := range list {
		reply.Albums = append(reply.Albums, &vl_pb.AlbumData{
			Id:         v.Id,
			Name:       v.Name,
			Intro:      v.Intro,
			Cover:      v.Cover,
			CreateTime: v.CreateTime,
		})
	}
	m.BuildRecAlbumsExtendInfo(ctx, reply.Albums)
	return nil
}

func (m *mainLogic) GetSubjectAlbumsByPage(ctx context.Context, req *vl_pb.ReqGetSubjectAlbumsByPage, reply *vl_pb.ResGetAlbumsByPage) error {
	reply.Albums = make([]*vl_pb.AlbumData, 0)
	list, err := dao.VoiceLoverAlbumSubjectDao.GetListBySubjectId(ctx, req.SubjectId, int(req.Page), int(req.Limit)+1)
	if err != nil {
		return err
	}
	if len(list) > int(req.Limit) {
		list = list[:req.Limit]
		reply.HasMore = true
	}
	albumIds := make([]uint64, 0)

	for _, v := range list {
		albumIds = append(albumIds, v.AlbumId)
	}
	albumList, err := dao.VoiceLoverAlbumDao.GetValidAlbumListByIds(ctx, albumIds)
	for _, v := range albumList {
		reply.Albums = append(reply.Albums, &vl_pb.AlbumData{
			Id:         v.Id,
			Name:       v.Name,
			Intro:      v.Intro,
			Cover:      v.Cover,
			CreateTime: v.CreateTime,
		})
	}
	m.BuildRecAlbumsExtendInfo(ctx, reply.Albums)
	return nil
}

func (m *mainLogic) GetRecSubjects(ctx context.Context, req *vl_pb.ReqGetRecSubjects, reply *vl_pb.ResGetRecSubjects) error {
	reply.Subjects = make([]*vl_pb.SubjectData, 0)
	list, err := dao.VoiceLoverSubjectDao.GetValidSubjectList(ctx, 0, 10)
	if err != nil {
		return err
	}
	subjectIds := make([]uint64, 0)
	for _, v := range list {
		subjectIds = append(subjectIds, v.Id)
		reply.Subjects = append(reply.Subjects, &vl_pb.SubjectData{
			Id:         v.Id,
			Name:       v.Name,
			CreateTime: v.CreateTime,
			Albums:     make([]*vl_pb.AlbumData, 0),
		})
	}

	albumSubjectRelList, err := dao.VoiceLoverAlbumSubjectDao.GetListBySubjectIds(ctx, subjectIds)
	if err != nil {
		return err
	}
	subjectAlbumsMap := make(map[uint64][]uint64)
	albumIds := make([]uint64, 0)
	for _, v := range albumSubjectRelList {
		subjectAlbumsMap[v.SubjectId] = append(subjectAlbumsMap[v.SubjectId], v.AlbumId)
		albumIds = append(albumIds, v.AlbumId)
	}
	albumList, err := dao.VoiceLoverAlbumDao.GetValidAlbumListByIds(ctx, albumIds)
	if err != nil {
		return err
	}
	albums := make([]*vl_pb.AlbumData, 0)
	albumsMap := make(map[uint64]*vl_pb.AlbumData)
	for _, v := range albumList {
		albums = append(albums, &vl_pb.AlbumData{
			Id:         v.Id,
			Name:       v.Name,
			Intro:      v.Intro,
			Cover:      v.Cover,
			CreateTime: v.CreateTime,
		})
	}
	m.BuildRecAlbumsExtendInfo(ctx, albums)
	for _, v := range albums {
		albumsMap[v.Id] = v
	}

	for _, v := range reply.Subjects {
		if _, ok := subjectAlbumsMap[v.Id]; !ok {
			continue
		}
		for _, albumId := range subjectAlbumsMap[v.Id] {
			if _, ok := albumsMap[albumId]; !ok {
				continue
			}
			v.Albums = append(v.Albums, albumsMap[albumId])
		}
	}
	return nil
}

func (m *mainLogic) BatchGetAlbumAudioCount(ctx context.Context, req *vl_pb.ReqBatchGetAlbumAudioCount, reply *vl_pb.ResBatchGetAlbumAudioCount) error {
	reply.AlbumCounts = make(map[uint64]uint32)
	wg := sync.WaitGroup{}
	mutex := &sync.Mutex{}
	for _, v := range req.AlbumIds {
		if _, ok := reply.AlbumCounts[v]; ok {
			continue
		}
		reply.AlbumCounts[v] = 0
		wg.Add(1)
		go func(albumId uint64) {
			defer wg.Done()
			total, _ := dao.VoiceLoverAudioAlbumDao.GetCountByAlbumId(ctx, albumId)
			mutex.Lock()
			reply.AlbumCounts[albumId] = uint32(total)
			mutex.Unlock()
		}(v)
	}
	wg.Wait()
	return nil
}

func (m *mainLogic) IsUserCollectAlbum(ctx context.Context, req *vl_pb.ReqIsUserCollectAlbum, reply *vl_pb.ResIsUserCollectAlbum) error {
	reply.IsCollect = false
	// 如果UserCollectAlbumKey存在 0=未收藏 1=已收藏
	// 如果UserCollectAlbumKey存在 从mysql查一遍 写缓存
	key := consts.UserCollectAlbumKey.Key(req.Uid, req.AlbumId)
	if m.rds.Exists(ctx, key).Val() == 1 {
		if m.rds.Get(ctx, key).Val() == "1" {
			reply.IsCollect = true
		}
	} else {
		data, err := dao.VoiceLoverUserCollectDao.GetInfoByUidAndTypeAndId(ctx, req.Uid, req.AlbumId, dao.CollectTypeAlbum)
		if err != nil {
			return err
		}
		if data.GetId() > 0 {
			reply.IsCollect = true
		}
		defer func(isCollect bool) {
			value := 0
			if isCollect {
				value = 1
			}
			_ = m.rds.Set(ctx, key, value, consts.UserCollectAlbumKey.Ttl()).Err()
		}(reply.IsCollect)
	}
	return nil
}

func (m *mainLogic) IsUserCollectAlbums(ctx context.Context, req *vl_pb.ReqIsUserCollectAlbums, reply *vl_pb.ResIsUserCollectAlbums) error {
	reply.IsCollects = make([]bool, 0)
	isCollectMap := make(map[uint64]bool)
	wg := sync.WaitGroup{}
	mutex := &sync.Mutex{}
	for _, v := range req.AlbumIds {
		if _, ok := isCollectMap[v]; ok {
			continue
		}
		isCollectMap[v] = false
		wg.Add(1)
		go func(albumId uint64) {
			defer wg.Done()
			tmpReply := &vl_pb.ResIsUserCollectAlbum{}
			_ = m.IsUserCollectAlbum(ctx, &vl_pb.ReqIsUserCollectAlbum{Uid: req.Uid, AlbumId: albumId}, tmpReply)
			mutex.Lock()
			isCollectMap[albumId] = tmpReply.GetIsCollect()
			mutex.Unlock()
		}(v)
	}
	wg.Wait()
	for _, v := range req.AlbumIds {
		if _, ok := isCollectMap[v]; ok {
			reply.IsCollects = append(reply.IsCollects, isCollectMap[v])
		} else {
			reply.IsCollects = append(reply.IsCollects, false)
		}
	}
	return nil
}

func (m *mainLogic) Collect(ctx context.Context, req *vl_pb.ReqCollect, reply *vl_pb.ResCollect) error {
	var err error
	if req.Type == 0 {
		// 处理专辑
		key := consts.UserCollectAlbumKey.Key(req.Uid, req.Id)
		if req.From == 1 {
			// 收藏
			_, err = dao.VoiceLoverUserCollectDao.Add(ctx, req.Uid, req.Id, dao.CollectTypeAlbum)
			_ = m.rds.Set(ctx, key, 1, consts.UserCollectAlbumKey.Ttl())
		} else if req.From == 0 {
			// 取消收藏
			err = dao.VoiceLoverUserCollectDao.Delete(ctx, req.Uid, req.Id, dao.CollectTypeAlbum)
			_ = m.rds.Set(ctx, key, 0, consts.UserCollectAlbumKey.Ttl())
		} else {
			return gerror.New(fmt.Sprintf("param req.From=%d invalid", req.From))
		}

	} else if req.Type == 1 {
		// 处理音频
		key := consts.UserCollectAudioKey.Key(req.Uid, req.Id)
		if req.From == 1 {
			// 收藏
			_, err = dao.VoiceLoverUserCollectDao.Add(ctx, req.Uid, req.Id, dao.CollectTypeAudio)
			_ = m.rds.Set(ctx, key, 1, consts.UserCollectAudioKey.Ttl())
		} else if req.From == 0 {
			// 取消收藏
			err = dao.VoiceLoverUserCollectDao.Delete(ctx, req.Uid, req.Id, dao.CollectTypeAudio)
			_ = m.rds.Set(ctx, key, 0, consts.UserCollectAudioKey.Ttl())
		} else {
			return gerror.New(fmt.Sprintf("param req.From=%d invalid", req.From))
		}
	} else {
		return gerror.New(fmt.Sprintf("param req.Type=%d invalid", req.Type))
	}
	if err != nil {
		g.Log().Errorf("mainLogic Collect req=%+v||error=%v", req, err)
		return err
	}
	return nil
}

func (m *mainLogic) GetAlbumCollectList(ctx context.Context, req *vl_pb.ReqGetAlbumCollectList, reply *vl_pb.ResGetAlbumCollectList) error {
	reply.List = make([]*vl_pb.AlbumData, 0)
	list, err := dao.VoiceLoverUserCollectDao.GetListByUidAndType(ctx, req.Uid, dao.CollectTypeAlbum, int(req.Page), int(req.Limit)+1)
	if err != nil {
		return err
	}
	if len(list) > int(req.Limit) {
		list = list[:req.Limit]
		reply.HasMore = true
	}
	albumIds := make([]uint64, 0)
	for _, v := range list {
		albumIds = append(albumIds, v.CollectId)
	}
	albumList, err := dao.VoiceLoverAlbumDao.GetValidAlbumListByIds(ctx, albumIds)
	if err != nil {
		return err
	}
	for _, v := range albumList {
		reply.List = append(reply.List, &vl_pb.AlbumData{
			Id:         v.Id,
			Name:       v.Name,
			Intro:      v.Intro,
			Cover:      v.Cover,
			CreateTime: v.CreateTime,
		})
	}
	m.BuildRecAlbumsExtendInfo(ctx, reply.List)
	return nil
}

func (m *mainLogic) GetAudioCollectList(ctx context.Context, req *vl_pb.ReqGetAudioCollectList, reply *vl_pb.ResGetAudioCollectList) error {
	reply.List = make([]*vl_pb.AudioSimpleData, 0)
	list, err := dao.VoiceLoverUserCollectDao.GetListByUidAndType(ctx, req.Uid, dao.CollectTypeAudio, int(req.Page), int(req.Limit)+1)
	if err != nil {
		return err
	}
	if len(list) > int(req.Limit) {
		list = list[:req.Limit]
		reply.HasMore = true
	}
	audioIds := make([]uint64, 0)
	for _, v := range list {
		audioIds = append(audioIds, v.CollectId)
	}
	audioList, err := dao.VoiceLoverAudioDao.GetValidAudioListByIds(ctx, audioIds)
	if err != nil {
		return err
	}
	for _, v := range audioList {
		reply.List = append(reply.List, &vl_pb.AudioSimpleData{
			Id:       v.Id,
			Title:    v.Title,
			Resource: v.Resource,
			Covers:   convertCoversToArray(v.Cover),
			Uid:      uint32(v.PubUid),
		})
	}
	return nil
}

func (m *mainLogic) GetAudioListByAlbumId(ctx context.Context, req *vl_pb.ReqGetAudioListByAlbumId, reply *vl_pb.ResGetAudioListByAlbumId) error {
	reply.Audios = make([]*vl_pb.AudioSimpleData, 0)
	list, err := dao.VoiceLoverAudioAlbumDao.GetListByAlbumId(ctx, req.AlbumId)
	if err != nil {
		return err
	}
	audioIds := make([]uint64, 0)
	for _, v := range list {
		audioIds = append(audioIds, v.AudioId)
	}
	// 并发获取播放数量和详情
	wg := sync.WaitGroup{}
	audioDetailMap := make(map[uint64]*vl_pb.AudioSimpleData)
	audioPlayCountMap := make(map[uint64]uint64)
	wg.Add(2)
	go func() {
		defer wg.Done()
		audioList, tErr := dao.VoiceLoverAudioDao.GetAudioDetailsByAudioIds(ctx, audioIds)
		if tErr != nil {
			g.Log().Errorf("mainLogic GetAudioListByAlbumId GetAudioDetailsByAudioIds error=%v", tErr)
			return
		}
		for _, audioInfo := range audioList {
			audioDetailMap[audioInfo.Id] = &vl_pb.AudioSimpleData{
				Id:       audioInfo.Id,
				Title:    audioInfo.Title,
				Resource: audioInfo.Resource,
				Covers:   convertCoversToArray(audioInfo.Cover),
				Seconds:  audioInfo.Seconds,
				Uid:      uint32(audioInfo.PubUid),
			}
		}
	}()
	go func() {
		defer wg.Done()
		keys := make([]string, 0)
		for _, audioId := range audioIds {
			keys = append(keys, consts.VoiceLoverAudioPlayCount.Key(audioId))
		}
		vals, tErr := m.rds.MGet(ctx, keys...).Result()
		if tErr != nil {
			g.Log().Errorf("mainLogic GetAudioListByAlbumId MGet error=%v", tErr)
			return
		}
		if len(vals) != len(audioIds) {
			g.Log().Errorf("mainLogic GetAudioListByAlbumId MGet result error")
			return
		}
		for i, audioId := range audioIds {
			audioPlayCountMap[audioId] = gconv.Uint64(vals[i])
		}
	}()
	wg.Wait()

	// 组装返回数据
	for _, audioId := range audioIds {
		if _, ok := audioDetailMap[audioId]; !ok {
			continue
		}
		if playCount, ok := audioPlayCountMap[audioId]; ok {
			audioDetailMap[audioId].PlayCount = playCount
			if playCount < 10000 {
				audioDetailMap[audioId].PlayCountDesc = fmt.Sprintf("%d", playCount)
			} else {
				audioDetailMap[audioId].PlayCountDesc = fmt.Sprintf("%.1fw", float64(playCount)/10000.0)
			}
		}
		reply.Audios = append(reply.Audios, audioDetailMap[audioId])
	}
	return nil
}

func (m *mainLogic) SubmitAudioComment(ctx context.Context, req *vl_pb.ReqAudioSubmitComment, reply *vl_pb.ResCommonPost) error {
	data := g.Map{
		"audio_id":    req.AudioId,
		"content":     req.Content,
		"create_time": time.Now().Unix(),
		"update_time": time.Now().Unix(),
		"uid":         req.Uid,
		"type":        req.Type,
		"address":     req.Address,
	}
	success, err := dao.VoiceLoverAudioCommentDao.Insert(ctx, data)
	if err == nil && success {
		reply.Success = true
	}
	return nil
}

func (m *mainLogic) GetAudioCommentList(ctx context.Context, req *vl_pb.ReqGetAudioCommentList, reply *vl_pb.ResCommentList) error {
	commentList, err := dao.VoiceLoverAudioCommentDao.GetList(ctx, req.AudioId, req.Offset, req.Size)
	g.Log().Printf("GetAudioCommentList_list=>%v", commentList)
	if err != nil || len(commentList) == 0 {
		return errors.New("暂无数据")
	}
	reqUids := &userpb.ReqUserProfiles{
		Fields: []string{"uid", "icon", "name"},
	}
	for _, v := range commentList {
		reqUids.Uids = append(reqUids.Uids, uint32(v.Uid))
	}

	userList, err := user.UserProfile.Mget(ctx, reqUids)
	userMap := make(map[uint32]*xianshi.EntityXsUserProfile, 0)
	for _, v := range userList.Data {
		userMap[v.Uid] = v
	}

	for _, v := range commentList {
		tmp := &vl_pb.Comment{
			Id:         v.Id,
			Content:    v.Content,
			CreateTime: v.CreateTime,
			Address:    v.Address,
		}
		if profile, ok := userMap[uint32(v.Uid)]; ok {
			tmp.UserInfo = &vl_pb.CommentUser{
				Name:  profile.Name,
				Avtar: profile.Icon,
			}
		}
		reply.List = append(reply.List, tmp)
	}
	return nil
}

func (m *mainLogic) SubmitAlbumComment(ctx context.Context, req *vl_pb.ReqAlbumSubmitComment, reply *vl_pb.ResCommonPost) error {
	data := g.Map{
		"album_id":    req.AlbumId,
		"content":     req.Content,
		"create_time": time.Now().Unix(),
		"update_time": time.Now().Unix(),
		"uid":         req.Uid,
		"address":     req.Address,
	}
	_, err := dao.VoiceLoverAlbumCommentDao.Insert(ctx, data)
	if err == nil {
		reply.Success = true
	}
	return nil
}

func (m *mainLogic) GetAlbumCommentList(ctx context.Context, req *vl_pb.ReqGetAlbumCommentList, reply *vl_pb.ResCommentList) error {
	commentList, err := dao.VoiceLoverAlbumCommentDao.GetList(ctx, req.AlbumId, req.Offset, req.Size)
	if err != nil || len(commentList) == 0 {
		return errors.New("暂无数据")
	}
	reqUids := &userpb.ReqUserProfiles{
		Fields: []string{"uid", "icon", "name"},
	}
	for _, v := range commentList {
		reqUids.Uids = append(reqUids.Uids, uint32(v.Uid))
	}

	userList, err := user.UserProfile.Mget(ctx, reqUids)
	userMap := make(map[uint32]*xianshi.EntityXsUserProfile, 0)
	for _, v := range userList.Data {
		userMap[v.Uid] = v
	}

	for _, v := range commentList {
		tmp := &vl_pb.Comment{
			Id:         v.Id,
			Content:    v.Content,
			CreateTime: v.CreateTime,
			Address:    v.Address,
		}
		if profile, ok := userMap[uint32(v.Uid)]; ok {
			tmp.UserInfo = &vl_pb.CommentUser{
				Name:  profile.Name,
				Avtar: profile.Icon,
			}
		}
		reply.List = append(reply.List, tmp)
	}

	return nil
}

func (m *mainLogic) GetAudioInfoById(ctx context.Context, req *vl_pb.ReqGetAudioDetail, reply *vl_pb.ResGetAudioDetail) error {
	row, err := dao.VoiceLoverAudioDao.GetAudioDetailByAudioId(ctx, req.Id)
	if err != nil || row == nil {
		return errors.New("暂无该记录")
	}
	//音频基础信息
	reply.Audio = &vl_pb.AudioData{
		Id:           row.Id,
		Title:        row.Title,
		Desc:         row.Desc,
		Covers:       []string{row.Cover},
		Resource:     row.Resource,
		Labels:       convertLabelsToArray(row.Labels),
		Uid:          uint32(row.PubUid),
		CreateTime:   row.CreateTime,
		EditDubs:     make([]*vl_pb.AudioEditData, 0),
		EditContents: make([]*vl_pb.AudioEditData, 0),
		EditPosts:    make([]*vl_pb.AudioEditData, 0),
		EditCovers:   make([]*vl_pb.AudioEditData, 0),
		Seconds:      row.GetSeconds(),
	}

	// 获取音频播放数量
	val := m.rds.Get(ctx, consts.VoiceLoverAudioPlayCount.Key(req.Id)).Val()
	reply.Audio.PlayCount = gconv.Uint64(val)
	if reply.Audio.PlayCount < 10000 {
		reply.Audio.PlayCountDesc = fmt.Sprintf("%d", reply.Audio.PlayCount)
	} else {
		reply.Audio.PlayCountDesc = fmt.Sprintf("%.1fw", float64(reply.Audio.PlayCount)/10000.0)
	}

	//专辑基础信息
	albumIds, err := dao.VoiceLoverAudioAlbumDao.GetAlbumIdsByAudioId(ctx, req.Id)
	if err == nil && len(albumIds) > 0 {
		albumInfoMap, _ := dao.VoiceLoverAlbumDao.GetValidAlbumListByIds(ctx, albumIds)
		for _, info := range albumInfoMap {
			reply.Album = append(reply.Album, &vl_pb.AlbumData{
				Id:    info.Id,
				Name:  info.Name,
				Intro: info.Intro,
				Cover: info.Cover,
			})
		}
		m.BuildRecAlbumsExtendInfo(ctx, reply.Album)
	}

	// 参与人信息
	partnerList, _ := dao.VoiceLoverAudioPartnerDao.GetAudioPartnerByAudioId(ctx, req.Id)
	for _, e := range partnerList {
		if e.Type == Dub {
			reply.Audio.EditDubs = append(reply.Audio.EditDubs, &vl_pb.AudioEditData{
				Uid:  uint32(e.Uid),
				Type: e.Type,
			})
		}
		if e.Type == Content {
			reply.Audio.EditContents = append(reply.Audio.EditContents, &vl_pb.AudioEditData{
				Uid:  uint32(e.Uid),
				Type: e.Type,
			})
		}
		if e.Type == Post {
			reply.Audio.EditPosts = append(reply.Audio.EditPosts, &vl_pb.AudioEditData{
				Uid:  uint32(e.Uid),
				Type: e.Type,
			})
		}
		if e.Type == Cover {
			reply.Audio.EditCovers = append(reply.Audio.EditCovers, &vl_pb.AudioEditData{
				Uid:  uint32(e.Uid),
				Type: e.Type,
			})
		}
	}
	return nil
}

func (m *mainLogic) UpdateReportStatus(ctx context.Context, req *vl_pb.ReqUpdateStatus, reply *vl_pb.ResUpdateStatus) error {
	var r bool
	var err error
	if req.Type == 0 {
		r, err = dao.VoiceLoverAlbumCommentDao.UpdateStatus(ctx, req.Id, req.Status)
	} else if req.Type == 1 {
		r, err = dao.VoiceLoverAudioCommentDao.UpdateStatus(ctx, req.Id, req.Status)
	}
	if err == nil && r {
		reply.Success = true
	}
	return nil
}

func (m *mainLogic) PlayStatReport(ctx context.Context, req *vl_pb.ReqPlayStatReport, reply *vl_pb.ResPlayStatReport) error {
	_ = m.rds.Incr(ctx, consts.VoiceLoverAlbumPlayCount.Key(req.AlbumId))
	_ = m.rds.Incr(ctx, consts.VoiceLoverAudioPlayCount.Key(req.AudioId))
	return nil
}

func (m *mainLogic) IsUserCollectAudio(ctx context.Context, req *vl_pb.ReqCollect, reply *vl_pb.ResIsUserCollectAudio) error {
	reply.IsCollect = false
	// 如果UserCollectAlbumKey存在 0=未收藏 1=已收藏
	// 如果UserCollectAlbumKey存在 从mysql查一遍 写缓存
	key := consts.UserCollectAudioKey.Key(req.Uid, req.Id)
	if m.rds.Exists(ctx, key).Val() == 1 {
		if m.rds.Get(ctx, key).Val() == "1" {
			reply.IsCollect = true
		}
	} else {
		data, err := dao.VoiceLoverUserCollectDao.GetInfoByUidAndTypeAndId(ctx, req.Uid, req.Id, dao.CollectTypeAudio)
		if err != nil {
			return err
		}
		if data.GetId() > 0 {
			reply.IsCollect = true
		}
		defer func(isCollect bool) {
			value := 0
			if isCollect {
				value = 1
			}
			_ = m.rds.Set(ctx, key, value, consts.UserCollectAudioKey.Ttl()).Err()
		}(reply.IsCollect)
	}
	return nil
}

func (m *mainLogic) GetValidAudioUsers(ctx context.Context, req *vl_pb.ReqGetValidAudioUsers, reply *vl_pb.ResGetValidAudioUsers) error {
	reply.Uids = make([]uint32, 0)
	val := m.rds.SMembers(ctx, consts.VoiceLoverAudioPostUids.Key()).Val()
	if len(val) == 0 {
		// 从mysql查询数据
		list, err := dao.VoiceLoverAudioDao.GetValidUidsByUid(ctx, req.Uid)
		if err != nil {
			g.Log().Errorf("mainLogic GetValidAudioUsers GetValidUidsByUid error=%v", err)
			return err
		}
		tmpUids := make([]uint64, 0)
		for _, v := range list {
			tmpUids = append(tmpUids, v.PubUid)
		}
		uids := utils.DistinctUint64Slice(tmpUids)
		reply.Uids = gconv.Uint32s(uids)
		if len(uids) != 0 {
			_ = m.rds.SAdd(ctx, consts.VoiceLoverAudioPostUids.Key(), uids)
		}
	} else {
		reply.Uids = gconv.Uint32s(val)
	}
	return nil
}

// GetActivity 获取活动详情
func (m *mainLogic) GetActivity(ctx context.Context, req *vl_pb.ReqGetActivity, reply *vl_pb.RespGetActivity) error {
	data, err := dao.VoiceLoverActivityDao.GetOne(ctx, req.GetId())
	if err != nil {
		g.Log().Errorf("get activity err: %v, id: %d", err, req.GetId())
		reply.Msg = err.Error()
		return err
	}

	reply.Success = true
	reply.Activity = &vl_pb.ActivityInfo{
		Id:          data.GetId(),
		Title:       data.GetTitle(),
		Intro:       data.GetIntro(),
		Cover:       data.GetCover(),
		StartTime:   data.GetStartTime(),
		EndTime:     data.GetEndTime(),
		RankAwardId: data.GetRankAwardId(),
		RuleUrl:     data.GetRuleUrl(),
	}
	return nil
}

// GetRankAward 获取排行奖励详情
func (m *mainLogic) GetRankAward(ctx context.Context, req *vl_pb.ReqGetRankAward, reply *vl_pb.RespGetRankAward) error {
	rankAward, err := dao.VoiceLoverActivityRankAwardDao.GetOne(ctx, req.GetId())
	if err != nil {
		g.Log().Errorf("get rank award err: %v, id: %d", err, req.GetId())
		reply.Msg = err.Error()
		return err
	}

	// 获取奖励包信息
	pkgReply := &vl_pb.RespGetAwardPackage{}
	if err := m.GetAwardPackage(ctx, &vl_pb.ReqGetAwardPackage{Id: rankAward.GetPackageId()}, pkgReply); err != nil {
		g.Log().Errorf("get award package err: %v, id: %d", err, rankAward.GetPackageId())
		reply.Msg = err.Error()
		return err
	}

	reply.Success = true
	reply.Package = pkgReply.GetData()
	return nil
}

// GetAwardPackage 获取奖励包详情
func (m *mainLogic) GetAwardPackage(ctx context.Context, req *vl_pb.ReqGetAwardPackage, reply *vl_pb.RespGetAwardPackage) error {
	// 获取奖励包信息
	pkg, err := dao.VoiceLoverAwardPackageDao.GetOne(ctx, req.GetId())
	if err != nil {
		g.Log().Errorf("get award package err: %v, id: %d", err, req.GetId())
		reply.Msg = err.Error()
		return err
	}
	g.Log().Infof("tanlian get award package, id: %d, pkg: %+v", req.GetId(), pkg)

	// 批量获取装扮信息
	var pretendIds []uint32
	awards := make(map[string]string)
	if err := json.Unmarshal([]byte(pkg.GetAwards()), &awards); err != nil {
		g.Log().Errorf("unmashal package awards err: %v, awards: %s", err, pkg.GetAwards())
		reply.Msg = err.Error()
		return err
	}
	pretendIds = gconv.Uint32s(strings.Split(awards["pretend"], ","))
	pretendRsp, err := client.Pretend.MGetPretends(ctx, pretendIds)
	if err != nil {
		g.Log().Errorf("batch get pretend info err: %v, ids: %v", err, pretendIds)
		reply.Msg = err.Error()
		return err
	}

	var pretends []*vl_pb.PretendInfo
	for _, v := range pretendIds {
		item, ok := pretendRsp.GetData()[v]
		if !ok {
			continue
		}
		pretends = append(pretends, &vl_pb.PretendInfo{
			Id:   v,
			Name: item.GetName(),
			Icon: item.GetUrlSource(),
		})
	}

	reply.Success = true
	reply.Data = &vl_pb.AwardPackageInfo{
		Id:       pkg.GetId(),
		Name:     pkg.GetName(),
		Pretends: pretends,
	}
	return nil
}
