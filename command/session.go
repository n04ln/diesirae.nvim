package command

import (
	"github.com/NoahOrberg/aoj.nvim/config"
	"github.com/h2non/gentleman"
	"github.com/neovim/go-client/nvim"
)

// NOTE: セッションが生きているかどうかの確認
func (a *AOJ) Self(v *nvim.Nvim, args []string) error {
	conf := config.GetConfig()

	cli := gentleman.New()
	cli.URL(conf.Endpoint)

	req := cli.Request()
	req.Method("GET")
	req.Path("/self")
	req.SetHeader("Cookie", a.Cookie)

	res, err := req.Send()
	if err != nil {
		return err
	}
	if !res.Ok {
		v.Command("echom 'session not exists'")
		return nil
	}

	v.Command("echom 'session exists'")
	return nil
}
