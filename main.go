package main

import (
	"github.com/NoahOrberg/diesirae.nvim/command"
	"github.com/neovim/go-client/nvim/plugin"
)

func main() {
	a, err := command.NewAOJ()
	if err != nil {
		panic("cannot use session")
	}

	plugin.Main(func(p *plugin.Plugin) error {
		p.HandleCommand(&plugin.CommandOptions{Name: "AojSubmit", NArgs: "+"}, a.SubmitAndCheckStatus)
		p.HandleCommand(&plugin.CommandOptions{Name: "AojStatus"}, a.Status)
		p.HandleCommand(&plugin.CommandOptions{Name: "AojSelf"}, a.Self)
		p.HandleCommand(&plugin.CommandOptions{Name: "AojSession"}, a.Session)
		p.HandleCommand(&plugin.CommandOptions{Name: "AojStatusList"}, a.StatusList)
		return nil
	})
}
