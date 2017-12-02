package command

import (
	"errors"

	"github.com/NoahOrberg/diesirae.nvim/nvimutil"
	"github.com/neovim/go-client/nvim"
)

// 既に提出したバッファであれば、その一番最近の結果を返す
func (a *AOJ) Status(v *nvim.Nvim, args []string) error {
	defer a.panicLog(v)

	nvimutil := nvimutil.New(v)

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
