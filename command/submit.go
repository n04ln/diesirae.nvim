package command

import (
	"errors"

	"github.com/neovim/go-client/nvim"
)

var ErrInvalidArgs = errors.New("invalid args")

func (a *AOJ) Submit(v *nvim.Nvim, args []string) error {
	// NOTE: args[0] -> problem id
	//       args[1] -> language
	//       args[2] -> ...
	// TODO: 提出用コマンド: CurrentBufferを、現在のファイルタイプから見て投げる(その場合上のlanguageいらなくなる)

	if len(args) != 2 {
		return ErrInvalidArgs
	}
	v.Command("echo '" + args[0] + " " + args[1] + "'")
	return nil
}
