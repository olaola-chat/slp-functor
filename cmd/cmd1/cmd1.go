package cmd1

import "github.com/olaola-chat/slp-functor/cmd/cmd1/internal/logic"

type CmdCmd01 struct {
}

func (c *CmdCmd01) Method1() {
	m := &logic.Method1{}
	m.Run()
}
