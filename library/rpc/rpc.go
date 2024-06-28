package rpc

import (
	"fmt"
	"time"

	"github.com/olaola-chat/rbp-functor/library/rpc/plugins"

	"github.com/gogf/gf/frame/g"
	"github.com/rcrowley/go-metrics"
	"github.com/rpcxio/libkv/store"
	"github.com/smallnest/rpcx/client"
	"github.com/smallnest/rpcx/server"
	"github.com/smallnest/rpcx/serverplugin"

	"github.com/gogf/gf/errors/gerror"
)

var config *discoverConfig

type discoverConfig struct {
	Type string
	Addr []string
	Path string
}

func init() {
	config = &discoverConfig{}
	err := g.Cfg().GetStruct("rpc.discover", config)
	if err != nil {
		panic(gerror.Wrap(err, "rpc discover config error"))
	}

}

// NewClientDiscover 根据server名字创建客户端发现服务配置
func NewClientDiscover(service string) (client.ServiceDiscovery, error) {
	return NewClientDiscoverFromConfig(service, config)
}

// NewClientDiscoverFromConfig 根据server和配置参数创建
func NewClientDiscoverFromConfig(service string, cfg *discoverConfig) (client.ServiceDiscovery, error) {
	switch cfg.Type {
	case "redis":
		return client.NewRedisDiscovery(
			cfg.Path,
			service,
			cfg.Addr,
			&store.Config{
				PersistConnection: true,
			},
		)
	case "consul":
		return client.NewConsulDiscovery(
			cfg.Path,
			service,
			cfg.Addr,
			nil,
		)
	}
	return nil, gerror.Newf("error discover type %s", cfg.Type)
}

// NewServer 创建rpc server服务
func NewServer(addr string, limit ...int64) *server.Server {
	fmt.Println("NewServer", addr, config.Type, config.Addr, config.Path)
	rpcServer := server.NewServer(
		server.WithReadTimeout(time.Second*3),
		server.WithWriteTimeout(time.Second*3),
	)
	//一秒产生1万个令牌，桶最大存储两倍于每秒产生的令牌
	var lt int64 = 10000
	if len(limit) > 0 && limit[0] > 0 {
		lt = limit[0]
	}
	rpcServer.Plugins.Add(plugins.NewInfoPlugin(lt, lt*2))
	rpcServer.Plugins.Add(serverplugin.OpenTracingPlugin{})
	rpcServer.DisableHTTPGateway = true

	switch config.Type {
	case "redis":
		discover := &serverplugin.RedisRegisterPlugin{
			ServiceAddress: "tcp@" + addr,
			RedisServers:   config.Addr,
			BasePath:       config.Path,
			Metrics:        metrics.NewRegistry(),
			UpdateInterval: time.Second * 3,
			Options: &store.Config{
				PersistConnection: true,
			},
		}
		err := discover.Start()
		if err != nil {
			panic(err)
		}
		rpcServer.Plugins.Add(discover)
	case "consul":
		discover := &serverplugin.ConsulRegisterPlugin{
			ServiceAddress: "tcp@" + addr,
			ConsulServers:  config.Addr,
			BasePath:       config.Path,
			Metrics:        metrics.DefaultRegistry,
			UpdateInterval: time.Second * 10, //这个更新的是Metrics
		}
		err := discover.Start()
		if err != nil {
			panic(err)
		}
		rpcServer.Plugins.Add(discover)
		//startMetrics()
	default:
		panic(gerror.Newf("error discover type %s", config.Type))
	}

	return rpcServer
}
