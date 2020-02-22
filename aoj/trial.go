package aoj

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/n04ln/diesirae.nvim/config"
	"github.com/n04ln/diesirae.nvim/util"
	"github.com/h2non/gentleman"
)

var ErrCompileError = errors.New("Compile Error. check your source code.")

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

	temp := "Serial %d:\nInput:\n%s\nExpected Output:\n%s\nActual Output:\n%s\n*****"
	var res string
	for _, ss := range s.Samples {
		res += fmt.Sprintf(temp, ss.Serial, ss.Input, ss.Output, ss.Actual)
	}

	for i, ss := range s.Samples {
		if ss.Actual != ss.Output {
			res = "Wrong Answer...\n Don't worry! this is testing sample I/O!\n*****" + res
			break
		}

		if i+1 == len(s.Samples) {
			res = "All cases AC!\n Good Job!\n*****" + res
		}
	}

	return res
}

func replaceBuildCommands(bc []string, bin, source string) []string {
	res := make([]string, 0, len(bc))
	for _, b := range bc {
		res = append(res, b)
	}

	targets := []string{
		"*bin*",
		"*source*",
	}

	replacements := []string{
		bin,
		source,
	}

	for i := 0; i < len(bc); i++ {
		for j := 0; j < len(targets); j++ {
			if targets[j] == res[i] {
				res[i] = replacements[j]
			}
		}
	}

	return res
}

func (samples *Samples) ExecSamples(ext, sourceCode string,
	buildCommands, execCommands []string, timeLimit int) (*string, error) {

	// tempfile用dir作成
	dir, err := ioutil.TempDir("", "diesirae")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(dir) // clean up

	// tempfile
	tmpfilepath := filepath.Join(dir, "tmpfile."+ext)
	if err := ioutil.WriteFile(tmpfilepath, []byte(sourceCode), 0666); err != nil {
		return nil, err
	}

	// build tempfile
	binpath := filepath.Join(dir, "tmp")
	buildCommands = replaceBuildCommands(buildCommands, binpath, tmpfilepath)
	if len(buildCommands) < 2 {
		if !(buildCommands == nil || len(buildCommands) == 0) {
			return nil, errors.New("invalid commands")
		}
	}
	if len(buildCommands) >= 2 {
		_, err = exec.Command(buildCommands[0], buildCommands[1:]...).Output()
		if err != nil {
			errStr := err.Error()
			return &errStr, ErrCompileError
		}
	}

	for i, sample := range samples.Samples {
		// run tempfile within Stdin
		var cmd *exec.Cmd
		execCommands = replaceBuildCommands(execCommands, binpath, tmpfilepath)
		if len(execCommands) < 2 {
			if len(execCommands) == 0 {
				return nil, errors.New("invalid commands")
			}
			cmd = exec.Command(execCommands[0])
		} else {
			cmd = exec.Command(execCommands[0], execCommands[1:]...)
		}
		if cmd == nil {
			return nil, errors.New("invalid commands")
		}
		stdin, err := cmd.StdinPipe()
		if err != nil {
			return nil, err
		}
		io.WriteString(stdin, sample.Input)
		stdin.Close()

		cout := make(chan string)
		ready, latch := make(chan struct{}), make(chan struct{})
		go func() {
			ready <- struct{}{}
			<-latch

			// 片方しか使わない
			bout, err := cmd.CombinedOutput()

			if err != nil {
				cout <- err.Error()
			} else {
				cout <- string(bout)
			}
		}()

		<-ready
		close(latch) // start

		// timeLimitまで待機
		time.Sleep(time.Duration(timeLimit) * time.Second)

		select {
		case out := <-cout:
			samples.Samples[i].Actual = string(out)
		default:
			samples.Samples[i].Actual = "no output\n" // Nothing output
			// 後処理
			if err := cmd.Process.Kill(); err != nil {
				return nil, err
			}
		}

	}

	return nil, nil
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
