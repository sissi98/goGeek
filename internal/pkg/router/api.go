package router

import (
	"context"

	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/os/glog"
)
type MyReq struct {
	*ghttp.Request
	logger  *glog.Logger    // 设置了请求上下文的 logger
	i18nCtx context.Context // 国际化组件用到的上下文变量

	Params interface{} // 请求的参数

	// 携带在请求中的与请求相关的参数
	GameId         int      // 游戏ID，由于大部分接口都有该参数，故抽离出来
	AttributeLangs []string // 所选择的语言
	Lang           string   // 语言，用于接口报错信息多语言的处理
	Tid            int      // 接口追踪ID，取自请求参数或自动根据时间戳生成
}

type MyApiFunc func(*MyReq)

type ApiConfig struct {
	Url              string // 在同一个api下不允许重复
	ApiFunc          MyApiFunc
	ParamModel       interface{}
	Middlewares      []func(*ghttp.Request) // 额外的中间件
	SkipMiddlewares  []func(*ghttp.Request) // 跳过的中间件
	ExtraMiddlewares []string               // 门户支持的额外的中间件
	// method string 想了下，算了，没啥用，用all即可
}

type GroupConfig struct {
	Middlewares []func(*ghttp.Request)
	ApiCfg      []ApiConfig
}
