package rpc

import (
	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/frame/g"
	"github.com/rpcxio/libkv/store"
	"github.com/smallnest/rpcx/client"
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
