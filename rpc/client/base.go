package client

import (
	"context"
	"fmt"
	"time"

	"github.com/olaola-chat/slp-functor/app/model"
	"github.com/olaola-chat/slp-functor/library/rpc"

	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/frame/gins"
	rc "github.com/smallnest/rpcx/client"
	"github.com/smallnest/rpcx/share"
)

type base struct {
	name string
}

// RPCType 客户端发送RPC的方式
type RPCType int

const (
	// RPCTypeCall 单个请求
	RPCTypeCall RPCType = iota
	// RPCTypeFork 向包含此服务的所有服务器发送请求
	// 如果有任何服务器返回响应但没有错误，Fork则将为此XClient返回
	// 如果所有服务器都返回错误，则Fork返回这些错误中的一个错误
	RPCTypeFork
	// RPCTypeBroadcast 可以将一个请求发送到这个服务的所有节点
	// 如果所有的节点都正常返回，没有错误的话， Broadcast将返回其中的一个节点的返回结果
	// 如果有节点返回错误的话，Broadcast将返回这些错误信息中的一个
	RPCTypeBroadcast
)

func (serv *base) call(ctx context.Context, method string, req interface{}, reply interface{}, timeout ...time.Duration) error {
	return serv.exec(ctx, RPCTypeCall, method, req, reply, timeout...)
}

// 先注释掉，过ci检测，需要再打开

//func (serv *base) fork(ctx context.Context, method string, req interface{}, reply interface{}, timeout ...time.Duration) error {
//	return serv.exec(ctx, RPCTypeFork, method, req, reply, timeout...)
//}

//func (serv *base) broadcast(ctx context.Context, method string, req interface{}, reply interface{}, timeout ...time.Duration) error {
//	return serv.exec(ctx, RPCTypeBroadcast, method, req, reply, timeout...)
//}

func (serv *base) exec(ctx context.Context, rpcType RPCType, method string, req interface{}, reply interface{}, timeout ...time.Duration) error {
	t := time.Second
	if len(timeout) > 0 {
		t = timeout[0]
	}

	ctxTimeout, cancel := context.WithTimeout(ctx, t)
	c := model.GetContext(ctx)
	if c != nil && model.GetOpentracingEnabled(ctx) {
		ctxTimeout = context.WithValue(ctx, share.OpentracingSpanServerKey, c.Span)
	}

	var err error
	client := serv.getClient()
	switch rpcType {
	case RPCTypeCall:
		err = client.Call(ctxTimeout, method, req, reply)
	case RPCTypeFork:
		err = client.Fork(ctxTimeout, method, req, reply)
	case RPCTypeBroadcast:
		err = client.Broadcast(ctxTimeout, method, req, reply)
	}

	cancel()
	return err
}

func (serv *base) getClient() rc.XClient {
	instanceKey := fmt.Sprintf("self-rpc-service-client.%s", serv.name)
	// GetOrSetFuncLock里面使用g.Log()会产生死锁
	result := gins.GetOrSetFuncLock(instanceKey, func() interface{} {
		c, err := rpc.NewClientUsePb(serv.name)
		if err == nil {
			return c
		}
		return nil
	})
	if result == nil {
		panic(gerror.New("get rpc client error"))
	}
	if client, ok := result.(rc.XClient); ok {
		return client
	}
	panic(gerror.New("get rpc client error"))
}
