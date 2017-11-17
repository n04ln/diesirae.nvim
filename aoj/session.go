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

func IsAliveSession(cookie string) bool {
	conf := config.GetConfig()

	cli := gentleman.New()
	cli.URL(conf.Endpoint)

	req := cli.Request()
	req.Method("GET")
	req.Path("/self")
	req.SetHeader("Cookie", cookie)

	res, err := req.Send()
	if err != nil {
		return false
	}
	if !res.Ok {
		return false
	}

	return true
}
