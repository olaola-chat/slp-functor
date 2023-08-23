package voice_lover

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/frame/g"
	"github.com/olaola-chat/rbp-proto/gen_pb/db/xianshi"

	"github.com/olaola-chat/rbp-proto/gen_pb/rpc/room"
	user_pb "github.com/olaola-chat/rbp-proto/gen_pb/rpc/user"
	vl_pb "github.com/olaola-chat/rbp-proto/gen_pb/rpc/voice_lover"
	rpcRoom "github.com/olaola-chat/rbp-proto/rpcclient/room"
	user_rpc "github.com/olaola-chat/rbp-proto/rpcclient/user"
	vl_rpc "github.com/olaola-chat/rbp-proto/rpcclient/voice_lover"

	"github.com/olaola-chat/rbp-library/redis"

	"github.com/olaola-chat/rbp-functor/app/pb"
	"github.com/olaola-chat/rbp-functor/app/query"
	//"github.com/olaola-chat/rbp-library/nsq"
)

var VoiceLoverService = &voiceLoverService{}

type voiceLoverService struct{}

func (serv *voiceLoverService) GetMainData(ctx context.Context, uid uint32) (*pb.RespVoiceLoverMain, error) {
	res := &pb.RespVoiceLoverMain{
		Success: true,
		Msg:     "",
		Data: &pb.VoiceLoverMain{
			RecAlbums:    make([]*pb.AlbumData, 0),
			RecBanners:   make([]*pb.BannerData, 0),
			RecUsers:     make([]*pb.UserData, 0),
			RecSubjects:  make([]*pb.SubjectData, 0),
			CommonAlbums: make([]*pb.AlbumData, 0),
		},
	}
	wg := sync.WaitGroup{}
	wg.Add(5)
	// 获取精选专辑推荐
	go func() {
		defer wg.Done()
		recAlbumList, err := vl_rpc.VoiceLoverMain.GetRecAlbums(ctx, &vl_pb.ReqGetRecAlbums{Uid: uid})
		if err != nil {
			g.Log().Errorf("voiceLoverService GetMainData GetRecAlbums error=%v", err)
			return
		}
		albumIds := make([]uint64, 0)
		for _, v := range recAlbumList.GetAlbums() {
			albumIds = append(albumIds, v.Id)
			res.Data.RecAlbums = append(res.Data.RecAlbums, &pb.AlbumData{
				Id:         v.Id,
				Title:      v.Name,
				Cover:      v.Cover,
				AudioTotal: v.AudioCount,
				PlayStats:  v.PlayCountDesc,
			})
		}
		// 判断用户是否已收藏专辑
		isCollectsRes, err := vl_rpc.VoiceLoverMain.IsUserCollectAlbums(ctx, &vl_pb.ReqIsUserCollectAlbums{Uid: uid, AlbumIds: albumIds})
		if err != nil {
			g.Log().Errorf("voiceLoverService GetMainData IsUserCollectAlbums error=%v", err)
			return
		}
		for i, v := range isCollectsRes.GetIsCollects() {
			res.Data.RecAlbums[i].IsCollect = v
		}
	}()
	// 获取banner推荐
	go func() {
		defer wg.Done()
	}()
	// 获取用户推荐
	go func() {
		defer wg.Done()
		recUids := []uint32{101000097}
		userInfosRes, err := user_rpc.UserProfile.Mget(ctx, &user_pb.ReqUserProfiles{Uids: recUids, Fields: []string{"name", "uid", "icon"}})
		if err != nil {
			g.Log().Errorf("voiceLoverService GetMainData Mget UserInfo error=%v", err)
			return
		}
		for _, v := range userInfosRes.GetData() {
			res.Data.RecUsers = append(res.Data.RecUsers, &pb.UserData{
				Uid:    v.Uid,
				Avatar: v.Icon,
				Name:   v.Name,
			})
		}
	}()
	// 获取话题推荐
	go func() {
		defer wg.Done()
		subjectList, err := vl_rpc.VoiceLoverMain.GetRecSubjects(ctx, &vl_pb.ReqGetRecSubjects{Uid: uid})
		if err != nil {
			g.Log().Errorf("voiceLoverService GetMainData GetRecSubjects error=%v", err)
			return
		}
		for _, v := range subjectList.GetSubjects() {
			subjectData := &pb.SubjectData{
				Id:     v.Id,
				Title:  v.Name,
				Albums: make([]*pb.AlbumData, 0),
			}
			for _, albumData := range v.Albums {
				subjectData.Albums = append(subjectData.Albums, &pb.AlbumData{
					Id:         albumData.Id,
					Title:      albumData.Name,
					Cover:      albumData.Cover,
					AudioTotal: albumData.AudioCount,
					PlayStats:  albumData.PlayCountDesc,
				})
			}
			res.Data.RecSubjects = append(res.Data.RecSubjects, subjectData)
		}
	}()
	// 获取普通专辑
	go func() {
		defer wg.Done()
		recAlbumList, err := vl_rpc.VoiceLoverMain.GetRecCommonAlbums(ctx, &vl_pb.ReqGetRecCommonAlbums{Uid: uid})
		if err != nil {
			g.Log().Errorf("voiceLoverService GetMainData GetRecCommonAlbums error=%v", err)
			return
		}
		for _, v := range recAlbumList.GetAlbums() {
			res.Data.CommonAlbums = append(res.Data.CommonAlbums, &pb.AlbumData{
				Id:         v.Id,
				Title:      v.Name,
				Cover:      v.Cover,
				AudioTotal: v.AudioCount,
				PlayStats:  v.PlayCountDesc,
			})
		}
	}()
	wg.Wait()
	return res, nil
}

func (serv *voiceLoverService) GetAlbumList(ctx context.Context, req *query.ReqAlbumList) (*pb.RespAlbumList, error) {
	res := &pb.RespAlbumList{
		Success: true,
		Msg:     "",
		Data: &pb.AlbumList{
			Albums:  make([]*pb.AlbumData, 0),
			HasMore: false,
		},
	}
	var albumsRes *vl_pb.ResGetAlbumsByPage
	var err error
	if req.Choice == 0 || req.Choice == 1 {
		// 查询默认或者精选专辑列表 直接查专辑表
		albumsRes, err = vl_rpc.VoiceLoverMain.GetAlbumsByPage(ctx, &vl_pb.ReqGetAlbumsByPage{
			Choice: req.Choice,
			Page:   req.Page,
			Limit:  req.Limit,
		})
		if err != nil {
			g.Log().Errorf("voiceLoverService GetAlbumList GetAlbumsByPage error=%v", err)
			return res, gerror.New("system error")
		}
	} else if req.Choice == 99 {
		// 查询专题下专辑列表
		albumsRes, err = vl_rpc.VoiceLoverMain.GetSubjectAlbumsByPage(ctx, &vl_pb.ReqGetSubjectAlbumsByPage{
			SubjectId: req.SubjectId,
			Page:      req.Page,
			Limit:     req.Limit,
		})
		if err != nil {
			g.Log().Errorf("voiceLoverService GetAlbumList GetSubjectAlbumsByPage error=%v", err)
			return res, gerror.New("system error")
		}
	} else {
		g.Log().Errorf("voiceLoverService GetAlbumList req.Choice=%d not supported", req.Choice)
		return res, gerror.New("param error")
	}
	res.Data.HasMore = albumsRes.GetHasMore()
	for _, v := range albumsRes.GetAlbums() {
		res.Data.Albums = append(res.Data.Albums, &pb.AlbumData{
			Id:         v.Id,
			Title:      v.Name,
			Cover:      v.Cover,
			AudioTotal: v.AudioCount,
			PlayStats:  v.PlayCountDesc,
		})
	}
	return res, nil
}

func (serv *voiceLoverService) GetAlbumDetail(ctx context.Context, uid uint32, albumId uint64) (*pb.RespAlbumDetail, error) {
	res := &pb.RespAlbumDetail{
		Success: true,
		Msg:     "",
		Data: &pb.AlbumDetail{
			Audios: make([]*pb.AudioData, 0),
		},
	}

	// 查询专辑主体信息
	albumInfoRes, err := vl_rpc.VoiceLoverMain.GetAlbumInfoById(ctx, &vl_pb.ReqGetAlbumInfoById{Id: albumId})
	if err != nil {
		g.Log().Errorf("voiceLoverService GetAlbumDetail GetAlbumInfoById error=%v", err)
		return res, gerror.New("system error")
	}
	if albumInfoRes.GetAlbum() == nil {
		g.Log().Errorf("voiceLoverService GetAlbumDetail GetAlbumInfoById empty||albumId=%d", albumId)
		return res, gerror.New("system error")
	}
	res.Data.Album = &pb.AlbumData{
		Id:         albumInfoRes.Album.Id,
		Title:      albumInfoRes.Album.Name,
		Cover:      albumInfoRes.Album.Cover,
		AudioTotal: albumInfoRes.Album.AudioCount,
	}

	// 专辑主体信息获取正常的话，并发获取其他数据
	wg := sync.WaitGroup{}
	wg.Add(3)
	// 用户是否已收藏
	go func() {
		defer wg.Done()
		isAlbumCollectRes, rErr := vl_rpc.VoiceLoverMain.IsUserCollectAlbum(ctx, &vl_pb.ReqIsUserCollectAlbum{
			AlbumId: albumId,
			Uid:     uid,
		})
		if rErr != nil {
			g.Log().Errorf("voiceLoverService GetAlbumDetail IsUserCollectAlbum error=%v", rErr)
			return
		}
		res.Data.IsCollected = isAlbumCollectRes.GetIsCollect()
	}()
	// 专辑评论数量
	go func() {
		defer wg.Done()
		albumCommentCountRes, rErr := vl_rpc.VoiceLoverMain.GetAlbumCommentCount(ctx, &vl_pb.ReqGetAlbumCommentCount{
			AlbumId: albumId,
		})
		if rErr != nil {
			g.Log().Errorf("voiceLoverService GetAlbumDetail GetAlbumCommentCount error=%v", rErr)
			return
		}
		commentCount := albumCommentCountRes.GetTotal()
		if commentCount < 10000 {
			res.Data.CommentCountDesc = fmt.Sprintf("%d", commentCount)
		} else {
			res.Data.CommentCountDesc = fmt.Sprintf("%.1fw", float64(commentCount)/10000.0)
		}
	}()
	// 获取音频列表
	go func() {
		defer wg.Done()
		audioListRes, rErr := vl_rpc.VoiceLoverMain.GetAudioListByAlbumId(ctx, &vl_pb.ReqGetAudioListByAlbumId{
			AlbumId: albumId,
		})
		if rErr != nil {
			g.Log().Errorf("voiceLoverService GetAlbumDetail GetAudioListByAlbumId error=%v", rErr)
			return
		}
		uids := make([]uint32, 0)
		for _, v := range audioListRes.GetAudios() {
			uids = append(uids, v.Uid)
			res.Data.Audios = append(res.Data.Audios, &pb.AudioData{
				Id:        v.Id,
				Title:     v.Title,
				Resource:  v.Resource,
				Covers:    v.Covers,
				Seconds:   v.Seconds,
				PlayStats: v.PlayCountDesc,
				UserInfo:  &pb.UserData{Uid: v.Uid},
			})
		}
		userInfosRes, _ := user_rpc.UserProfile.Mget(ctx, &user_pb.ReqUserProfiles{Uids: uids, Fields: []string{"name", "uid", "icon"}})
		userMap := make(map[uint32]*xianshi.EntityXsUserProfile, 0)
		for _, v := range userInfosRes.GetData() {
			userMap[v.Uid] = v
		}
		for _, v := range res.Data.Audios {
			if _, ok := userMap[v.UserInfo.Uid]; !ok {
				continue
			}
			v.UserInfo.Name = userMap[v.UserInfo.Uid].Name
			v.UserInfo.Avatar = userMap[v.UserInfo.Uid].Icon
		}
	}()
	wg.Wait()
	return res, nil
}

func (serv *voiceLoverService) GetAudioCommentList(ctx context.Context, audioId uint64, page, limit uint32) *pb.RespAudioComments {
	ret := &pb.RespAudioComments{
		Success: true,
		Msg:     "",
		Data: &pb.AudioComments{
			Comments: make([]*pb.CommentData, 0),
			HasMore:  false,
		},
	}
	if page <= 1 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}
	offset := (page - 1) * limit
	rows, err := vl_rpc.VoiceLoverMain.GetAudioCommentList(ctx, &vl_pb.ReqGetAudioCommentList{
		AudioId: audioId,
		Offset:  int32(offset),
		Size:    limit + 1,
	})
	if err != nil || len(rows.List) == 0 {
		ret.Msg = "暂无数据"
		return ret
	}

	ret.Success = true
	for k, v := range rows.List {
		if k >= int(limit) {
			ret.Data.HasMore = true
			break
		}
		tmp := &pb.CommentData{
			Id:      v.Id,
			Comment: v.Content,
			Address: v.Address,
			Date:    time.Unix(int64(v.CreateTime), 0).Local().Format("2006-01-02"),
		}

		if v.UserInfo != nil {
			tmp.UserInfo = &pb.UserData{
				Name:   v.UserInfo.Name,
				Avatar: v.UserInfo.Avtar,
			}
		}
		ret.Data.Comments = append(ret.Data.Comments, tmp)
	}

	return ret
}

func (serv *voiceLoverService) GetAlbumCommentList(ctx context.Context, albumId uint64, page, limit uint32) *pb.RespAlbumComments {
	ret := &pb.RespAlbumComments{
		Success: true,
		Msg:     "",
		Data: &pb.AlbumComments{
			Comments: make([]*pb.CommentData, 0),
			HasMore:  false,
		},
	}
	if page <= 1 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}
	offset := (page - 1) * limit
	rows, err := vl_rpc.VoiceLoverMain.GetAlbumCommentList(ctx, &vl_pb.ReqGetAlbumCommentList{
		AlbumId: albumId,
		Offset:  int32(offset),
		Size:    limit + 1,
	})
	if err != nil || len(rows.List) == 0 {
		ret.Msg = "暂无数据"
		return ret
	}
	ret.Success = true
	for k, v := range rows.List {
		if k >= int(limit) {
			ret.Data.HasMore = true
			break
		}

		tmp := &pb.CommentData{
			Id:      v.Id,
			Comment: v.Content,
			Address: v.Address,
			Date:    time.Unix(int64(v.CreateTime), 0).Local().Format("2006-01-02"),
		}
		if v.UserInfo != nil {
			tmp.UserInfo = &pb.UserData{
				Name:   v.UserInfo.Name,
				Avatar: v.UserInfo.Avtar,
			}
		}
		ret.Data.Comments = append(ret.Data.Comments, tmp)
	}

	return ret
}

func (serv *voiceLoverService) GetAudioDetail(ctx context.Context, uid uint32, audioId uint64) *pb.RespAudioDetail {
	res := &pb.RespAudioDetail{
		Success: true,
		Data:    &pb.AudioDetail{},
	}

	// 查询专辑主体信息
	detail, err := vl_rpc.VoiceLoverMain.GetAudioInfoById(ctx, &vl_pb.ReqGetAudioDetail{
		Id:  audioId,
		Uid: uid,
	})
	g.Log().Infof("GetAudioInfoById_r: %v", detail)
	if err != nil || detail == nil || detail.Audio == nil {
		res.Msg = "暂无数据"
		return res
	}
	var item []*pb.AlbumData
	for _, v := range detail.Album {
		item = append(item, &pb.AlbumData{
			Id:         v.Id,
			Cover:      v.Cover,
			Title:      v.Intro,
			PlayStats:  v.PlayCountDesc,
			AudioTotal: v.AudioCount,
		})
	}
	res.Data = &pb.AudioDetail{
		Audio: &pb.AudioData{
			Id:       detail.Audio.Id,
			Title:    detail.Audio.Title,
			Covers:   detail.Audio.Covers,
			Resource: detail.Audio.Resource,
			Seconds:  detail.Audio.Seconds,
		},
		Audios: item,
	}
	profile, err := user_rpc.UserProfile.Get(ctx, &user_pb.ReqUserProfile{
		Uid:    detail.Audio.Uid,
		Fields: []string{"name", "icon"},
	})
	if err == nil && profile != nil {
		res.Data.Audio.UserInfo = &pb.UserData{
			Uid:    detail.Audio.Uid,
			Name:   profile.Name,
			Avatar: profile.Icon,
		}
	}

	//是否关注了
	follow, err := user_rpc.UserProfile.CheckFollow(ctx, &user_pb.ReqCheckFollow{
		Uid:   uid,
		ToUid: detail.Audio.Uid,
	})
	if err == nil && follow != nil {
		res.Data.IsFollow = follow.IsFollow
	}

	//是否在房间
	roomInfo, err := rpcRoom.RoomInfo.InRoom(ctx, &room.ReqUid{Uid: detail.Audio.Uid})
	if err == nil && roomInfo.Rid > 0 {
		res.Data.RoomId = roomInfo.Rid
	}

	//是否收藏了
	collected, _ := vl_rpc.VoiceLoverMain.IsUserCollectAudio(ctx, &vl_pb.ReqCollect{
		Id:   res.Data.Audio.Id,
		Uid:  uid,
		Type: 1,
	})
	res.Data.IsCollected = collected.IsCollect

	return res
}

func (serv *voiceLoverService) Report(ctx context.Context, uniqueId uint32, desc string) *pb.CommonResp {
	//todo 发送nsq给后台
	//nsq.NewNsqClient().SendBytes()

	return &pb.CommonResp{
		Success: true,
	}
}

func (serv *voiceLoverService) SubmitAudioComment(ctx context.Context, req *vl_pb.ReqAudioSubmitComment) *pb.RespCommentAudio {
	ret := &pb.RespCommentAudio{}
	key := fmt.Sprintf("submit.audio.comment.%d", req.Uid)
	rds := redis.NewMutex("cache", key)
	success, err := rds.TryLockWithTtl(ctx, time.Second*3)
	if err != nil || !success {
		ret.Msg = "请勿频繁操作"
		return ret
	}
	_, err = vl_rpc.VoiceLoverMain.SubmitAudioComment(ctx, req)
	if err != nil {
		g.Log().Errorf("err: %v,%v", req, err)
		ret.Msg = "提交失败"
		return ret
	}
	ret.Success = true

	return ret
}

func (serv *voiceLoverService) SubmitAlbumComment(ctx context.Context, req *vl_pb.ReqAlbumSubmitComment) *pb.RespCommentAlbum {
	ret := &pb.RespCommentAlbum{}
	key := fmt.Sprintf("submit.album.comment.%d", req.Uid)
	rds := redis.NewMutex("cache", key)
	success, err := rds.TryLockWithTtl(ctx, time.Second*3)
	g.Log().Printf("debug: %v,%v,%v", success, err, req.Uid)
	if err != nil || !success {
		ret.Msg = "请勿频繁操作"
		return ret
	}
	_, err = vl_rpc.VoiceLoverMain.SubmitAlbumComment(ctx, req)
	if err != nil {
		g.Log().Errorf("err: %v,%v", req, err)
		ret.Msg = "提交失败"
		return ret
	}
	ret.Success = true

	return ret
}

func (serv *voiceLoverService) GetCollectAlbumList(ctx context.Context, uid uint32, req *query.ReqCollectAlbumList) (*pb.RespCollectAlbumList, error) {
	res := &pb.RespCollectAlbumList{
		Success: true,
		Msg:     "",
		Data: &pb.CollectAlbumList{
			List:    make([]*pb.AlbumData, 0),
			HasMore: false,
		},
	}
	albumListRes, err := vl_rpc.VoiceLoverMain.GetAlbumCollectList(ctx, &vl_pb.ReqGetAlbumCollectList{Uid: uid, Limit: req.Limit, Page: req.Page})
	if err != nil {
		g.Log().Errorf("voiceLoverService GetCollectAlbumList GetAlbumCollectList error=%v", err)
		return res, err
	}
	res.Data.HasMore = albumListRes.GetHasMore()
	for _, v := range albumListRes.GetList() {
		res.Data.List = append(res.Data.List, &pb.AlbumData{
			Id:         v.Id,
			Title:      v.Name,
			Cover:      v.Cover,
			AudioTotal: v.AudioCount,
			PlayStats:  v.PlayCountDesc,
		})
	}
	return res, nil
}

func (serv *voiceLoverService) GetCollectAudioList(ctx context.Context, uid uint32, req *query.ReqCollectAudioList) (*pb.RespCollectAudioList, error) {
	res := &pb.RespCollectAudioList{
		Success: true,
		Msg:     "",
		Data: &pb.CollectAudioList{
			List:    make([]*pb.AudioData, 0),
			HasMore: false,
		},
	}
	audioListRes, err := vl_rpc.VoiceLoverMain.GetAudioCollectList(ctx, &vl_pb.ReqGetAudioCollectList{Uid: uid, Limit: req.Limit, Page: req.Page})
	if err != nil {
		g.Log().Errorf("voiceLoverService GetCollectAudioList GetAudioCollectList error=%v", err)
		return res, err
	}
	res.Data.HasMore = audioListRes.GetHasMore()
	uids := make([]uint32, 0)
	for _, v := range audioListRes.GetList() {
		uids = append(uids, v.Uid)
		res.Data.List = append(res.Data.List, &pb.AudioData{
			Id:       v.Id,
			Title:    v.Title,
			Resource: v.Resource,
			Covers:   v.Covers,
			UserInfo: &pb.UserData{
				Uid: v.Uid,
			},
		})
	}
	userInfosRes, err := user_rpc.UserProfile.Mget(ctx, &user_pb.ReqUserProfiles{Uids: uids, Fields: []string{"name", "uid", "icon"}})
	if err != nil {
		g.Log().Errorf("voiceLoverService GetMainData Mget UserInfo error=%v", err)
		return res, nil
	}
	userMap := make(map[uint32]*xianshi.EntityXsUserProfile, 0)
	for _, v := range userInfosRes.GetData() {
		userMap[v.Uid] = v
	}
	for _, v := range res.Data.List {
		if _, ok := userMap[v.UserInfo.Uid]; !ok {
			continue
		}
		v.UserInfo.Name = userMap[v.UserInfo.Uid].Name
		v.UserInfo.Avatar = userMap[v.UserInfo.Uid].Icon
	}
	return res, nil
}
