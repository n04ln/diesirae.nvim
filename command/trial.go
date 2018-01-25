package command

import (
	"errors"
	"net/url"

	"github.com/NoahOrberg/diesirae.nvim/aoj"
	"github.com/NoahOrberg/diesirae.nvim/config"
	"github.com/NoahOrberg/diesirae.nvim/util"
	"github.com/neovim/go-client/nvim"
)

type CompileError struct {
	s string
}

func (c *CompileError) String() string {
	return "CompileError:\n" + c.s
}

// Vim-Function definition:
func (a *AOJ) Trial(v *nvim.Nvim, args []string) error {
	defer a.panicLog(v)
	nimvle := nimvleNew(v)

	if len(args) != 1 {
		nimvle.Log(util.ErrInvalidArgs.Error())
		return errors.New("invalid args")
	}

	if a.IsValidCookie == false {
		nimvle.Log(util.ErrInvalidCookie.Error())
		return errors.New("you should execute :AojSession")
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
			return errors.New("no such id")
		}

		problemId = ids[0]
	}

	extension, err := nimvle.CurrentBufferFilenameExtension()
	if err != nil {
		return err
	}
	fileType, err := transLanguage(extension)
	if err != nil {
		return err
	}

	sourceCode, err := nimvle.GetContentFromCurrentBuffer()
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
			return nimvle.ShowScratchBuffer(*a.ScratchBuffer, &CompileError{s: *output})
		}

		return err
	}

	// よしなにScratchBufferに表示
	if a.ScratchBuffer == nil {
		a.ScratchBuffer, err = nimvle.NewScratchBuffer(config.GetConfig().ResultBufferName)
		if err != nil {
			return err
		}
	}
	return nimvle.ShowScratchBuffer(*a.ScratchBuffer, samples)
}
