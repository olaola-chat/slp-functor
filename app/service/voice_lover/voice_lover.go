package voice_lover

import (
	"context"
	"errors"
	"sync"

	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/frame/g"

	vl_pb "github.com/olaola-chat/rbp-proto/gen_pb/rpc/voice_lover"
	vl_rpc "github.com/olaola-chat/rbp-proto/rpcclient/voice_lover"

	"github.com/olaola-chat/rbp-functor/app/pb"
	"github.com/olaola-chat/rbp-functor/app/query"
)

var VoiceLoverService = &voiceLoverService{}

type voiceLoverService struct{}

func (serv *voiceLoverService) GetMainData(ctx context.Context, uid uint32) (*pb.RespVoiceLoverMain, error) {
	res := &pb.RespVoiceLoverMain{
		Success:      true,
		Msg:          "",
		RecAlbums:    make([]*pb.AlbumData, 0),
		RecBanners:   make([]*pb.BannerData, 0),
		RecUsers:     make([]*pb.UserData, 0),
		RecSubjects:  make([]*pb.SubjectData, 0),
		CommonAlbums: make([]*pb.AlbumData, 0),
	}
	wg := sync.WaitGroup{}
	// 获取精选专辑推荐
	go func() {
		wg.Add(1)
		defer wg.Done()
		recAlbumList, err := vl_rpc.VoiceLoverMain.GetRecAlbums(ctx, &vl_pb.ReqGetRecAlbums{Uid: uid})
		if err != nil {
			g.Log().Errorf("voiceLoverService GetMainData GetRecAlbums error=%v", err)
			return
		}
		for _, v := range recAlbumList.GetAlbums() {
			res.RecAlbums = append(res.RecAlbums, &pb.AlbumData{
				Id:         v.Id,
				Title:      v.Name,
				Cover:      v.Cover,
				AudioTotal: v.AudioCount,
			})
		}
	}()
	// 获取话题推荐
	go func() {
		wg.Add(1)
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
				})
			}
		}
	}()
	wg.Wait()
	return res, nil
}

func (serv *voiceLoverService) GetAlbumList(ctx context.Context, req *query.ReqAlbumList) (*pb.RespAlbumList, error) {
	res := &pb.RespAlbumList{
		Success: true,
		Msg:     "",
		Albums:  make([]*pb.AlbumData, 0),
		HasMore: false,
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
	res.HasMore = albumsRes.GetHasMore()
	for _, v := range albumsRes.GetAlbums() {
		res.Albums = append(res.Albums, &pb.AlbumData{
			Id:         v.Id,
			Title:      v.Name,
			Cover:      v.Cover,
			AudioTotal: v.AudioCount,
		})
	}
	return res, nil
}

func (serv *voiceLoverService) GetAlbumDetail(ctx context.Context, uid uint32, albumId uint64) (*pb.RespAlbumDetail, error) {
	res := &pb.RespAlbumDetail{
		Success: true,
		Msg:     "",
		Audios:  make([]*pb.AudioData, 0),
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
	res.Album = &pb.AlbumData{
		Id:         albumInfoRes.Album.Id,
		Title:      albumInfoRes.Album.Name,
		Cover:      albumInfoRes.Album.Cover,
		AudioTotal: albumInfoRes.Album.AudioCount,
	}

	// 专辑主体信息获取正常的话，并发获取其他数据
	wg := sync.WaitGroup{}
	// 用户是否已收藏
	go func() {
		wg.Add(1)
		defer wg.Done()
		isAlbumCollectRes, rErr := vl_rpc.VoiceLoverMain.IsUserCollectAlbum(ctx, &vl_pb.ReqIsUserCollectAlbum{
			AlbumId: albumId,
			Uid:     uid,
		})
		if rErr != nil {
			g.Log().Errorf("voiceLoverService GetAlbumDetail IsUserCollectAlbum error=%v", rErr)
			return
		}
		res.IsCollected = isAlbumCollectRes.GetIsCollect()
	}()
	// 专辑评论数量
	go func() {
		wg.Add(1)
		defer wg.Done()
		albumCommentCountRes, rErr := vl_rpc.VoiceLoverMain.GetAlbumCommentCount(ctx, &vl_pb.ReqGetAlbumCommentCount{
			AlbumId: albumId,
		})
		if rErr != nil {
			g.Log().Errorf("voiceLoverService GetAlbumDetail GetAlbumCommentCount error=%v", rErr)
			return
		}
		res.CommentCount = albumCommentCountRes.GetTotal()
	}()
	// 获取音频列表
	go func() {
		wg.Add(1)
		defer wg.Done()
		audioListRes, rErr := vl_rpc.VoiceLoverMain.GetAudioListByAlbumId(ctx, &vl_pb.ReqGetAudioListByAlbumId{
			AlbumId: albumId,
		})
		if rErr != nil {
			g.Log().Errorf("voiceLoverService GetAlbumDetail GetAudioListByAlbumId error=%v", rErr)
			return
		}
		for _, v := range audioListRes.GetAudios() {
			res.Audios = append(res.Audios, &pb.AudioData{
				Id:        v.Id,
				Title:     v.Title,
				Resource:  v.Resource,
				Covers:    v.Covers,
				Seconds:   v.Seconds,
				PlayStats: "",
			})
		}
	}()
	wg.Wait()
	return res, nil
}

func (serv *voiceLoverService) GetAudioCommentList(ctx context.Context, audioId uint64) (*pb.RespAudioComments, error) {
	ret := &pb.RespAudioComments{}
	rows, err := vl_rpc.VoiceLoverMain.GetAudioCommentList(ctx, &vl_pb.ReqGetAudioCommentList{
		AudioId: audioId,
	})
	if err != nil || len(rows.List) == 0 {
		return nil, errors.New("暂无数据")
	}

	ret.Success = true
	for _, v := range rows.List {
		ret.Comments = append(ret.Comments, &pb.CommentData{
			Id: v.Id,
		})
	}

	return ret, nil
}

func (serv *voiceLoverService) GetAlbumCommentList(ctx context.Context, albumId uint64) (*pb.RespAlbumComments, error) {
	ret := &pb.RespAlbumComments{}
	rows, err := vl_rpc.VoiceLoverMain.GetAlbumCommentList(ctx, &vl_pb.ReqGetAlbumCommentList{
		AlbumId: albumId,
	})
	if err != nil || len(rows.List) == 0 {
		return nil, errors.New("暂无数据")
	}
	ret.Success = true
	for _, v := range rows.List {
		ret.Comments = append(ret.Comments, &pb.CommentData{
			Id: v.Id,
		})
	}

	return ret, nil
}
