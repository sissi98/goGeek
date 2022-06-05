package router

import (
	"goGeek/api"
	"goGeek/internal/pkg/router"

	"github.com/gogf/gf/net/ghttp"
)

func init() {
	router.Setup(myApi, []func(*ghttp.Request){})
}

var myApi = map[string]router.GroupConfig{
	"web-api": {
		Middlewares: []func(*ghttp.Request){},
		ApiCfg: []router.ApiConfig{
			{
				Url:        "sissi/v1/hello-world",
				ApiFunc:    api.hello_world.HelloWorld,
				ParamModel: nil,
			},
			
		},
	},
}
