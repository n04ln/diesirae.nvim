package command

import (
	"errors"
	"fmt"
	"net/url"

	"github.com/NoahOrberg/diesirae.nvim/aoj"
	"github.com/NoahOrberg/diesirae.nvim/config"
	"github.com/NoahOrberg/diesirae.nvim/util"
	"github.com/neovim/go-client/nvim"
)

func changeExtToLanguageName(ex string) (string, error) {
	var language string
	switch ex {
	// Languages: C, Haskell, Go, C++14, JAVA, C#, D, Go, Rust, Ruby,
	//            Python3, JavaScript, Haskell, Scala, PHP, OCaml, Kotlin
	case "c":
		language = "C"
	case "hs":
		language = "Haskell"
	case "go":
		language = "Go"
	case "cpp":
		language = "C++14"
	case "java":
		language = "JAVA"
	case "cs":
		language = "C#"
	case "d":
		language = "D"
	case "rs":
		language = "Rust"
	case "rb":
		language = "Ruby"
	case "py":
		language = "Python3"
	case "js":
		language = "JavaScript"
	case "scala":
		language = "Scala"
	case "php":
		language = "PHP"
	case "ml":
		language = "OCaml"
	case "kt":
		language = "Kotlin"
	default:
		if ex == "" {
			return "", fmt.Errorf("cannot identify file type")
		}
		return "", fmt.Errorf("cannot submit this file: .%s", ex)
	}

	return language, nil
}

// Vim-Function definition:
//   第一引数で問題のタイトルを指定する。
func (a *AOJ) SubmitAndCheckStatus(v *nvim.Nvim, args []string) error {
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

	extension, err := nimvle.CurrentBufferFilenameExtension()
	if err != nil {
		nimvle.Log(err.Error())
		return err
	}

	language, err := changeExtToLanguageName(extension)
	if err != nil {
		nimvle.Log(err.Error())
		return err
	}

	sourceCode, err := nimvle.GetContentFromCurrentBuffer()
	if err != nil {
		nimvle.Log(err.Error())
		return err
	}

	token, err := aoj.Submit(a.Cookie, problemId, language, sourceCode)
	if err != nil {
		nimvle.Log(err.Error())
		return err
	}

	stat, err := aoj.Status(a.Cookie, token, problemId)
	if err != nil {
		return err
	}

	buf, err := v.CurrentBuffer()
	if err != nil {
		nimvle.Log(err.Error())
		return err
	}

	a.SetStatusByBuffer(buf, stat)

	return nimvle.ShowScratchBuffer(*a.ScratchBuffer, stat)
}
