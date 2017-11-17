package nvimutil

import "github.com/neovim/go-client/nvim"

type Nvimutil struct {
	v *nvim.Nvim
}

func New(v *nvim.Nvim) *Nvimutil {
	return &Nvimutil{
		v: v,
	}
}

func (n *Nvimutil) Log(message string) error {
	return n.v.Command("echom 'aoj.nvim: " + message + "'")
}
