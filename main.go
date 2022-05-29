package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"time"

	"golang.org/x/sync/errgroup"
)

//*启动 HTTP server
func StartHttpServer(srv *http.Server) error {
	http.HandleFunc("/sissi", Sissi)
	http.HandleFunc("/zhudi", Zhudi)
	http.HandleFunc("/cancel", Cancel)
	http.HandleFunc("/timeout", TimeOut)
	fmt.Println("http server start")
	err := srv.ListenAndServe()
	return err
}

// 增加一个 HTTP hanlder
func Sissi(w http.ResponseWriter, req *http.Request) {
	fmt.Println("正在访问localhost:9091/sissi接口")
	io.WriteString(w, "hello, sissi!\n")
}

// 增加一个 HTTP hanlder
func Zhudi(w http.ResponseWriter, req *http.Request) {
	fmt.Println("正在访问localhost:9091/zhudi接口")
	io.WriteString(w, "hello,zhudi!\n")
}

func Cancel(w http.ResponseWriter, req *http.Request){
	fmt.Println("正在访问localhost:9091/cancel接口")
	io.WriteString(w, "Canel!\n")
}
func TimeOut(w http.ResponseWriter, req *http.Request){
	fmt.Println("正在访问localhost:9091/timeout接口")
	time.Sleep(time.Second*80)
	io.WriteString(w, "cannot write!\n")
}

//1. 基于 errgroup 实现一个 http server 的启动和关闭 ，
//以及 linux signal 信号的注册和处理，要保证能够一个退出，全部注销退出。
func main() {
	ctx := context.Background()
	// 定义 withCancel -> cancel() 方法 去取消下游的 Context
	ctx, cancel := context.WithTimeout(ctx,time.Second*100)
	//* 使用 errgroup 进行 goroutine 取消
	group, errCtx := errgroup.WithContext(ctx)
	//http server
	srv := &http.Server{Addr: ":9091"}

	group.Go(func() error {
		return StartHttpServer(srv)
	})

	group.Go(func() error {
		//Done() 返回一个 channel，可以表示 context 被取消的信号：
		//当这个 channel 被关闭时，说明 context 被取消了。
		//注意，这是一个只读的channel。 
		//读一个关闭的 channel 会读出相应类型的零值。并且源码里没有地方会向这个 channel 里面塞入值。换句话说，这是一个 receive-only 的 channel。
		//因此在子协程里读这个 channel，除非被关闭，否则读不出来任何东西。也正是利用了这一点，子协程从 channel 里读出了值（零值）后，就可以做一些收尾工作，尽快退出。
		<-errCtx.Done() //阻塞。因为 cancel、timeout、deadline 都可能导致 Done 被 close
		fmt.Println("http server stop")
		//* 关闭 http server
		return srv.Shutdown(errCtx) 
	})

	channel := make(chan os.Signal, 1) //这里要用 buffer 为1的 chan
	signal.Notify(channel)

	group.Go(func() error {
		for {
			select {
			case <-errCtx.Done(): // 因为 cancel、timeout、deadline 都可能导致 Done 被 close
				fmt.Println("Done")
				//cancel:context canceled
				//context deadline exceeded
				return errCtx.Err()
			case <-channel: // 因为 kill -9 或其他而终止
				cancel()
			}
		}
        return nil
	})

	if err := group.Wait(); err != nil {
		fmt.Println("group error: ", err)
	}
	fmt.Println("all group done!")

}
