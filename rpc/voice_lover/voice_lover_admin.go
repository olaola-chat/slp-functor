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
