package voice_lover

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/gogf/gf/util/gconv"

	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/frame/g"
	"github.com/olaola-chat/rbp-proto/gen_pb/db/xianshi"

	"github.com/olaola-chat/rbp-proto/gen_pb/rpc/room"
	user_pb "github.com/olaola-chat/rbp-proto/gen_pb/rpc/user"
	vl_pb "github.com/olaola-chat/rbp-proto/gen_pb/rpc/voice_lover"
	rpcRoom "github.com/olaola-chat/rbp-proto/rpcclient/room"
	user_rpc "github.com/olaola-chat/rbp-proto/rpcclient/user"
	vl_rpc "github.com/olaola-chat/rbp-proto/rpcclient/voice_lover"

	friend_pb "github.com/olaola-chat/rbp-proto/gen_pb/rpc/friends"
	friend_rpc "github.com/olaola-chat/rbp-proto/rpcclient/friends"

	redisV8 "github.com/go-redis/redis/v8"
	"github.com/olaola-chat/rbp-library/redis"

	"github.com/olaola-chat/rbp-functor/app/pb"
	"github.com/olaola-chat/rbp-functor/app/query"
	//"github.com/olaola-chat/rbp-library/nsq"
)

var VoiceLoverService = &voiceLoverService{}

type voiceLoverService struct{}

func (serv *voiceLoverService) GetMainData(ctx context.Context, uid, ver uint32) (*pb.RespVoiceLoverMain, error) {
	res := &pb.RespVoiceLoverMain{
		Success: true,
		Msg:     "",
		Data: &pb.VoiceLoverMain{
			RecAlbums:    make([]*pb.AlbumData, 0),
			RecBanners:   make([]*pb.BannerData, 0),
			RecUsers:     make([]*pb.UserData, 0),
			RecSubjects:  make([]*pb.SubjectData, 0),
			CommonAlbums: make([]*pb.AlbumData, 0),
			IsAnchor:     false,
		},
	}
	wg := sync.WaitGroup{}
	wg.Add(7)
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
		recBannersRes, err := vl_rpc.VoiceLoverMain.GetRecBanners(ctx, &vl_pb.ReqGetRecBanners{Uid: uid})
		if err != nil {
			g.Log().Errorf("voiceLoverService GetMainData GetRecBanners error=%v", err)
			return
		}
		for _, v := range recBannersRes.GetBanners() {
			// 低版本不显示声恋挑战活动入口
			if ver == 0 && strings.Contains(v.Schema, "page=sl_activity") {
				continue
			}
			res.Data.RecBanners = append(res.Data.RecBanners, &pb.BannerData{
				Id:          uint32(v.Id),
				ImgUrl:      v.Cover,
				RedirectUrl: v.Schema,
			})
		}
	}()
	// 获取用户推荐
	go func() {
		defer wg.Done()
		postUidsRes, _ := vl_rpc.VoiceLoverMain.GetValidAudioUsers(ctx, &vl_pb.ReqGetValidAudioUsers{Uid: uid})
		postUids := postUidsRes.GetUids()
		if len(postUids) == 0 {
			return
		}
		inRoomRes, _ := rpcRoom.RoomInfo.MgetInRoom(ctx, &room.ReqUids{Uids: postUids})
		inRoomMap := inRoomRes.GetData()
		inRoomUids := make([]uint32, 0)
		notInRoomUids := make([]uint32, 0)
		for _, v := range postUids {
			if rid, ok := inRoomMap[v]; ok {
				if rid > 0 {
					inRoomUids = append(inRoomUids, v)
				} else {
					notInRoomUids = append(notInRoomUids, v)
				}
			}
		}
		recUids := make([]uint32, 0)
		if len(inRoomUids) >= 5 {
			recUids = append(recUids, inRoomUids[:5]...)
		} else {
			recUids = append(recUids, inRoomUids...)
		}
		showNum := 5
		if len(recUids) < showNum {
			left := showNum - len(recUids)
			if len(notInRoomUids) >= left {
				recUids = append(recUids, notInRoomUids[:left]...)
			} else {
				recUids = append(recUids, notInRoomUids...)
			}
		}
		if len(recUids) == 0 {
			return
		}
		userInfosRes, err := user_rpc.UserProfile.Mget(ctx, &user_pb.ReqUserProfiles{Uids: recUids, Fields: []string{"name", "uid", "icon"}})
		if err != nil {
			g.Log().Errorf("voiceLoverService GetMainData Mget UserInfo error=%v", err)
			return
		}
		userInfosMap := make(map[uint32]*xianshi.EntityXsUserProfile)
		for _, v := range userInfosRes.GetData() {
			userInfosMap[v.Uid] = v
		}
		for _, v := range recUids {
			if _, ok := userInfosMap[v]; !ok {
				continue
			}
			rid := uint32(0)
			if value, ok := inRoomMap[v]; ok && value > 0 {
				rid = inRoomMap[v]
			}
			res.Data.RecUsers = append(res.Data.RecUsers, &pb.UserData{
				Uid:    userInfosMap[v].Uid,
				Avatar: userInfosMap[v].Icon,
				Name:   userInfosMap[v].Name,
				Rid:    rid,
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
	// 判断是否有有效工会
	go func() {
		defer wg.Done()
		isBrokerUserRes, _ := user_rpc.UserProfile.IsValidBrokerUser(ctx, &user_pb.ReqIsValidBrokerUser{Uid: uid})
		if isBrokerUserRes.GetResult() {
			res.Data.IsAnchor = true
		}
	}()
	// 获取全区动态数据
	go func() {
		defer wg.Done()
		// 获取排名最高的前10个声音作品
		rc := redis.RedisClient("user")
		rankKey := rc.Get(ctx, "rbp.voice.lover.audio.key").Val()
		if rankKey == "" {
			return
		}
		vals := rc.ZRevRangeByScore(ctx, rankKey, &redisV8.ZRangeBy{
			Min:   "0",
			Max:   "+inf",
			Count: 10,
		}).Val()
		if len(vals) == 0 {
			return
		}
		audioIds := gconv.Uint32s(vals)
		g.Log().Infof("tanlian get top rank audios: %v", audioIds)

		// 获取声音详情
		rsp, err := vl_rpc.VoiceLoverMain.BatchGetAudioInfo(ctx, &vl_pb.ReqBatchGetAudioInfo{AudioId: audioIds})
		if err != nil {
			g.Log().Errorf("batch get audio info err: %v, audio_ids: %v", err, audioIds)
			return
		}
		g.Log().Infof("tanlian rsp: %+v", rsp)
		for _, v := range rsp.GetItems() {
			audio := &pb.AudioData{
				Id:         uint64(v.GetId()),
				Title:      v.GetTitle(),
				Resource:   v.GetResource(),
				Covers:     []string{v.GetCover()},
				Seconds:    v.GetSeconds(),
				PlayStats:  formatPlayStats(v.GetPlayCnt()),
				UserInfo:   nil,
				Desc:       v.GetDesc(),
				CreateTime: uint64(v.GetCreateTime()),
				Partners:   nil,
			}
			res.Data.Audios = append(res.Data.Audios, audio)
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
		PlayStats:  albumInfoRes.Album.PlayCountDesc,
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

		// 批量判断用户是否收藏了音频
		var audioIds []uint64
		for _, v := range audioListRes.GetAudios() {
			audioIds = append(audioIds, v.GetId())
		}
		collectRsp, err := vl_rpc.VoiceLoverMain.BatchCheckUserCollect(ctx, &vl_pb.ReqBatchCheckUserCollect{Uid: uid, AudioId: gconv.Uint32s(audioIds)})
		if err != nil || !collectRsp.GetSuccess() {
			g.Log().Errorf("voiceLoverService BatchCheckUserCollect err: %v, uid: %d, audio_ids: %v", err, uid, audioIds)
		}

		// 批量获取音频的收藏数量
		numRsp, err := vl_rpc.VoiceLoverMain.BatchGetCollectNum(ctx, &vl_pb.ReqBatchGetCollectNum{CollectId: gconv.Uint32s(audioIds)})
		if err != nil || !collectRsp.GetSuccess() {
			g.Log().Errorf("voiceLoverService BatchGetCollectNum err: %v, uid: %d, audio_ids: %v", err, uid, audioIds)
		}

		uids := make([]uint32, 0)
		for _, v := range audioListRes.GetAudios() {
			uids = append(uids, v.Uid)
			res.Data.Audios = append(res.Data.Audios, &pb.AudioData{
				Id:         v.Id,
				Title:      v.Title,
				Resource:   v.Resource,
				Covers:     v.Covers,
				Seconds:    v.Seconds,
				PlayStats:  v.PlayCountDesc,
				UserInfo:   &pb.UserData{Uid: v.Uid},
				IsCollect:  collectRsp.GetCollectInfo()[uint32(v.Id)],
				CollectNum: numRsp.GetNums()[uint32(v.Id)],
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

	// 查询音频主体信息
	detail, err := vl_rpc.VoiceLoverMain.GetAudioInfoById(ctx, &vl_pb.ReqGetAudioDetail{
		Id:  audioId,
		Uid: uid,
	})
	if err != nil || detail == nil || detail.Audio == nil {
		res.Success = false
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
			Id:         detail.Audio.Id,
			Title:      detail.Audio.Title,
			Desc:       detail.Audio.Desc,
			Covers:     detail.Audio.Covers,
			Resource:   detail.Audio.Resource,
			Seconds:    detail.Audio.Seconds,
			CreateTime: detail.Audio.CreateTime,
			Labels:     detail.Audio.Labels,
			UserInfo:   &pb.UserData{Uid: detail.Audio.Uid},
			Partners:   make([]*pb.AudioPartner, 0),
			PlayStats:  detail.Audio.PlayCountDesc,
		},
		Albums: item,
	}
	uids := make([]uint32, 0)
	uids = append(uids, detail.Audio.Uid)
	// 处理参与人信息
	for _, v := range detail.Audio.EditDubs {
		uids = append(uids, v.Uid)
		res.Data.Audio.Partners = append(res.Data.Audio.Partners, &pb.AudioPartner{
			User: &pb.UserData{Uid: v.Uid},
			Tag:  "配音",
		})
	}
	for _, v := range detail.Audio.EditContents {
		uids = append(uids, v.Uid)
		res.Data.Audio.Partners = append(res.Data.Audio.Partners, &pb.AudioPartner{
			User: &pb.UserData{Uid: v.Uid},
			Tag:  "文案",
		})
	}
	for _, v := range detail.Audio.EditCovers {
		uids = append(uids, v.Uid)
		res.Data.Audio.Partners = append(res.Data.Audio.Partners, &pb.AudioPartner{
			User: &pb.UserData{Uid: v.Uid},
			Tag:  "封面设计",
		})
	}
	for _, v := range detail.Audio.EditPosts {
		uids = append(uids, v.Uid)
		res.Data.Audio.Partners = append(res.Data.Audio.Partners, &pb.AudioPartner{
			User: &pb.UserData{Uid: v.Uid},
			Tag:  "后期",
		})
	}
	userInfosRes, err := user_rpc.UserProfile.Mget(ctx, &user_pb.ReqUserProfiles{
		Uids:   uids,
		Fields: []string{"name", "icon", "uid"},
	})
	if err != nil {
		g.Log().Errorf("voiceLoverService GetAudioDetail Mget userinfo error=%v", err)
	} else {
		userInfosMap := make(map[uint32]*xianshi.EntityXsUserProfile)
		for _, v := range userInfosRes.GetData() {
			userInfosMap[v.Uid] = v
		}
		if userInfo, ok := userInfosMap[res.Data.Audio.UserInfo.Uid]; ok {
			res.Data.Audio.UserInfo.Name = userInfo.Name
			res.Data.Audio.UserInfo.Avatar = userInfo.Icon
		}
		for _, v := range res.Data.Audio.Partners {
			if userInfo, ok := userInfosMap[v.User.Uid]; ok {
				v.User.Name = userInfo.Name
				v.User.Avatar = userInfo.Icon
			}
		}
	}
	// 粉丝数量
	friendRes, _ := friend_rpc.Friend.Count(ctx, &friend_pb.ReqFriendCount{Uid: res.Data.Audio.UserInfo.Uid, Follow: true})
	res.Data.Audio.UserInfo.FansNum = friendRes.GetFollow()

	//是否关注了
	//followRes, _ := friend_rpc.Friend.IsFollow(ctx, &friend_pb.ReqIsFollow{Uid: res.Data.Audio.UserInfo.Uid, Uids: []uint32{uid}})
	followRes, err := friend_rpc.Friend.IsFollow(ctx, &friend_pb.ReqIsFollow{Uid: 101000023, Uids: []uint32{101000097}})
	if err != nil {
		g.Log().Errorf("voiceLoverService GetAudioDetail IsFollow error=%v", err)
	} else {
		g.Log().Debugf("testlw||uid=%d||toUid=%d||res=%+v", uid, res.Data.Audio.UserInfo.Uid, followRes.GetData())
		if len(followRes.GetData()) == 1 && followRes.Data[0].GetFollow() == 1 {
			res.Data.IsFollow = true
		}
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

	// 收藏数
	if collectNumRsp, _ := vl_rpc.VoiceLoverMain.BatchGetCollectNum(ctx, &vl_pb.ReqBatchGetCollectNum{CollectId: []uint32{uint32(audioId)}}); collectNumRsp.GetSuccess() {
		res.Data.Audio.CollectNum = collectNumRsp.GetNums()[uint32(audioId)]
	}

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

func (serv *voiceLoverService) ShareAlbumInfo(ctx context.Context, uid uint32, req *query.ReqShareAlbum) (*pb.RespShareInfo, error) {
	res := &pb.RespShareInfo{
		Success: true,
		Msg:     "",
		Data: &pb.ShareData{
			ShareTitle: "少年骨骼清奇，快跟我一起闯荡彩虹星球直播江湖！",
			ShareDesc:  "",
			ShareUrl:   "https://www.caihongmeng.com",
			ShareIcon:  "",
		},
	}
	userInfo, err := user_rpc.UserProfile.Get(ctx, &user_pb.ReqUserProfile{Uid: uid, Fields: []string{"name", "uid", "icon"}})
	if err != nil {
		return res, err
	}
	res.Data.ShareIcon = userInfo.GetIcon()
	return res, nil
}

func (serv *voiceLoverService) ShareAudioInfo(ctx context.Context, uid uint32, req *query.ReqShareAudio) (*pb.RespShareInfo, error) {
	res := &pb.RespShareInfo{
		Success: true,
		Msg:     "",
		Data: &pb.ShareData{
			ShareTitle: "少年骨骼清奇，快跟我一起闯荡彩虹星球直播江湖！",
			ShareDesc:  "",
			ShareUrl:   "https://www.caihongmeng.com",
			ShareIcon:  "",
		},
	}
	userInfo, err := user_rpc.UserProfile.Get(ctx, &user_pb.ReqUserProfile{Uid: uid, Fields: []string{"name", "uid", "icon"}})
	if err != nil {
		return res, err
	}
	res.Data.ShareIcon = userInfo.GetIcon()
	return res, nil
}

func formatPlayStats(playCnt uint32) string {
	if playCnt < 10000 {
		return fmt.Sprintf("%d", playCnt)
	}
	return fmt.Sprintf("%.1fw", float64(playCnt)/10000.0)
}
