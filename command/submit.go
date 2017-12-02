package command

import (
	"github.com/NoahOrberg/diesirae.nvim/aoj"
	"github.com/NoahOrberg/diesirae.nvim/nvimutil"
	"github.com/NoahOrberg/diesirae.nvim/util"
	"github.com/neovim/go-client/nvim"
)

// Vim-Command definition:
// Exコマンドの第一引数で問題のタイトルを指定する。
func (a *AOJ) SubmitAndCheckStatus(v *nvim.Nvim, args []string) error {
	defer a.panicLog(v)

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

	stat, err := aoj.Status(a.Cookie, token, problemId)
	if err != nil {
		return err
	}

	mes, err := stat.CheckAC()
	if err != nil {
		return err
	}

	buf, err := v.CurrentBuffer()
	if err != nil {
		return err
	}

	a.SetStatusByBuffer(buf, stat)

	var scratch *nvim.Buffer
	if a.ScratchBuffer == nil {
		scratch, err = nvimutil.NewScratchBuffer("AOJ Status")
		a.ScratchBuffer = scratch
	} else {
		scratch = a.ScratchBuffer
	}

	nvimutil.SetContentToBuffer(*scratch, stat.String())

	nvimutil.Log(mes)

	return nil
}
