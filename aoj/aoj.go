package aoj

import (
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

	return res.String(), nil // TODO: return token only
}
