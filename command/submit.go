package command

import (
	"errors"
	"net/url"

	"github.com/NoahOrberg/diesirae.nvim/aoj"
	"github.com/NoahOrberg/diesirae.nvim/nvimutil"
	"github.com/neovim/go-client/nvim"
)

// Vim-Function definition:
//   第一引数で問題のタイトルを指定する。
func (a *AOJ) SubmitAndCheckStatus(v *nvim.Nvim, args []string) error {
	if len(args) != 1 {
		return errors.New("invalid args")
	}

	if a.IsValidCookie == false {
		return errors.New("you should execute :AojSession")
	}
	defer a.panicLog(v)

	nvimutil := nvimutil.New(v)

	input := args[0]
	var problemId string
	// ここでは、URLでくるか、問題の題名だけでくるか、両方を受容する
	// TODO: 変更される余地ありかもなので、ここは要観察。現行版のAOJはid=XXXXでクエリパラメータ渡してるのでいいが、他の場合は要修正。
	if u, err := url.ParseRequestURI(input); err != nil {
		problemId = input
	} else {
		ids, ok := u.Query()["id"]
		if !ok || len(ids) == 0 {
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

	// よしなにScratchBufferに表示
	if err := a.showScratchBuffer(nvimutil, stat); err != nil {
		return err
	}

	return nil
}
