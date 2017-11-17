package command

import (
	"github.com/NoahOrberg/aoj.nvim/config"
	"github.com/NoahOrberg/aoj.nvim/nvimutil"
	"github.com/h2non/gentleman"
	"github.com/neovim/go-client/nvim"
)

// セッションが生きているかどうかの確認
func (a *AOJ) Self(v *nvim.Nvim, args []string) error {
	if ok := isAliveSession(a.Cookie); !ok {
		nvimutil.Log("session not exists")
		return nil
	}
	nvimutil.Log("session exists")
	return nil
}

// セッションを張り直す
func (a *AOJ) Session(v *nvim.Nvim, args []string) error {
	conf := config.GetConfig()

	if ok := isAliveSession(a.Cookie); ok {
		nvimutil.Log("session exists")
		return nil
	}

	if cookie, err := reconnectSession(); err != nil {
		a.Cookie = cookie
		nvimutil.Log("session reconnect!")
		return nil
	}
}

func isAliveSession(cookie string) bool {
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
		return false
	}

	return true
}

func reconnectSession() (string, error) {
	conf := config.GetConfig()

	return Session(conf.ID, conf.RawPassword)
}
