package main

import (
	"github.com/olaola-chat/rbp-functor/cmd/cmd1"

	"github.com/olaola-chat/rbp-library/server/cmd"
)

func main() {
	servers := []interface{}{
		&cmd1.CmdCmd01{},
	}

	cmd.Run(servers)
}
