package router

import (
	"context"

	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/os/gctx"
)

// Setup 注册各项目自己的 api 和 全局前置中间件
// myApi 具体项目的路由配置
// myMiddlewares 具体项目所需的全局前置中间件（最早经过的中间件）
func Setup(myApi map[string]GroupConfig, myMiddlewares []func(*ghttp.Request)) {
	s := g.Server()
	s.Group("/", func(group *ghttp.RouterGroup) {
		myGroup := &myGroup{
			group: group,
		}
		for api, cfg := range myApi {
			for _, apiCfg := range cfg.ApiCfg {
				apiCfgClone := apiCfg
				myGroup.Clone().
					ALL("/"+api+"/"+apiCfgClone.Url, apiCfgClone.ApiFunc, apiCfgClone.ParamModel)
			}
		}
	})
	ctx := gctx.New()

	go sendRouterToGate(ctx, myApi)
}
func sendRouterToGate(ctx context.Context, apis map[string]GroupConfig) {

}