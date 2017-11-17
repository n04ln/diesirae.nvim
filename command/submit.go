package command

import (
	"github.com/NoahOrberg/aoj.nvim/aoj"
	"github.com/NoahOrberg/aoj.nvim/nvimutil"
	"github.com/NoahOrberg/aoj.nvim/util"
	"github.com/neovim/go-client/nvim"
)

// Exコマンドの第一引数で問題のタイトルを指定する。
func (a *AOJ) Submit(v *nvim.Nvim, args []string) error {
	if len(args) != 1 {
		return util.ErrInvalidArgs
	}

	nvimutil := nvimutil.New(v)

	problemId := args[0]

	language, err := nvimutil.CurrentBufferFileType()
	if err != nil {
		return err
	}

	sourceCode, err := nvimutil.GetContentFromCurrentBuffer()
	if err != nil {
		return err
	}

	token, err := aoj.Submit(a.Cookie, problemId, language, sourceCode)
	if err != nil {
		return err
	}

	res, err := aoj.Status(a.Cookie, token)
	if err != nil {
		return err
	}

	mes, err := res.CheckAC()
	if err != nil {
		return err
	}
	nvimutil.Log(mes)

	return nil
}
