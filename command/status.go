package command

import (
	"errors"
	"fmt"

	"github.com/neovim/go-client/nvim"
)

// Vim-Command definition:
// 既に提出したバッファであれば、その一番最近の結果を返す
func (a *AOJ) Status(v *nvim.Nvim, args []string) error {
	defer a.panicLog(v)

	nimvle := nimvleNew(v)

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
		scratch, err = nimvle.NewScratchBuffer("AOJ Status")
		a.ScratchBuffer = scratch
	} else {
		scratch = a.ScratchBuffer
	}

	nimvle.SetContentToBuffer(*scratch, stat.String())

	nimvle.Log(mes)

	return nil
}

// Status用ScratchBufferにStatusのリストを吐く.(diesirae上で参照するStatus番号も)
//  e.g.
//   0 - 2017/01/01 12:00:12 ITP1_1_A
//   1 - 2017/01/01 12:00:21 ITP1_1_A
//   2 - 2017/01/01 12:01:26 ITP1_1_A
//   3 - 2017/01/01 12:01:34 ITP1_1_A
//   ...
func (a *AOJ) StatusList(v *nvim.Nvim, args []string) error {
	defer a.panicLog(v)

	nimvle := nimvleNew(v)

	var output string
	for b, status := range a.SubmittedStatuses {
		for i := 0; i < len(status); i++ {
			output += fmt.Sprintf("%v - %v %s\n", b, status[i].Time, status[i].ProblemId)
			output += fmt.Sprintf("           %v\n", status[i].String())
		}
	}

	if a.ScratchBuffer == nil {
		// NOTE: ScratchBufferがないなら、未提出扱いにする
		return errors.New("not yet submitted")
	}

	err := nimvle.SetContentToBuffer(*a.ScratchBuffer, output)
	if err != nil {
		return err
	}

	return nil
}
