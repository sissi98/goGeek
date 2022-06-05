package api

import (
	"fmt"

	"github.com/gogf/gf/net/ghttp"
)

func hello(r *ghttp.Request){
	fmt.Println("hello")
}