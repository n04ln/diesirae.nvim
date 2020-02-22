package aoj

import (
	"encoding/json"

	"github.com/n04ln/diesirae.nvim/config"
	"github.com/n04ln/diesirae.nvim/util"
	"gopkg.in/h2non/gentleman.v2"
	"gopkg.in/h2non/gentleman.v2/plugins/body"
)

type SubmitRes struct {
	Token string `json:"token"`
}

func Submit(cookie, problemId, language, sourceCode string) (string, error) {
	conf := config.GetConfig()

	cli := gentleman.New()
	cli.URL(conf.API)

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
