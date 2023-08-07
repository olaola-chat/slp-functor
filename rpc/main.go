package main

import (
	"github.com/olaola-chat/rbp-functor/rpc/act1"

	"github.com/olaola-chat/rbp-library/server/rpc"
)

func main() {
	servers := map[string]*rpc.ServerCfg{
		"act1": {
			RegisterName: "Activity.Act1",
			Server:       act1.NewAct1,
		},
	}

	rpc.Run(servers)
}
