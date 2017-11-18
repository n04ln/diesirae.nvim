package command

import (
	"fmt"
	"strings"

	"github.com/NoahOrberg/aoj.nvim/aoj"
	"github.com/NoahOrberg/aoj.nvim/config"
	"github.com/neovim/go-client/nvim"
)

type AOJ struct {
	Cookie            string                                    // NOTE: ログインしたら、ここにクッキーをいれる
	Config            config.AOJConfig                          // NOTE: 環境変数から取得した情報格納
	SubmittedStatuses map[nvim.Buffer]([]*aoj.SubmissionStatus) // NOTE: 提出したとき、あとからそれを確認できる用にするため、キーをバッファ番号にして確認用Tokenを保存する
	ScratchBuffer     *nvim.Buffer
}

func (a *AOJ) StatusMessage(stat *aoj.SubmissionStatus) string {
	caseVerdictTemp := `testcase: %v, Memory: %vkB, CpuTime: %vs, Status: %v`
	messages := make([]string, 0, len(stat.CaseVerdicts)+3)
	for _, cv := range stat.CaseVerdicts {
		messages = append(messages, fmt.Sprintf(caseVerdictTemp, cv.Label, cv.Memory, cv.CpuTime, cv.Status))
	}

	if stat.CompileError != "" {
		messages = append(messages, fmt.Sprintf("CompileError: %v", stat.CompileError))
	}

	if stat.RuntimeError != "" {
		messages = append(messages, fmt.Sprintf("RuntimeError: %v", stat.RuntimeError))
	}

	if stat.UserOutput != "" {
		messages = append(messages, fmt.Sprintf("UserOutput: %v", stat.UserOutput))
	}

	return strings.Join(messages, "\n")
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

func NewAOJ() (*AOJ, error) {
	conf := config.GetConfig()

	cookie, err := aoj.Session(conf.ID, conf.RawPassword)
	if err != nil {
		return nil, err
	}

	return &AOJ{
		Cookie:            cookie,
		Config:            conf,
		SubmittedStatuses: map[nvim.Buffer]([]*aoj.SubmissionStatus){},
	}, nil
}
