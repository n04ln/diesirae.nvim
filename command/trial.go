package command

import (
	"errors"
	"net/url"

	"github.com/NoahOrberg/diesirae.nvim/aoj"
	"github.com/NoahOrberg/diesirae.nvim/nvimutil"
	"github.com/neovim/go-client/nvim"
)

type CompileError struct {
	s string
}

func (c *CompileError) String() string {
	return "CompileError:\n" + c.s
}

// Vim-Command definition:
func (a *AOJ) Trial(v *nvim.Nvim, args []string) error {
	defer a.panicLog(v)

	nvimutil := nvimutil.New(v)

	var problemId string
	input, err := nvimutil.Input("problem id")
	if input == "" {
		return nil
	}
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

	fileType, err := nvimutil.CurrentBufferFileType()
	if err != nil {
		return err
	}

	sourceCode, err := nvimutil.GetContentFromCurrentBuffer()
	if err != nil {
		return err
	}

	// sampleコード表示
	samples, err := aoj.GetSampleInputOutput(problemId)
	if err != nil {
		return err
	}

	// 実行
	if output, err := samples.ExecSamples(fileType, sourceCode); err != nil {
		if err == aoj.ErrCompileError {
			if output == nil {
				*output = "unexpected error"
			}

			// コンパイルエラーなので、その旨をScratchBuffer経由で報告する
			if err := a.showScratchBuffer(nvimutil, &CompileError{s: *output}); err != nil {
				return err
			}
			return nil
		}

		return err
	}

	// よしなにScratchBufferに表示
	if err := a.showScratchBuffer(nvimutil, samples); err != nil {
		return err
	}

	return nil
}
