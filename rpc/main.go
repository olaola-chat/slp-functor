package main

import (
	"github.com/olaola-chat/rbp-functor/rpc/act1"
	"github.com/olaola-chat/rbp-functor/rpc/voice_lover"

	"github.com/olaola-chat/rbp-library/server/rpc"
)

func main() {
	servers := map[string]*rpc.ServerCfg{
		"act1": {
			RegisterName: "Activity.Act1",
			Server:       act1.NewAct1,
		},
		"voice_lover_main": {
			RegisterName: "Functor.VoiceLover.Main",
			Server:       voice_lover.NewVoiceLoverMain,
		},
		"voice_lover_admin": {
			RegisterName: "Functor.VoiceLover.Admin",
			Server:       voice_lover.NewVoiceLoverAdmin,
		},
	}

	rpc.Run(servers)
}
