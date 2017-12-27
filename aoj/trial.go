package aoj

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"

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

func (samples *Samples) ExecSamples(fileType, sourceCode string) error {
	// TODO: とりあえずGoだけ
	switch fileType {
	case "Go":
		fp, err := ioutil.TempFile("", "diesirae")
		if err != nil {
			return err
		}
		defer os.Remove(fp.Name())
		defer fp.Close()
		defer os.Remove(fp.Name() + ".go")

		if err := os.Rename(fp.Name(), fp.Name()+".go"); err != nil {
			return err
		}

		if _, err := fp.Write([]byte(sourceCode)); err != nil {
			return err
		}

		for i, sample := range samples.Samples {
			cmd := exec.Command("go", "run", fp.Name()+".go")
			stdin, err := cmd.StdinPipe()
			if err != nil {
				return err
			}
			_, err = io.WriteString(stdin, sample.Input)
			if err != nil {
				return err
			}
			err = stdin.Close()
			if err != nil {
				return err
			}
			out, err := cmd.Output()
			if err != nil {
				return err
			}
			samples.Samples[i].Actual = string(out)
		}
	default:
		return errors.New("only support Golang :)")
	}

	return nil
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
