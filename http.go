package main

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"time"
)

// getRequest 处理 HTTP 请求的回调函数
func (rf *Raft) getRequest(writer http.ResponseWriter, request *http.Request) {
	// 解析请求的表单数据
	request.ParseForm()

	// 检查请求的 URL 是否包含 "message" 参数，并且当前有领导者节点
	// 例如，http://localhost:8080/req?message=ohmygod
	if len(request.Form["message"]) > 0 && rf.currentLeader != "null" && len(request.Form["count"]) > 0 {
		message := request.Form["message"][0] // 提取消息内容
		var count int
		count, _ = strconv.Atoi(request.Form["count"][0])

		for i := 0; i < count; i++ {
			m := LogEntry{
				// Command: fmt.Sprintf("%s", message),
				Command: message,
			}
			go rf.Boradcast(m)
			// 接收到消息后，直接将其转发给领导者节点进行处理

		}
		fmt.Println("http监听到了消息")
		writer.Write([]byte("收到")) // 向客户端返回确认消息
	}

	//算术操作，例如：http://localhost:8080/req?op=add&a=1
	// 	if len(request.Form["op"]) > 0 && len(request.Form["a"]) > 0 && len(request.Form["b"]) > 0 {
	// 		op := request.Form["op"][0]
	// 		a, _ := strconv.Atoi(request.Form["a"][0])
	// 	switch op {
	// 		case "+":
	// 			rf.lacalnum += a
	// 			fmt.Println("加法", a, rf.lacalnum)
	// 		case "-":
	// 			rf.lacalnum -= a
	// 			fmt.Println("减法", a, rf.lacalnum)
	// 		case "*":
	// 			rf.lacalnum *= a
	// 			fmt.Println("乘法", a, rf.lacalnum)
	// 		case "/":
	// 			if a != 0 {
	// 				rf.lacalnum /= a
	// 				fmt.Println("除法", a, rf.lacalnum)
	// 			} else {
	// 				fmt.Println("错误：除数不能为零")
	// 			}
	// 		default:
	// 			fmt.Printf("不支持的操作符：%s\n", op)
	// }
	// 		writer.Write([]byte("操作已执行")) // 向客户端返回确认消息

	// 	}
}

var HandleFuncFlag = true

func (rf *Raft) httpListen() {
	// 创建一个 http.Server 实例
	rf.server = &http.Server{
		Addr:    ":8080",
		Handler: nil, // 使用默认的 http.DefaultServeMux
	}

	if HandleFuncFlag {
		HandleFuncFlag = false
		http.HandleFunc("/req", rf.getRequest)
	}
	go func() {
		fmt.Println("监听8080")
		if err := rf.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Println("Server error:", err)
		}
	}()
}

func (rf *Raft) stopListening() {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	if err := rf.server.Shutdown(ctx); err != nil {
		fmt.Println("Shutdown error:", err)
	} else {
		fmt.Println("监听已停止")
	}
}
