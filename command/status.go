package command

import "github.com/neovim/go-client/nvim"

func (a *AOJ) Status(v *nvim.Nvim, args []string) error {
	// TODO: 提出情報から、CurrentBufferのTokenを抜き、それがACかWAかなど見る。
	// NOTE: 提出情報はa.SubmittedTokensに格納するようにする。詳しい情報などを見るために、WEBページのリンクも載せる
	return nil
}
