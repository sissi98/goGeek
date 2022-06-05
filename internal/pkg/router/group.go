package router

import (
	"fmt"
	"reflect"

	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/i18n/gi18n"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/os/gtime"
)

type myGroup struct {
	group *ghttp.RouterGroup
}

func (g *myGroup) Middleware(handlers ...func(r *ghttp.Request)) *myGroup {
	g.group.Middleware(handlers...)
	return g
}

func (g *myGroup) Clone() *myGroup {
	return &myGroup{
		group: g.group.Clone(),
	}
}

func (g *myGroup) ALL(pattern string, apiFunc MyApiFunc, paramModel interface{}) {
	g.group.ALL(pattern, func(r *ghttp.Request) {
		CallApi(r, apiFunc, paramModel)
	})
}

func CallApi(r *ghttp.Request, apiFunc MyApiFunc, paramModel interface{}) {
	req := NewReq(r)

	t := reflect.TypeOf(paramModel)
	valueType := t.Elem()            // 得到结构体对象的类型
	newPtr := reflect.New(valueType) // 产生指向此结构体类型的指针
	p := newPtr.Interface()

	// 接口参数校验
	if err := r.Parse(p); err != nil {
		fmt.Println(err)
	}
	req.Params = p

	apiFunc(r)
}
func NewReq(gfr *ghttp.Request) *MyReq {
	req := &MyReq{
		Request:        gfr,
		GameId:         gfr.GetInt("game_id"),
		AttributeLangs: gfr.GetArray("attribute_langs", nil),
		Tid:            gfr.GetInt("trace_id", gtime.TimestampMicro()),
		logger:         g.Log().Ctx(gfr.Context()),
	}

	req.Lang = gfr.GetString("lang", "en") // 系统语言，默认使用英文
	// 确定国际化组件用到的上下文变量
	req.i18nCtx = gi18n.WithLanguage(gfr.Context(), req.Lang)
	return req
}

