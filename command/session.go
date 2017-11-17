package command

import (
	"github.com/NoahOrberg/aoj.nvim/aoj"
	"github.com/NoahOrberg/aoj.nvim/config"
	"github.com/NoahOrberg/aoj.nvim/nvimutil"
	"github.com/neovim/go-client/nvim"
)

// セッションが生きているかどうかの確認
func (a *AOJ) Self(v *nvim.Nvim, args []string) error {
	nvimutil := nvimutil.New(v)
	if ok := aoj.IsAliveSession(a.Cookie); !ok {
		nvimutil.Log("session not exists")
		return nil
	}
	nvimutil.Log("session exists")
	return nil
}

// セッションを張り直す
func (a *AOJ) Session(v *nvim.Nvim, args []string) error {
	nvimutil := nvimutil.New(v)

	if ok := aoj.IsAliveSession(a.Cookie); ok {
		nvimutil.Log("session exists")
		return nil
	}

	if cookie, err := reconnectSession(); err != nil {
		a.Cookie = cookie
		nvimutil.Log("session reconnect!")
		return nil
	}

	return nil
}

func reconnectSession() (string, error) {
	conf := config.GetConfig()

	return aoj.Session(conf.ID, conf.RawPassword)
}
