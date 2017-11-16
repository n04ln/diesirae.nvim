package command

import (
	"github.com/NoahOrberg/aoj.nvim/util"
	"github.com/neovim/go-client/nvim"
)

func (a *AOJ) Session(v *nvim.Nvim, args []string) error {
	// TODO: 引数で、ユーザ名、パスワードをもらって、セッションを得る

	if len(args) != 2 {
		return util.ErrInvalidArgs
	}

	return nil
}

func (a *AOJ) Self(v *nvim.Nvim, args []string) error {
	// TODO: セッションが生きているかどうかの確認
	return nil
}
