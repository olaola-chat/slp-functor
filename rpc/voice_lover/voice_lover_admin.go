package voice_lover

import (
	"context"
	"strings"

	vl_pb "github.com/olaola-chat/rbp-proto/gen_pb/rpc/voice_lover"

	"github.com/olaola-chat/rbp-functor/rpc/voice_lover/internal/dao"
)

func NewVoiceLoverAdmin() interface{} {
	return &VoiceLoverAdmin{}
}

type VoiceLoverAdmin struct {
}

func (v *VoiceLoverAdmin) GetAudioEdit(ctx context.Context, req *vl_pb.ReqGetAudioEdit, reply *vl_pb.ResGetAudioEdit) error {
	res, err := dao.VoiceLoverAudioPartnerDao.GetAudioPartnerByAudioId(ctx, req.GetId())
	if err != nil {
		return err
	}
	for _, r := range res {
		reply.Edits = append(reply.Edits, &vl_pb.AudioEditData{
			Uid:  r.Uid,
			Type: r.Type,
		})
	}
	return nil
}

func (v *VoiceLoverAdmin) GetAudioDetail(ctx context.Context, req *vl_pb.ReqGetAudioDetail, reply *vl_pb.ResGetAudioDetail) error {
	res, err := dao.VoiceLoverAudioDao.GetAudioDetailByAudioId(ctx, req.GetId())
	if err != nil {
		return err
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
	reply.Audio = &vl_pb.AudioData{
		Id:          res.Id,
		Uid:         res.PubUid,
		Resource:    res.Resource,
		Covers:      covers,
		Title:       res.Title,
		Desc:        res.Title,
		Labels:      labels,
		AuditStatus: int32(res.AuditStatus),
		CreateTime:  res.CreateTime,
		OpUid:       res.OpUid,
	}
	return nil
}
