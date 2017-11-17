package nvimutil

import (
	"fmt"
	"strings"

	"github.com/neovim/go-client/nvim"
)

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

func (n *Nvimutil) CurrentBufferFileType() (string, error) {
	buf, err := n.v.CurrentBuffer()
	if err != nil {
		return "", err
	}

	bufferName, err := n.v.BufferName(buf)
	if err != nil {
		return "", err
	}

	dotName := strings.Split(bufferName, ".")[len(strings.Split(bufferName, "."))-1]

	var language string
	switch dotName {
	case "c":
		language = "C"
	case "hs": // NOTE: maybe not correct return "Haskell"? need fix
		language = "Haskell"
	default:
		return "", fmt.Errorf("cannot submit this file: .%s", dotName)
	}

	return language, nil
}
