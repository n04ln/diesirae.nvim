package aoj

import (
	"github.com/n04ln/diesirae.nvim/config"
	"github.com/n04ln/diesirae.nvim/util"
	"gopkg.in/h2non/gentleman.v2"
	"gopkg.in/h2non/gentleman.v2/plugins/body"
)

func Session(id, rawPassword string) (string, error) {
	conf := config.GetConfig()

	cli := gentleman.New()
	cli.URL(conf.API)

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
	cli.URL(conf.API)

	req := cli.Request()
	req.Method("GET")
	req.Path("/self")
	req.SetHeader("Cookie", cookie)
	req.SetHeader("Content-Type", "application/json")

	res, err := req.Send()
	if err != nil {
		return false
	}
	if !res.Ok {
		return false
	}

	return true
}
