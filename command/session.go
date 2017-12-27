package command

import (
	"errors"

	"github.com/NoahOrberg/diesirae.nvim/aoj"
	"github.com/NoahOrberg/diesirae.nvim/config"
	"github.com/NoahOrberg/diesirae.nvim/nvimutil"
	"github.com/neovim/go-client/nvim"
)

// Vim-Command definition:
// セッションが生きているかどうかの確認
func (a *AOJ) Self(v *nvim.Nvim, args []string) error {
	defer a.panicLog(v)

	nvimutil := nvimutil.New(v)

	if ok := aoj.IsAliveSession(a.Cookie); !ok {
		nvimutil.Log("session not exists. you should execute :AojSession")
		a.IsValidCookie = false
		return nil
	}
	nvimutil.Log("session exists")

	return nil
}

// Vim-Command definition:
// セッションを張り直す
func (a *AOJ) Session(v *nvim.Nvim, args []string) error {
	defer a.panicLog(v)

	nvimutil := nvimutil.New(v)

	if ok := aoj.IsAliveSession(a.Cookie); ok {
		nvimutil.Log("session exists")
		return nil
	}

	cookie, err := reconnectSession()
	if err != nil {
		a.IsValidCookie = false
		nvimutil.Log("session cannot reconnect...")
		return nil
	}

	a.Cookie = cookie
	a.IsValidCookie = true
	nvimutil.Log("session reconnect!")
	return nil
}

func reconnectSession() (string, error) {
	conf := config.GetConfig()

	cookie, err := aoj.Session(conf.ID, conf.RawPassword)
	if err != nil {
		return "", err
	}
	if len(cookie) == 0 {
		return "", errors.New("session is nothing")
	}

	return cookie, nil
}
