package command

import (
	"errors"

	"github.com/NoahOrberg/aoj.nvim/nvimutil"
	"github.com/neovim/go-client/nvim"
)

func (a *AOJ) Status(v *nvim.Nvim, args []string) error {
	buf, err := v.CurrentBuffer()
	if err != nil {
		return err
	}

	stat, ok := a.GetRecentStatusByBuffer(buf)
	if !ok {
		return errors.New("not submited this buffer yet")
	}

	mes, err := stat.CheckAC()
	if err != nil {
		return err
	}

	nvimutil := nvimutil.New(v)

	nvimutil.Log(mes)

	return nil
}
