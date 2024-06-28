package client

import (
	"context"

	"github.com/olaola-chat/rbp-functor/app/pb"
)

// PayConsume 定一个payConsume单例
var Pretend *pretendConsume = &pretendConsume{
	&base{
		name: "rpc.pretend",
	},
}

type pretendConsume struct {
	*base
}

func (cli *pretendConsume) MGetPretends(ctx context.Context, ids []uint32) (resp *pb.MgetPretendsResp, err error) {
	resp = &pb.MgetPretendsResp{}
	req := &pb.MgetPretendsReq{
		PretendIds: ids,
	}
	return resp, cli.call(ctx, "MGetPretends", req, resp)
}
