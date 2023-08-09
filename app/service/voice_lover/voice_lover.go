package voice_lover

import (
	"context"
	"sync"

	"github.com/gogf/gf/frame/g"
	vl_pb "github.com/olaola-chat/rbp-proto/gen_pb/rpc/voice_lover"
	vl_rpc "github.com/olaola-chat/rbp-proto/rpcclient/voice_lover"

	"github.com/olaola-chat/rbp-functor/app/pb"
)

var VoiceLoverService = &voiceLoverService{}

type voiceLoverService struct{}

func (serv *voiceLoverService) GetMainData(ctx context.Context, uid uint32) (*pb.RespVoiceLoverMain, error) {
	res := &pb.RespVoiceLoverMain{
		RecAlbums:    make([]*pb.AlbumRecData, 0),
		RecBanners:   make([]*pb.BannerRecData, 0),
		RecUsers:     make([]*pb.UserRecData, 0),
		RecSubjects:  make([]*pb.SubjectRecData, 0),
		CommonAlbums: make([]*pb.AlbumRecData, 0),
	}
	wg := sync.WaitGroup{}
	go func() {
		defer wg.Done()
		_, err := vl_rpc.VoiceLoverMain.GetRecAlbums(ctx, &vl_pb.ReqGetRecAlbums{Uid: 1})
		if err != nil {
			g.Log().Errorf("voiceLoverService GetMainData GetRecAlbums error=%v", err)
			return
		}
	}()
	return res, nil
}
