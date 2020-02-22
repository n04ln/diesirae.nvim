package command

import (
	"errors"
	"net/url"

	"github.com/n04ln/diesirae.nvim/aoj"
	"github.com/n04ln/diesirae.nvim/config"
	"github.com/n04ln/diesirae.nvim/util"
	"github.com/neovim/go-client/nvim"
)

type CompileError struct {
	s string
}

func (c *CompileError) String() string {
	return "CompileError:\n" + c.s
}

// Vim-Function definition:
func (a *AOJ) Trial(v *nvim.Nvim, args []string) (err error) {
	defer a.panicLog(v)
	if len(args) > 0 && args[0] == "" {
		return nil
	}

	nimvle := nimvleNew(v)
	done = make(chan string)

	if a.ScratchBuffer == nil {
		a.ScratchBuffer, err = nimvle.NewScratchBuffer(config.GetConfig().ResultBufferName)
		if err != nil {
			nimvle.Log(err.Error())
			return
		}
	}

	go drawLoadingCycle(nimvle, a.ScratchBuffer)

	if len(args) != 1 {
		err = util.ErrInvalidArgs
		nimvle.Log(err.Error())
		flushLoadingCycle(nimvle, a.ScratchBuffer, err)
		return
	}

	if a.IsValidCookie == false {
		err = util.ErrInvalidCookie
		nimvle.Log(err.Error())
		flushLoadingCycle(nimvle, a.ScratchBuffer, err)
		return
	}

	input := args[0]
	var problemId string
	// ここでは、URLでくるか、問題の題名だけでくるか、両方を受容する
	// TODO: 変更される余地ありかもなので、ここは要観察。現行版のAOJはid=XXXXでクエリパラメータ渡してるのでいいが、他の場合は要修正。
	if u, err2 := url.ParseRequestURI(input); err2 != nil {
		problemId = input
	} else {
		ids, ok := u.Query()["id"]
		if !ok || len(ids) == 0 {
			err = errors.New("no such id")
			nimvle.Log(err.Error())
			flushLoadingCycle(nimvle, a.ScratchBuffer, err)
			return
		}

		problemId = ids[0]
	}

	extension, err := nimvle.CurrentBufferFilenameExtension()
	if err != nil {
		nimvle.Log(err.Error())
		flushLoadingCycle(nimvle, a.ScratchBuffer, err)
		return
	}

	sourceCode, err := nimvle.GetContentFromCurrentBuffer()
	if err != nil {
		nimvle.Log(err.Error())
		flushLoadingCycle(nimvle, a.ScratchBuffer, err)
		return
	}

	// 現状、無限に返ってこない可能性があるため並列処理に回す
	go func() {
		defer a.panicLog(v)

		// sampleコード表示
		samples, err := aoj.GetSampleInputOutput(problemId)
		if err != nil {
			nimvle.Log(err.Error())
			flushLoadingCycle(nimvle, a.ScratchBuffer, err)
			return
		}

		// description を取得し、時間を入手
		desc, err := aoj.GetDescriptionByProblemId(a.Cookie, problemId)
		if err != nil {
			nimvle.Log(err.Error())
			flushLoadingCycle(nimvle, a.ScratchBuffer, err)
			return
		}

		// 実行
		vsConfig, err := config.GetVSConfig(v)
		buildCommand, execCommand := make([]string, 0), make([]string, 0)
		if cm, ok := vsConfig.Commands[extension]; ok {
			buildCommand = cm.BuildCommand
			execCommand = cm.ExecCommand
		} else if len(buildCommand) == 0 && len(execCommand) == 0 {
			err := errors.New("no such build config of " + extension + " language in g:" + util.ConfigVarName)
			nimvle.Log(err.Error())
			flushLoadingCycle(nimvle, a.ScratchBuffer, err)
			return
		}
		if output, err := samples.ExecSamples(extension, sourceCode,
			buildCommand, execCommand, desc.TimeLimit); err != nil {
			if err == aoj.ErrCompileError {
				if output == nil {
					output = new(string)
					*output = "unexpected error"
				}

				// コンパイルエラーなので、その旨をScratchBuffer経由で報告する
				err = nimvle.ShowScratchBuffer(*a.ScratchBuffer, &CompileError{s: *output})
				if err != nil {
					nimvle.Log(err.Error())
					flushLoadingCycle(nimvle, a.ScratchBuffer, err)
					return
				}
			}

			done <- ""
			return
		}

		done <- ""
		err = nimvle.ShowScratchBuffer(*a.ScratchBuffer, samples)
		if err != nil {
			nimvle.Log(err.Error())
			flushLoadingCycle(nimvle, a.ScratchBuffer, err)
			return
		}
	}()

	return nil
}
