package api

import (
	"fmt"

	"github.com/gogf/gf/net/ghttp"
)
type HelloWorldApi struct{}

func (*HelloApi)HelloWorld(r *ghttp.Request){
	fmt.Println("hello world")
}