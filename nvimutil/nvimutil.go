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
		if dotName == "" {
			return "", fmt.Errorf("cannot identify file type")
		}
		return "", fmt.Errorf("cannot submit this file: .%s", dotName)
	}

	return language, nil
}

func (n *Nvimutil) GetContentFromCurrentBuffer() (string, error) {
	buf, err := n.v.CurrentBuffer()
	if err != nil {
		return "", err
	}

	lines, err := n.v.BufferLines(buf, 0, -1, true)
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

func (n *Nvimutil) GetContentFromBuffer(buf nvim.Buffer) (string, error) {
	lines, err := n.v.BufferLines(buf, 0, -1, true)
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

func (n *Nvimutil) SetContentToBuffer(buf nvim.Buffer, content string) error {
	var byteContent [][]byte

	tmp := strings.Split(content, "\n")
	for _, c := range tmp {
		byteContent = append(byteContent, []byte(c))
	}

	return n.v.SetBufferLines(buf, 0, -1, true, byteContent)
}

func (n *Nvimutil) NewScratchBuffer() (*nvim.Buffer, error) {
	bwin, err := n.v.CurrentWindow()
	if err != nil {
		return nil, err
	}

	if err := n.v.Command("new"); err != nil {
		return nil, err
	}

	scratchBuf, err := n.v.CurrentBuffer()
	if err != nil {
		return nil, err
	}

	if err := n.v.SetBufferOption(scratchBuf, "buftype", "nofile"); err != nil {
		return nil, err
	}

	if err := n.v.Command("setlocal noswapfile"); err != nil {
		return nil, err
	}

	if err := n.v.SetBufferOption(scratchBuf, "undolevels", -1); err != nil {
		return nil, err
	}

	if err := n.v.SetBufferName(scratchBuf, "Aizu Online Judge"); err != nil {
		return nil, err
	}

	win, err := n.v.CurrentWindow()
	if err != nil {
		return nil, err
	}

	if err := n.v.SetWindowHeight(win, 15); err != nil {
		return nil, err
	}

	if err := n.v.SetCurrentWindow(bwin); err != nil {
		return nil, err
	}

	return &scratchBuf, nil
}
