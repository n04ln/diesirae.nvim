package aoj

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/NoahOrberg/aoj.nvim/config"
	"github.com/NoahOrberg/aoj.nvim/util"
	"github.com/h2non/gentleman"
	"github.com/h2non/gentleman/plugins/body"
)

func Session(id, rawPassword string) (string, error) {
	conf := config.GetConfig()

	cli := gentleman.New()
	cli.URL(conf.Endpoint)

	req := cli.Request()
	req.Method("POST")
	req.Path("/session")
	data := map[string]string{
		"id":       id,
		"password": rawPassword,
	}
	req.Use(body.JSON(data))

	res, err := req.Send()
	if err != nil {
		return "", err
	}
	if !res.Ok {
		return "", util.ErrResponseIsNotOK
	}

	if len(res.Cookies) != 1 {
		return "", util.ErrCookieIsNotFound
	}

	return res.Cookies[0].Raw, nil
}

type SubmitRes struct {
	Token string `json:"token"`
}

func Submit(cookie, problemId, language, sourceCode string) (string, error) {
	conf := config.GetConfig()

	cli := gentleman.New()
	cli.URL(conf.Endpoint)

	req := cli.Request()
	req.Method("POST")
	req.Path("/submissions")
	data := map[string]string{
		"problemId":  problemId,
		"language":   language,
		"sourceCode": sourceCode,
	}
	req.Use(body.JSON(data))
	req.SetHeader("Cookie", cookie)

	res, err := req.Send()
	if err != nil {
		return "", err
	}
	if !res.Ok {
		return "", util.ErrResponseIsNotOK
	}

	var r SubmitRes
	if err := json.Unmarshal([]byte(res.String()), &r); err != nil {
		return "", err
	}

	return r.Token, nil
}

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
}

func Status(cookie, token string) (*SubmissionStatus, error) {
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

	return &substatus, nil
}

func getStatusByJudgeId(cookie string, judgeId int64) (*gentleman.Response, error) {
	conf := config.GetConfig()
	cli := gentleman.New()
	cli.URL(conf.Endpoint)

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
	cli.URL(conf.Endpoint)

	req := cli.Request()
	req.Method("GET")
	req.Path("/submission_records/recent")
	req.SetHeader("Cookie", cookie)

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
