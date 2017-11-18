package nvimutil

import "github.com/neovim/go-client/nvim"

func Log(v *nvim.Nvim, message string) error {
	return v.Command("echom 'diesirae.nvim: " + message + "'")
}
