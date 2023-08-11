package voice_lover

import (
	"context"
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
		Success:   true,
		Msg:       "",
		RecAlbums: make([]*pb.AlbumData, 0),
		HasMore:   false,
	}
	if req.Choice == 0 || req.Choice == 1 {
		// 查询默认或者精选专辑列表 直接查专辑表
	} else if req.Choice == 99 {
		// 查询专题下专辑列表
	} else {
		g.Log().Errorf("voiceLoverService GetAlbumList req.Choice=%d not supported", req.Choice)
		return res, gerror.New("param error")
	}
	return res, nil
}
