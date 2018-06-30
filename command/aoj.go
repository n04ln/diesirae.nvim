package command

import (
	"fmt"
	"time"

	"github.com/NoahOrberg/diesirae.nvim/aoj"
	"github.com/NoahOrberg/diesirae.nvim/config"
	"github.com/NoahOrberg/nimvle.nvim/nimvle"
	"github.com/neovim/go-client/nvim"
)

var (
	done chan string
	ss   = []string{
		"\\", "|", "/", "-", "\\", "|", "/", "-",
	}
)

func nimvleNew(v *nvim.Nvim) *nimvle.Nimvle {
	return nimvle.New(v, "diesirae.nvim")
}

type AOJ struct {
	Cookie             string                                    // NOTE: ログインしたら、ここにクッキーをいれる
	Config             config.AOJConfig                          // NOTE: 環境変数から取得した情報格納
	SubmittedStatuses  map[nvim.Buffer]([]*aoj.SubmissionStatus) // NOTE: 提出したとき、あとからそれを確認できる用にするため、キーをバッファ番号にして確認用Tokenを保存する
	ScratchBuffer      *nvim.Buffer                              // NOTE: Statusなどを吐く
	DebugScratchBuffer *nvim.Buffer                              // NOTE: debug用. panicの情報などを吐く
	IsValidCookie      bool
}

func NewAOJ() (*AOJ, error) {
	conf := config.GetConfig()

	cookie, err := aoj.Session(conf.ID, conf.RawPassword)
	if err != nil {
		return &AOJ{
			Config:            conf,
			SubmittedStatuses: map[nvim.Buffer]([]*aoj.SubmissionStatus){},
		}, err
	}

	return &AOJ{
		Cookie:            cookie,
		Config:            conf,
		SubmittedStatuses: map[nvim.Buffer]([]*aoj.SubmissionStatus){},
		IsValidCookie:     true,
	}, nil
}

func (a *AOJ) SetStatusByBuffer(buf nvim.Buffer, stat *aoj.SubmissionStatus) {
	a.SubmittedStatuses[buf] = append(a.SubmittedStatuses[buf], stat)
}

func (a *AOJ) GetRecentStatusByBuffer(buf nvim.Buffer) (*aoj.SubmissionStatus, bool) {
	stats, ok := a.SubmittedStatuses[buf]
	if !ok {
		return nil, false
	}

	return stats[len(stats)-1], true
}

func (a *AOJ) panicLog(v *nvim.Nvim) {
	// only when debug mode
	if config.GetConfig().Mode == "debug" {
		n := nimvleNew(v)

		err := recover()

		if a.DebugScratchBuffer == nil && err != nil {
			a.DebugScratchBuffer, _ = n.NewScratchBuffer("DEBUG")
		}

		if err != nil {
			n.SetContentToBuffer(*a.DebugScratchBuffer, fmt.Sprintf("%v", err))
		}
	}
}

type loading uint64

func (s *loading) String() string {
	tmp := loading(uint64(*s) + 1)
	*s = tmp
	return ss[uint64(*s)%uint64(len(ss))]
}

func drawLoadingCycle(nimvle *nimvle.Nimvle, scratchBuf *nvim.Buffer) {
	l := new(loading)
	for {
		select {
		case s := <-done:
			if s == "" {
				return
			}
			err := nimvle.SetContentToBuffer(*scratchBuf, s)
			if err != nil {
				nimvle.Log(err.Error())
				return
			}
			return
		default:
			err := nimvle.SetStringerContentToBuffer(*scratchBuf, l)
			if err != nil {
				nimvle.Log(err.Error())
				return
			}
			time.Sleep(100 * time.Millisecond)
		}
	}
}

type empty struct{}

func (e *empty) String() string {
	return ""
}

func flushLoadingCycle(nimvle *nimvle.Nimvle, scratch *nvim.Buffer, err error) {
	if err == nil {
		return
	}

	done <- ""
	e := &empty{}
	nimvle.ShowScratchBuffer(*scratch, e)
}
