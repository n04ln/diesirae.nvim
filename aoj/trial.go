package aoj

import (
	"encoding/json"
	"fmt"

	"github.com/NoahOrberg/diesirae.nvim/config"
	"github.com/NoahOrberg/diesirae.nvim/util"
	"github.com/h2non/gentleman"
)

type SampleInputoutput struct {
	ProblemID string `json:"problemId"`
	Serial    int    `json:"serial"`
	Input     string `json:"in"`
	Output    string `json:"out"`
	Actual    string
}

type Samples struct {
	Samples []*SampleInputoutput
}

func (s *Samples) String() string {
	if len(s.Samples) == 0 {
		return "no test cases"
	}

	temp := "Serial %d:\nInput:\n%s===\nOutput:\n%s===\nActual:\n%s===\n"
	var res string
	for _, ss := range s.Samples {
		res += fmt.Sprintf(temp, ss.Serial, ss.Input, ss.Output, ss.Actual)
	}

	for i, ss := range s.Samples {
		if ss.Actual != ss.Output {
			res = "Wrong Answer...\n===" + res
			break
		}

		if i+1 == len(s.Samples) {
			res = "All cases AC!\n===" + res
		}
	}

	return res
}

func GetSampleInputOutput(problemId string) (*Samples, error) {
	conf := config.GetConfig()

	cli := gentleman.New()
	cli.URL(conf.DataAPI)

	req := cli.Request()
	req.Method("GET")
	req.Path("/testcases/samples/" + problemId)

	res, err := req.Send()
	if err != nil {
		return nil, err
	}
	if !res.Ok {
		return nil, util.ErrResponseIsNotOK
	}

	samples := make([]*SampleInputoutput, 0, 0)
	if err := json.Unmarshal(res.Bytes(), &samples); err != nil {
		return nil, err
	}

	return &Samples{
		Samples: samples,
	}, nil
}
