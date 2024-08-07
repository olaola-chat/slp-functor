package client

import (
	"context"

	"github.com/olaola-chat/rbp-functor/app/pb"
)

// IM 定一个imSender单例
var IM *imSender = &imSender{
	&base{
		name: "IM.Sender",
	},
}

type imSender struct {
	*base
}

// CheckDirty 敏感词检查
func (serv *imSender) CheckDirty(ctx context.Context, req *pb.ReqCheckDirty) (*pb.RepCheckDirty, error) {
	reply := &pb.RepCheckDirty{}
	err := serv.call(ctx, "CheckDirty", req, reply)
	return reply, err
}
