package voice_lover

import (
	"context"
	"fmt"

	"github.com/gogf/gf/frame/g"
	"github.com/olaola-chat/rbp-functor/app/pb"
	"github.com/olaola-chat/rbp-proto/gen_pb/rpc/voice_lover"
	vl_rpc "github.com/olaola-chat/rbp-proto/rpcclient/voice_lover"
)

var ActivitySrv = &activitySrv{}

type activitySrv struct{}

func (a *activitySrv) GetInfo(ctx context.Context, id uint32) (*pb.VoiceLoverActivityMain, error) {
	activity, award, err := a.getInfo(ctx, id)
	if err != nil {
		return nil, err
	}
	return &pb.VoiceLoverActivityMain{Activity: activity, Award: award}, nil
}

func (a *activitySrv) getInfo(ctx context.Context, id uint32) (*pb.VoiceLoverActivity, *pb.VoiceLoverActivityAward, error) {
	// 获取活动详情
	activityResp, err := vl_rpc.VoiceLoverMain.GetActivity(ctx, &voice_lover.ReqGetActivity{Id: id})
	if err != nil || !activityResp.GetSuccess() {
		g.Log().Errorf("get activiy info err: %v, id: %d, msg: %s", err, id, activityResp.GetMsg())
		return nil, nil, fmt.Errorf("err: %v, msg: %s", err, activityResp.GetMsg())
	}
	activity := activityResp.GetActivity()
	g.Log().Infof("id: %d, activity: %+v", id, activity)

	// 获取排行奖励信息
	rankAwardResp, err := vl_rpc.VoiceLoverMain.GetRankAward(ctx, &voice_lover.ReqGetRankAward{Id: activity.GetRankAwardId()})
	if err != nil || !rankAwardResp.GetSuccess() {
		g.Log().Errorf("get activity rank award err: %v, rank_award_id: %d, msg: %s", err, activity.GetRankAwardId(), rankAwardResp.GetMsg())
		return nil, nil, fmt.Errorf("err: %v, msg: %s", err, rankAwardResp.GetMsg())
	}
	pkg := rankAwardResp.GetPackage()

	// 装扮列表
	var pretends []*pb.Pretend
	for _, v := range pkg.GetPretends() {
		pretends = append(pretends, &pb.Pretend{
			Id:   v.GetId(),
			Name: v.GetName(),
			Icon: v.GetIcon(),
		})
	}

	// 排名信息
	var ranks []*pb.RankInfo
	for _, v := range rankAwardResp.GetRanks() {
		ranks = append(ranks, &pb.RankInfo{
			Type:      v.GetType(),
			Start:     v.GetStart(),
			End:       v.GetEnd(),
			Days:      v.GetDays(),
			AwardName: pkg.GetName(),
		})
	}

	return &pb.VoiceLoverActivity{
		Id:          activity.GetId(),
		Title:       activity.GetTitle(),
		Intro:       activity.GetIntro(),
		Cover:       activity.GetCover(),
		StartTime:   activity.GetStartTime(),
		EndTime:     activity.GetEndTime(),
		RankAwardId: activity.GetRankAwardId(),
		RuleUrl:     activity.GetRuleUrl(),
	}, &pb.VoiceLoverActivityAward{Pretends: pretends, Ranks: ranks}, nil
}
