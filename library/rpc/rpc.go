package rpc

import (
	"fmt"
	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/frame/g"
	"github.com/rpcxio/libkv/store"
	"github.com/smallnest/rpcx/client"
	"os"
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
	if config.Type == "consul" && len(config.Addr) == 0 {
		consulAgentIp := os.Getenv("CONSUL_AGENT_IP")
		if consulAgentIp == "" {
			panic(gerror.Wrap(err, "rpc discover config error"))
		}
		config.Addr = []string{fmt.Sprintf("%s:%d", consulAgentIp, 8500)}
	}
	g.Log().Infof("found rpc discover config, type = %s , path = %s , addrs = %v",
		config.Type, config.Path, config.Addr)
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
