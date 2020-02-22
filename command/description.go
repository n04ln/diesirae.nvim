package command

import (
	"errors"
	"net/url"

	"github.com/n04ln/diesirae.nvim/aoj"
	"github.com/n04ln/diesirae.nvim/config"
	"github.com/n04ln/diesirae.nvim/util"
	"github.com/neovim/go-client/nvim"
)

// Vim-Command definition:
// 問題概要取得
func (a *AOJ) Description(v *nvim.Nvim, args []string) error {
	nimvle := nimvleNew(v)

	if len(args) != 1 {
		nimvle.Log(util.ErrInvalidArgs.Error())
		return errors.New("invalid args")
	}

	if a.IsValidCookie == false {
		nimvle.Log(util.ErrInvalidCookie.Error())
		return errors.New("you should execute :AojSession")
	}
	defer a.panicLog(v)

	if a.ScratchBuffer == nil {
		var err error
		a.ScratchBuffer, err = nimvle.NewScratchBuffer(config.GetConfig().ResultBufferName)
		if err != nil {
			nimvle.Log(err.Error())
			return err
		}
	}

	input := args[0]
	var problemId string
	// ここでは、URLでくるか、問題の題名だけでくるか、両方を受容する
	// TODO: 変更される余地ありかもなので、ここは要観察。現行版のAOJはid=XXXXでクエリパラメータ渡してるのでいいが、他の場合は要修正。
	if u, err := url.ParseRequestURI(input); err != nil {
		problemId = input
	} else {
		ids, ok := u.Query()["id"]
		if !ok || len(ids) == 0 {
			nimvle.Log(err.Error())
			return errors.New("no such id")
		}

		problemId = ids[0]
	}

	desc, err := aoj.GetDescriptionByProblemId(a.Cookie, problemId)
	if err != nil {
		nimvle.Log(err.Error())
		return err
	}

	return nimvle.ShowScratchBuffer(*a.ScratchBuffer, desc)
}
