package aoj

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/n04ln/diesirae.nvim/config"
	"github.com/n04ln/diesirae.nvim/util"
	"gopkg.in/h2non/gentleman.v2"
)

type errStatus struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type caseVerdicts struct {
	// NOTE: I DO NOT know the following types
	// CaseName   XXXXX `json:"caseName"`
	// InputSize  XXXXX `json:"inputSize"`
	// OutputSize XXXXX `json:"outputSize"`
	CpuTime int64  `json:"cpuTime"`
	Memory  int64  `json:"memory"`
	Serial  int    `json:"serial"`
	Label   string `json:"label"`
	Status  string `json:"status"`
}

type SubmissionStatus struct {
	CaseVerdicts []caseVerdicts `json:"caseVerdicts"`
	CompileError string         `json:"compileError"`
	RuntimeError string         `json:"runtimeError"`
	UserOutput   string         `json:"userOutput"`
	ProblemId    string
	Time         time.Time
}

func (stat *SubmissionStatus) String() string {
	isAllCasesAC := true
	caseVerdictTemp := `testcase: %v, Memory: %vkB, CpuTime: %vs, Status: %v`
	messages := make([]string, 0, len(stat.CaseVerdicts)+3)
	for _, cv := range stat.CaseVerdicts {
		messages = append(messages, fmt.Sprintf(caseVerdictTemp, cv.Label, cv.Memory, cv.CpuTime, cv.Status))
		if cv.Status != "AC" {
			isAllCasesAC = false
		}
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

	var comment string
	if isAllCasesAC {
		comment = "\n✅ CONGRATULATION!"
	} else {
		comment = "\n❌ NOT BE ACCEPTED!"
	}

	return "Submission result:\n" + strings.Join(messages, "\n") + comment
}

func (s *SubmissionStatus) CheckAC() (string, error) {
	for _, caseVerdict := range s.CaseVerdicts {
		if caseVerdict.Status != "AC" {
			return fmt.Sprintf("testcase %s: %s", caseVerdict.Label, caseVerdict.Status), nil
		}
	}
	return fmt.Sprintf("All case: AC"), nil
}

func Status(cookie, token, problemId string) (*SubmissionStatus, error) {
	recents, err := getRecentSubmissionRecords(cookie)
	if err != nil {
		return nil, err
	}

	var submission recentSubmission
	for _, r := range recents {
		if r.Token == token {
			submission = r
		}
	}

	var cnt int
RETRY:
	res, err := getStatusByJudgeId(cookie, submission.JudgeId)
	if err != nil {
		return nil, err
	}
	if !res.Ok {
		if cnt <= 120 {
			time.Sleep(3 * time.Second)
			cnt++
			goto RETRY
		}
		return nil, util.ErrResponseIsNotOK
	}

	var substatus SubmissionStatus
	if err := json.Unmarshal([]byte(res.String()), &substatus); err != nil {
		return nil, err
	}
	substatus.Time = time.Now()
	substatus.ProblemId = problemId

	return &substatus, nil
}

func getStatusByJudgeId(cookie string, judgeId int64) (*gentleman.Response, error) {
	conf := config.GetConfig()
	cli := gentleman.New()
	cli.URL(conf.API)

	req := cli.Request()
	req.Method("GET")
	path := "/verdicts/" + fmt.Sprintf("%d", judgeId)
	req.Path(path)
	req.SetHeader("Cookie", cookie)

	return req.Send()
}

type recentSubmission struct {
	Accuracy       string `json:"accuracy"`
	CodeSize       int    `json:"codeSize"`
	CpuTime        int    `json:"cpuTime"`
	JudgeDate      int64  `json:"judgeDate"`
	JudgeId        int64  `json:"judgeId"`
	JudgeType      int    `json:"judgeType"`
	Language       string `json:"language"`
	Memory         int    `json:"memory"`
	ProblemId      string `json:"problemId"`
	ProblemTitle   string `json:"problemTitle"`
	Status         int    `json:"status"`
	SubmissionDate int64  `json:"submissionDate"`
	Token          string `json:"token"`
	UserId         string `json:"userId"`
}

func getRecentSubmissionRecords(cookie string) ([]recentSubmission, error) {
	conf := config.GetConfig()

	cli := gentleman.New()
	cli.URL(conf.API)

	req := cli.Request()
	req.Method("GET")
	req.Path("/submission_records/recent")
	req.SetHeader("Cookie", cookie)
	req.SetHeader("Content-Type", "application/json")

	res, err := req.Send()
	if err != nil {
		return nil, err
	}
	if !res.Ok {
		return nil, util.ErrResponseIsNotOK
	}

	var r []recentSubmission
	if err := json.Unmarshal([]byte(res.String()), &r); err != nil {
		return nil, err
	}

	return r, nil
}
