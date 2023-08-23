package voice_lover

import (
	"context"

	"github.com/gogf/gf/frame/g"
	vl_pb "github.com/olaola-chat/rbp-proto/gen_pb/rpc/voice_lover"

	"github.com/olaola-chat/rbp-functor/rpc/voice_lover/internal/logic"
)

func NewVoiceLoverAdmin() interface{} {
	return &VoiceLoverAdmin{}
}

type VoiceLoverAdmin struct {
}

func (v *VoiceLoverAdmin) GetAudioDetail(ctx context.Context, req *vl_pb.ReqGetAudioDetail, reply *vl_pb.ResGetAudioDetail) error {
	data, err := logic.AdminLogic.GetAudioDetail(ctx, req.GetId())
	if err != nil {
		return err
	}
	reply.Audio = data
	g.Log().Infof("GetAudioDetail data = %v", data)
	return nil
}

func (v *VoiceLoverAdmin) UpdateAudio(ctx context.Context, req *vl_pb.ReqUpdateAudio, reply *vl_pb.ResBase) error {
	return logic.AdminLogic.UpdateAudio(ctx, req)
}

func (v *VoiceLoverAdmin) AuditAudio(ctx context.Context, req *vl_pb.ReqAuditAudio, reply *vl_pb.ResAuditAudio) error {
	return logic.AdminLogic.AuditAudio(ctx, req)
}

func (v *VoiceLoverAdmin) CreateAlbum(ctx context.Context, req *vl_pb.ReqCreateAlbum, reply *vl_pb.ResCreateAlbum) error {
	id, err := logic.AdminLogic.CreateAlbum(ctx, req)
	if err != nil {
		return err
	}
	reply.Id = id
	return nil
}

func (v *VoiceLoverAdmin) DelAlbum(ctx context.Context, req *vl_pb.ReqDelAlbum, reply *vl_pb.ResDelAlbum) error {
	err := logic.AdminLogic.DelAlbum(ctx, req)
	if err != nil {
		return err
	}
	return nil
}

func (v *VoiceLoverAdmin) UpdateAlbum(ctx context.Context, req *vl_pb.ReqUpdateAlbum, reply *vl_pb.ResUpdateAlbum) error {
	err := logic.AdminLogic.UpdateAlbum(ctx, req)
	if err != nil {
		return err
	}
	return nil
}

func (v *VoiceLoverAdmin) GetAlbumDetail(ctx context.Context, req *vl_pb.ReqGetAlbumDetail, reply *vl_pb.ResGetAlbumDetail) error {
	albums, err := logic.AdminLogic.GetAlbumDetail(ctx, req)
	if err != nil {
		return err
	}
	reply.Albums = albums
	return nil
}

func (v *VoiceLoverAdmin) GetAlbumList(ctx context.Context, req *vl_pb.ReqGetAlbumList, reply *vl_pb.ResGetAlbumList) error {
	albums, total, err := logic.AdminLogic.GetAlbumList(ctx, req)
	if err != nil {
		return err
	}
	reply.Albums = albums
	reply.Total = total
	return nil
}

func (v *VoiceLoverAdmin) AudioCollect(ctx context.Context, req *vl_pb.ReqAudioCollect, reply *vl_pb.ResAudioCollect) error {
	return logic.AdminLogic.AudioCollect(ctx, req)
}

func (v *VoiceLoverAdmin) CreateSubject(ctx context.Context, req *vl_pb.ReqCreateSubject, reply *vl_pb.ResCreateSubject) error {
	id, err := logic.AdminLogic.CreateSubject(ctx, req)
	if err != nil {
		return err
	}
	reply.Id = id
	return nil
}

func (v *VoiceLoverAdmin) UpdateSubject(ctx context.Context, req *vl_pb.ReqUpdateSubject, reply *vl_pb.ResUpdateSubject) error {
	return logic.AdminLogic.UpdateSubject(ctx, req)
}

func (v *VoiceLoverAdmin) DelSubject(ctx context.Context, req *vl_pb.ReqDelSubject, reply *vl_pb.ResDelSubject) error {
	return logic.AdminLogic.DelSubject(ctx, req)
}

func (v *VoiceLoverAdmin) GetSubjectDetail(ctx context.Context, req *vl_pb.ReqGetSubjectDetail, reply *vl_pb.ResGetSubjectDetail) error {
	data, err := logic.AdminLogic.GetSubjectDetail(ctx, req)
	if err != nil {
		return err
	}
	reply.Subjects = data
	return nil
}

func (v *VoiceLoverAdmin) GetSubjectList(ctx context.Context, req *vl_pb.ReqGetSubjectList, reply *vl_pb.ResGetSubjectList) error {
	data, total, err := logic.AdminLogic.GetSubjectList(ctx, req)
	if err != nil {
		return err
	}
	reply.Subjects = data
	reply.Total = total
	return nil
}

func (v *VoiceLoverAdmin) AlbumCollect(ctx context.Context, req *vl_pb.ReqAlbumCollect, reply *vl_pb.ResAlbumCollect) error {
	err := logic.AdminLogic.AlbumCollect(ctx, req)
	if err != nil {
		return err
	}
	return nil
}

func (v *VoiceLoverAdmin) GetAlbumCollect(ctx context.Context, req *vl_pb.ReqGetAlbumCollect, reply *vl_pb.ResGetAlbumCollect) error {
	return logic.AdminLogic.GetAlbumCollect(ctx, req, reply)
}

func (v *VoiceLoverAdmin) AlbumChoice(ctx context.Context, req *vl_pb.ReqAlbumChoice, reply *vl_pb.ResAlbumChoice) error {
	return logic.AdminLogic.AlbumChoice(ctx, req)
}

func (v *VoiceLoverAdmin) GetAlbumChoice(ctx context.Context, req *vl_pb.ReqGetAlbumChoice, reply *vl_pb.ResGetAlbumChoice) error {
	list, err := logic.AdminLogic.GetAlbumChoice(ctx, req)
	if err != nil {
		return err
	}
	res := make([]*vl_pb.AlbumData, 0)
	for _, l := range list {
		res = append(res, &vl_pb.AlbumData{
			Id:         l.Id,
			Name:       l.Name,
			CreateTime: l.CreateTime,
		})
	}
	reply.Albums = res
	return nil
}

func (v *VoiceLoverAdmin) GetBannerList(ctx context.Context, req *vl_pb.ReqGetBannerList, reply *vl_pb.ResGetBannerList) error {
	list, count, err := logic.AdminLogic.GetBannerList(ctx, req)
	if err != nil {
		return err
	}
	res := make([]*vl_pb.BannerData, 0)
	for _, l := range list {
		res = append(res, &vl_pb.BannerData{
			Id:         l.Id,
			Title:      l.Title,
			Cover:      l.Cover,
			Schema:     l.Schema,
			OpUid:      l.OpUid,
			StartTime:  l.StartTime,
			EndTime:    l.EndTime,
			CreateTime: l.CreateTime,
			Sort:       l.Sort,
		})
	}
	reply.Banners = res
	reply.Total = count
	return nil
}

func (v *VoiceLoverAdmin) CreateBanner(ctx context.Context, req *vl_pb.ReqCreateBanner, reply *vl_pb.ResCreateBanner) error {
	id, err := logic.AdminLogic.CreateBanner(ctx, req)
	if err != nil {
		return err
	}
	reply.Id = id
	return nil
}

func (v *VoiceLoverAdmin) UpdateBanner(ctx context.Context, req *vl_pb.ReqUpdateBanner, reply *vl_pb.ResUpdateBanner) error {
	return logic.AdminLogic.UpdateBanner(ctx, req)
}

func (v *VoiceLoverAdmin) GetBannerDetail(ctx context.Context, req *vl_pb.ReqGetBannerDetail, reply *vl_pb.ResGetBannerDetail) error {
	b, err := logic.AdminLogic.GetBannerDetail(ctx, req)
	if err != nil {
		return err
	}
	if b == nil {
		return nil
	}
	reply.Banner = &vl_pb.BannerData{
		Id:         b.Id,
		Title:      b.Title,
		Cover:      b.Cover,
		Schema:     b.Schema,
		OpUid:      b.OpUid,
		StartTime:  b.StartTime,
		EndTime:    b.EndTime,
		Sort:       b.Sort,
		CreateTime: b.CreateTime,
	}
	return nil
}
