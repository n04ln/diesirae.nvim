package command

import (
	"errors"
	"net/url"

	"github.com/NoahOrberg/diesirae.nvim/aoj"
	"github.com/NoahOrberg/diesirae.nvim/config"
	"github.com/NoahOrberg/diesirae.nvim/nvimutil"
	"github.com/neovim/go-client/nvim"
)

// Vim-Command definition:
// Exコマンドの第一引数で問題のタイトルを指定する。
func (a *AOJ) SubmitAndCheckStatus(v *nvim.Nvim, args []string) error {
	defer a.panicLog(v)

	nvimutil := nvimutil.New(v)

	var problemId string
	input, err := nvimutil.Input("problem id")
	// ここでは、URLでくるか、問題の題名だけでくるか、両方を受容する
	// TODO: 変更される余地ありかもなので、ここは要観察。現行版のAOJはid=XXXXでクエリパラメータ渡してるのでいいが、他の場合は要修正。
	if u, err := url.ParseRequestURI(input); err != nil {
		problemId = input
	} else {
		ids, ok := u.Query()["id"]
		if !ok {
			return errors.New("no such id")
		}

		problemId = ids[0]
	}

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

	buf, err := v.CurrentBuffer()
	if err != nil {
		return err
	}

	a.SetStatusByBuffer(buf, stat)

	var opened bool
	var scratch *nvim.Buffer
	conf := config.GetConfig()
	if a.ScratchBuffer == nil {
		scratch, err = nvimutil.NewScratchBuffer(conf.ResultBufferName)
		a.ScratchBuffer = scratch
		opened = true
	} else {
		scratch = a.ScratchBuffer
	}

	nvimutil.SetContentToBuffer(*scratch, stat.String())

	winls, err := nvimutil.GetWindowList()
	if err != nil {
		return err
	}

	if !opened {
		for _, bufname := range winls {
			if bufname == conf.ResultBufferName {
				opened = true
				break
			}
		}
	}

	if !opened {
		nvimutil.SplitOpenBuffer(*scratch)
	}

	return nil
}
