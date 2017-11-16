package command

import (
	"github.com/NoahOrberg/aoj.nvim/aoj"
	"github.com/NoahOrberg/aoj.nvim/config"
)

type AOJ struct {
	Cookie          string           // NOTE: ログインしたら、ここにクッキーをいれる
	Config          config.AOJConfig // NOTE: 環境変数から取得した情報格納
	SubmittedTokens map[int]string   // NOTE: 提出したとき、あとからそれを確認できる用にするため、キーをバッファ番号にして確認用Tokenを保存する
}

func NewAOJ() (*AOJ, error) {
	conf := config.GetConfig()

	cookie, err := aoj.Session(conf.ID, conf.RawPassword)
	if err != nil {
		return nil, err
	}

	return &AOJ{
		Cookie: cookie,
		Config: conf,
	}, nil
}
