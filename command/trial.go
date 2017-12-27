package command

import (
	"errors"
	"io"
	"io/ioutil"
	"net/url"
	"os"
	"os/exec"

	"github.com/NoahOrberg/diesirae.nvim/aoj"
	"github.com/NoahOrberg/diesirae.nvim/config"
	"github.com/NoahOrberg/diesirae.nvim/nvimutil"
	"github.com/neovim/go-client/nvim"
)

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
	// TODO: とりあえずGoだけ
	switch fileType {
	case "Go":
		fp, err := ioutil.TempFile("", "diesirae")
		if err != nil {
			return err
		}
		defer os.Remove(fp.Name())
		defer fp.Close()
		defer os.Remove(fp.Name() + ".go")

		if err := os.Rename(fp.Name(), fp.Name()+".go"); err != nil {
			return err
		}

		if _, err := fp.Write([]byte(sourceCode)); err != nil {
			return err
		}

		for i, sample := range samples.Samples {
			cmd := exec.Command("go", "run", fp.Name()+".go")
			stdin, err := cmd.StdinPipe()
			if err != nil {
				return err
			}
			_, err = io.WriteString(stdin, sample.Input)
			if err != nil {
				return err
			}
			err = stdin.Close()
			if err != nil {
				return err
			}
			out, err := cmd.Output()
			if err != nil {
				return err
			}
			samples.Samples[i].Actual = string(out)
		}
	default:
		return errors.New("only support Golang :)")
	}

	// ScratchBufferを別ウィンドウで開いていればいいが、開かれていない場合などの処理
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

	nvimutil.SetContentToBuffer(*scratch, samples.String())

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
