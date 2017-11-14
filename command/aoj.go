package command

type AOJ struct {
	Cookie          string         // NOTE: ログインしたら、ここにクッキーをいれる
	Config          config.config  // NOTE: 環境変数から取得した情報格納
	SubmittedTokens map[int]string // NOTE: 提出したとき、あとからそれを確認できる用にするため、キーをバッファ番号にして確認用Tokenを保存する
}

func NewAOJ() *AOJ {
	// TODO: Configのセット、その他諸々の初期設定
	// NOTE: ログインは提出前にやるようにするので、ここではSessionのCookieセットは行わない
}
