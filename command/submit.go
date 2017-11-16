package command

import (
	"github.com/NoahOrberg/aoj.nvim/aoj"
	"github.com/NoahOrberg/aoj.nvim/util"
	"github.com/neovim/go-client/nvim"
)

func (a *AOJ) Submit(v *nvim.Nvim, args []string) error {
	// NOTE: args[0] -> problem id
	//       args[1] -> language
	//       args[2] -> ...
	// TODO: 提出用コマンド: CurrentBufferを、現在のファイルタイプから見て投げる(その場合上のlanguageいらなくなる)

	v.Command("echo '" + a.Cookie + "'")

	if len(args) != 2 {
		return util.ErrInvalidArgs
	}

	problemId := "ITP1_1_A"
	language := "C"
	sourceCode := "#include \nint main(){\n printf(\"Hello World\\n\");\n return 0;\n}"

	v.Command("echom '" + a.Config.ID + "'")
	v.Command("echom '" + a.Config.Endpoint + "'")

	token, err := aoj.Submit(a.Cookie, problemId, language, sourceCode)
	if err != nil {
		return err
	}

	v.Command("echo '" + token + "'")

	return nil
}
