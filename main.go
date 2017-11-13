package main

import (
	"github.com/NoahOrberg/aoj.nvim/command"
	"github.com/neovim/go-client/nvim/plugin"
)

func main() {
	a := &command.AOJ{}
	plugin.Main(func(p *plugin.Plugin) error {
		p.HandleCommand(&plugin.CommandOptions{Name: "AojSubmit", NArgs: "*"}, a.Submit)
		return nil
	})
}
