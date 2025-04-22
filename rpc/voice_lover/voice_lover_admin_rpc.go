package voice_lover

import (
	"context"

	vl_pb "github.com/olaola-chat/slp-proto/gen_pb/rpc/voice_lover"

	"github.com/olaola-chat/slp-functor/rpc/voice_lover/internal/logic"
)

func (v *VoiceLoverAdmin) AdminAudioList(ctx context.Context, req *vl_pb.ReqAdminAudioList, reply *vl_pb.ResAdminAudioList) error {
	return logic.AdminLogic.AdminAudioList(ctx, req, reply)
}

func (v *VoiceLoverAdmin) AdminAudioDetail(ctx context.Context, req *vl_pb.ReqAdminAudioDetail, reply *vl_pb.ResAdminAudioDetail) error {
	return logic.AdminLogic.AdminAudioDetail(ctx, req, reply)
}
func (v *VoiceLoverAdmin) AdminAudioUpdate(ctx context.Context, req *vl_pb.ReqAdminAudioUpdate, reply *vl_pb.ResAdminAudioUpdate) error {
	return logic.AdminLogic.AdminAudioUpdate(ctx, req, reply)
}

func (v *VoiceLoverAdmin) AdminAudioAudit(ctx context.Context, req *vl_pb.ReqAdminAudioAudit, reply *vl_pb.ResAdminAudioAudit) error {
	return logic.AdminLogic.AdminAudioAudit(ctx, req, reply)
}

func (v *VoiceLoverAdmin) AdminAudioAuditReason(ctx context.Context, req *vl_pb.ReqAdminAudioAuditReason, reply *vl_pb.ResAdminAudioAuditReason) error {
	return logic.AdminLogic.AdminAudioAuditReason(ctx, reply)
}

func (v *VoiceLoverAdmin) AdminAlbumCreate(ctx context.Context, req *vl_pb.ReqAdminAlbumCreate, reply *vl_pb.ResAdminAlbumCreate) error {
	return logic.AdminLogic.AdminAlbumCreate(ctx, req, reply)
}

func (v *VoiceLoverAdmin) AdminAlbumUpdate(ctx context.Context, req *vl_pb.ReqAdminAlbumUpdate, reply *vl_pb.ResAdminAlbumUpdate) error {
	return logic.AdminLogic.AdminAlbumUpdate(ctx, req, reply)
}

func (v *VoiceLoverAdmin) AdminAlbumDel(ctx context.Context, req *vl_pb.ReqAdminAlbumDel, reply *vl_pb.ResAdminAlbumDel) error {
	return logic.AdminLogic.AdminAlbumDel(ctx, req, reply)
}

func (v *VoiceLoverAdmin) AdminAlbumDetail(ctx context.Context, req *vl_pb.ReqAdminAlbumDetail, reply *vl_pb.ResAdminAlbumDetail) error {
	return logic.AdminLogic.AdminAlbumDetail(ctx, req, reply)
}

func (v *VoiceLoverAdmin) AdminAlbumList(ctx context.Context, req *vl_pb.ReqAdminAlbumList, reply *vl_pb.ResAdminAlbumList) error {
	return logic.AdminLogic.AdminAlbumList(ctx, req, reply)
}

func (v *VoiceLoverAdmin) AdminAudioCollectList(ctx context.Context, req *vl_pb.ReqAdminAudioCollectList, reply *vl_pb.ResAdminAudioCollectList) error {
	return logic.AdminLogic.AdminAudioCollectList(ctx, req, reply)
}

func (v *VoiceLoverAdmin) AdminAudioCollect(ctx context.Context, req *vl_pb.ReqAdminAudioCollect, reply *vl_pb.ResAdminAudioCollect) error {
	return logic.AdminLogic.AdminAudioCollect(ctx, req, reply)
}

func (v *VoiceLoverAdmin) AdminSubjectCreate(ctx context.Context, req *vl_pb.ReqAdminSubjectCreate, reply *vl_pb.ResAdminSubjectCreate) error {
	return logic.AdminLogic.AdminSubjectCreate(ctx, req, reply)
}

func (v *VoiceLoverAdmin) AdminSubjectUpdate(ctx context.Context, req *vl_pb.ReqAdminSubjectUpdate, reply *vl_pb.ResAdminSubjectUpdate) error {
	return logic.AdminLogic.AdminSubjectUpdate(ctx, req, reply)
}

func (v *VoiceLoverAdmin) AdminSubjectDel(ctx context.Context, req *vl_pb.ReqAdminSubjectDel, reply *vl_pb.ResAdminSubjectDel) error {
	return logic.AdminLogic.AdminSubjectDel(ctx, req, reply)
}

func (v *VoiceLoverAdmin) AdminSubjectList(ctx context.Context, req *vl_pb.ReqAdminSubjectList, reply *vl_pb.ResAdminSubjectList) error {
	return logic.AdminLogic.AdminSubjectList(ctx, req, reply)

}

func (v *VoiceLoverAdmin) AdminAlbumCollect(ctx context.Context, req *vl_pb.ReqAdminAlbumCollect, reply *vl_pb.ResAdminAlbumCollect) error {
	return logic.AdminLogic.AdminAlbumCollect(ctx, req, reply)

}

func (v *VoiceLoverAdmin) AdminAlbumCollectList(ctx context.Context, req *vl_pb.ReqAdminAlbumCollectList, reply *vl_pb.ResAdminAlbumCollectList) error {
	return logic.AdminLogic.AdminAlbumCollectList(ctx, req, reply)

}

func (v *VoiceLoverAdmin) AdminSubjectDetail(ctx context.Context, req *vl_pb.ReqAdminSubjectDetail, reply *vl_pb.ResAdminSubjectDetail) error {
	return logic.AdminLogic.AdminSubjectDetail(ctx, req, reply)

}

func (v *VoiceLoverAdmin) AdminAlbumChoice(ctx context.Context, req *vl_pb.ReqAdminAlbumChoice, reply *vl_pb.ResAdminAlbumChoice) error {
	return logic.AdminLogic.AdminAlbumChoice(ctx, req, reply)

}

func (v *VoiceLoverAdmin) AdminAlbumChoiceList(ctx context.Context, req *vl_pb.ReqAdminAlbumChoiceList, reply *vl_pb.ResAdminAlbumChoiceList) error {
	return logic.AdminLogic.AdminAlbumChoiceList(ctx, req, reply)
}

func (v *VoiceLoverAdmin) AdminBannerList(ctx context.Context, req *vl_pb.ReqAdminBannerList, reply *vl_pb.ResAdminBannerList) error {
	return logic.AdminLogic.AdminBannerList(ctx, req, reply)
}

func (v *VoiceLoverAdmin) AdminBannerCreate(ctx context.Context, req *vl_pb.ReqAdminBannerCreate, reply *vl_pb.ResAdminBannerCreate) error {
	return logic.AdminLogic.AdminBannerCreate(ctx, req, reply)
}

func (v *VoiceLoverAdmin) AdminBannerUpdate(ctx context.Context, req *vl_pb.ReqAdminBannerUpdate, reply *vl_pb.ResAdminBannerUpdate) error {
	return logic.AdminLogic.AdminBannerUpdate(ctx, req, reply)
}

func (v *VoiceLoverAdmin) AdminBannerDetail(ctx context.Context, req *vl_pb.ReqAdminBannerDetail, reply *vl_pb.ResAdminBannerDetail) error {
	return logic.AdminLogic.AdminBannerDetail(ctx, req, reply)
}

// AdminAddActivity 添加/更新活动
func (v *VoiceLoverAdmin) AdminAddActivity(ctx context.Context, req *vl_pb.ReqAdminAddActivity, reply *vl_pb.RespAdminAddActivity) error {
	id, err := logic.AdminLogic.AddActivity(ctx, req)
	if err != nil {
		reply.Msg = err.Error()
		return err
	}
	reply.Success = true
	reply.Id = id
	return nil
}

// AdminAddAwardPackage 添加/更新奖励包
func (v *VoiceLoverAdmin) AdminAddAwardPackage(ctx context.Context, req *vl_pb.ReqAdminAddAwardPackage, reply *vl_pb.RespAdminAddAwardPackage) error {
	id, err := logic.AdminLogic.AddActivityAwardPackage(ctx, req)
	if err != nil {
		reply.Msg = err.Error()
		return err
	}
	reply.Success = true
	reply.Id = id
	return nil
}

// AdminAddRankAward 添加/更新排行奖励
func (v *VoiceLoverAdmin) AdminAddRankAward(ctx context.Context, req *vl_pb.ReqAdminAddRankAward, reply *vl_pb.RespAdminAddRankAward) error {
	id, err := logic.AdminLogic.AddActivityRankAward(ctx, req)
	if err != nil {
		reply.Msg = err.Error()
		return err
	}
	reply.Success = true
	reply.Id = id
	return nil
}

// AdminActivityList 获取活动列表
func (v *VoiceLoverAdmin) AdminActivityList(ctx context.Context, req *vl_pb.ReqAdminActivityList, reply *vl_pb.RespAdminActivityList) error {
	items, total, err := logic.AdminLogic.AdminActivityList(ctx, req)
	if err != nil {
		reply.Msg = err.Error()
		return err
	}
	reply.Success = true
	reply.Data = items
	reply.Total = uint32(total)
	return nil
}

// AdminAwardPackageList 获取奖励包列表
func (v *VoiceLoverAdmin) AdminAwardPackageList(ctx context.Context, req *vl_pb.ReqAdminAwardPackageList, reply *vl_pb.RespAdminAwardPackageList) error {
	items, total, err := logic.AdminLogic.AdminAwardPackageList(ctx, req)
	if err != nil {
		reply.Msg = err.Error()
		return err
	}
	reply.Success = true
	reply.Data = items
	reply.Total = uint32(total)
	return nil
}

// AdminRankAwardList 获取排行奖励列表
func (v *VoiceLoverAdmin) AdminRankAwardList(ctx context.Context, req *vl_pb.ReqAdminRankAwardList, reply *vl_pb.RespAdminRankAwardList) error {
	items, total, err := logic.AdminLogic.AdminRankAwardList(ctx, req)
	if err != nil {
		reply.Msg = err.Error()
		return err
	}
	reply.Success = true
	reply.Data = items
	reply.Total = uint32(total)
	return nil
}

// AdminActivityDelete 删除活动
func (v *VoiceLoverAdmin) AdminActivityDelete(ctx context.Context, req *vl_pb.ReqAdminActivityDelete, reply *vl_pb.RespAdminActivityDelete) error {
	if err := logic.AdminLogic.AdminActivityDelete(ctx, req.GetId()); err != nil {
		reply.Msg = err.Error()
		return err
	}
	reply.Success = true
	return nil
}

// AdminAwardPackageDelete 删除奖励包
func (v *VoiceLoverAdmin) AdminAwardPackageDelete(ctx context.Context, req *vl_pb.ReqAdminAwardPackageDelete, reply *vl_pb.RespAdminAwardPackageDelete) error {
	if err := logic.AdminLogic.AdminAwardPackageDelete(ctx, req.GetId()); err != nil {
		reply.Msg = err.Error()
		return err
	}
	reply.Success = true
	return nil
}

// AdminRankAwardDelete 删除奖励排行
func (v *VoiceLoverAdmin) AdminRankAwardDelete(ctx context.Context, req *vl_pb.ReqAdminRankAwardDelete, reply *vl_pb.RespAdminRankAwardDelete) error {
	if err := logic.AdminLogic.AdminRankAwardDelete(ctx, req.GetId()); err != nil {
		reply.Msg = err.Error()
		return err
	}
	reply.Success = true
	return nil
}
