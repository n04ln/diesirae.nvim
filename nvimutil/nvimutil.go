package nvimutil

import (
	"fmt"
	"strings"

	"github.com/neovim/go-client/nvim"
)

const (
	getWindowList = `diesirae#getWindowList()`
)

type Nvimutil struct {
	v *nvim.Nvim
}

func New(v *nvim.Nvim) *Nvimutil {
	return &Nvimutil{
		v: v,
	}
}

func (n *Nvimutil) Log(message interface{}) error {
	return n.v.Command("echom 'diesirae.nvim: " + fmt.Sprintf("%v", message) + "'")
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
	case "hs":
		language = "Haskell"
	case "go":
		language = "Go"
	case "cpp":
		language = "C++14"
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

func (n *Nvimutil) GetWindowList() (map[string]string, error) {
	res := make(map[string]string)

	if err := n.v.Eval(getWindowList, &res); err != nil {
		return nil, err
	}

	return res, nil
}

func (n *Nvimutil) SplitOpenBuffer(buf nvim.Buffer) error {
	var bwin nvim.Window
	var win nvim.Window

	b := n.v.NewBatch()
	b.CurrentWindow(&bwin)
	b.Command(fmt.Sprintf("sb %d", buf))
	b.CurrentWindow(&win)
	b.SetWindowHeight(win, 15)
	if err := b.Execute(); err != nil {
		return err
	}

	return n.v.SetCurrentWindow(bwin)
}

func (n *Nvimutil) NewScratchBuffer(bufferName string) (*nvim.Buffer, error) {
	var scratchBuf nvim.Buffer
	var bwin nvim.Window
	var win nvim.Window

	b := n.v.NewBatch()
	b.CurrentWindow(&bwin)
	b.Command("silent! execute 'new' '" + bufferName + "'")
	b.CurrentBuffer(&scratchBuf)
	b.SetBufferOption(scratchBuf, "buftype", "nofile")
	b.SetBufferOption(scratchBuf, "bufhidden", "hide")
	b.Command("setlocal noswapfile")
	b.Command("setlocal nobuflisted")
	b.SetBufferOption(scratchBuf, "undolevels", -1)
	b.CurrentWindow(&win)
	b.SetWindowHeight(win, 15)

	if err := b.Execute(); err != nil {
		return nil, err
	}

	if err := n.v.SetCurrentWindow(bwin); err != nil {
		return nil, err
	}

	return &scratchBuf, nil
}

func (n *Nvimutil) Input(ask string) (string, error) {
	var input string
	if err := n.v.Eval(`input("`+ask+`: ")`, &input); err != nil {
		return "", err
	}

	return input, nil
}
