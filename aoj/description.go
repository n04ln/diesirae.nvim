package aoj

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/n04ln/diesirae.nvim/config"
	"github.com/n04ln/diesirae.nvim/util"
	"gopkg.in/h2non/gentleman.v2"
)

type Commentary struct {
	Filter  []string `json:"filter"`
	Pattern string   `json:"pattern"`
	Type    string   `json:"type"`
}

type Description struct {
	// TODO: Unknown type
	// "bookmarks": null,
	// "source": null,
	// "recommend": null,
	Commentaries    []*Commentary `json:"commentaries"`
	CreatedAt       int64         `json:"created_at"`
	HTML            string        `json:"html"`
	IsSolved        bool          `json:"isSolved"`
	Language        string        `json:"language"`
	MemoryLimit     int           `json:"memory_limit"`
	ProblemId       string        `json:"problem_id"`
	Recommendations int           `json:"recommendations"`
	Score           float64       `json:"score"`
	ServerTime      int           `json:"server_time"`
	SolvedUser      int           `json:"solvedUser"`
	SuccessRate     float64       `json:"successRate"`
	TimeLimit       int           `json:"time_limit"` // sec
}

func (d *Description) String() string {
	tmpl := "ProblemId: %v\nIsSolved?: %v\nCreatedAt: %v\nTimeLimit: %v sec\n\nplease see below:\n http://judge.u-aizu.ac.jp/onlinejudge/description.jsp?id=%v"

	var problemId string
	problemId = d.ProblemId

	var isSolvedMsg string
	if d.IsSolved {
		isSolvedMsg = "yes"
	} else {
		isSolvedMsg = "no"
	}

	t := time.Unix(d.CreatedAt/1000, 0)
	var createdAt string
	createdAt = t.Format("2006-01-02")

	var timeLimit int
	timeLimit = d.TimeLimit

	return fmt.Sprintf(tmpl, problemId, isSolvedMsg, createdAt, timeLimit, problemId)
}

func GetDescriptionByProblemId(cookie, id string) (*Description, error) {
	conf := config.GetConfig()

	cli := gentleman.New()
	cli.URL(conf.API)

	req := cli.Request()
	req.Method("GET")
	req.Path("/resources/descriptions/ja/" + id)
	req.SetHeader("Cookie", cookie)

	res, err := req.Send()
	if err != nil {
		return nil, err
	}
	if !res.Ok {
		return nil, util.ErrResponseIsNotOK
	}

	var description Description
	if err := json.Unmarshal(res.Bytes(), &description); err != nil {
		return nil, err
	}

	return &description, nil
}
