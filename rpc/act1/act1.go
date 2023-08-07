package act1

import (
	"context"

	"github.com/olaola-chat/rbp-functor/rpc/act1/internal/logic"
)

func NewAct1() interface{} {
	return &Act1{}
}

type Act1 struct {
}

type xxpb struct {
}

func (a *Act1) Method1(ctx context.Context, req *xxpb, reply *xxpb) error {

	l := &logic.Act1Logic{}
	l.Method1()

	return nil
}
