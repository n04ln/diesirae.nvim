package command

import (
	"fmt"

	"github.com/NoahOrberg/aoj.nvim/aoj"
	"github.com/NoahOrberg/aoj.nvim/nvimutil"
	"github.com/NoahOrberg/aoj.nvim/util"
	"github.com/neovim/go-client/nvim"
)

// Exコマンドの第一引数で問題のタイトルを指定する。
func (a *AOJ) Submit(v *nvim.Nvim, args []string) error {
	if len(args) != 1 {
		return util.ErrInvalidArgs
	}

	nvimutil := nvimutil.New(v)

	problemId := args[0]

	language, err := nvimutil.CurrentBufferFileType()
	if err != nil {
		return err
	}

	buf, err := v.CurrentBuffer()
	if err != nil {
		return err
	}

	sourceCode, err := getContentFromBuffer(v, buf)
	if err != nil {
		return err
	}

	token, err := aoj.Submit(a.Cookie, problemId, language, sourceCode)
	if err != nil {
		return err
	}

	res, err := aoj.Status(a.Cookie, token)
	if err != nil {
		return err
	}

	mes, err := checkAC(res)
	if err != nil {
		return err
	}
	nvimutil.Log(mes)

	return nil
}

func checkAC(res *aoj.SubmissionStatus) (string, error) {
	for _, caseVerdict := range res.CaseVerdicts {
		if caseVerdict.Status != "AC" {
			return fmt.Sprintf("testcase %s: %s", caseVerdict.Label, caseVerdict.Status), nil
		}
	}
	return fmt.Sprintf("All case: AC"), nil
}

func getContentFromBuffer(v *nvim.Nvim, buf nvim.Buffer) (string, error) {
	lines, err := v.BufferLines(buf, 0, -1, true)
	if err != nil {
		return "", err
	}

	var content string
	for i, c := range lines {
		content += string(c)
		if i < len(lines)-1 {
			content += "\n"
		}
	}

	return content, nil
}
