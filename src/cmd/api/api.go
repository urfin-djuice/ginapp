package main

import (
	"oko/pkg/account"
	"oko/pkg/action"
	"oko/pkg/domain"
	"oko/pkg/ginapp"
	"oko/pkg/ginapp/controller"
	"oko/pkg/links"
	"oko/pkg/proxy"
	"oko/pkg/repost"
	"oko/pkg/rss"
	"oko/pkg/rule"
	"oko/pkg/trigger"
	"oko/pkg/valid"
)

// @title OKO API
// @version 1.0
// @BasePath /api

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

func main() {
	var api = ginapp.App{
		RootRote:     "/api/",
		RootHandlers: controller.HandlerList{},
		Ctrls: []controller.Ctrl{
			account.NewController(),
			domain.NewController(),
			links.NewController(),
			repost.NewController(),
			proxy.NewController(),
			action.NewController(),
			rule.NewController(),
			trigger.NewController(),
			rss.NewController(),
		},
		Validators: valid.Validators,
	}
	api.Init()
	api.Do()
}
