package test

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net"
	"net/http"
	"testing"
)

// TestWebSocket 测试WebSocket收发消息
func TestWebSocket(t *testing.T) {
	server := gin.Default()
	server.GET("/test/ws", func(ctx *gin.Context) {
		// 获取ws连接
		conn, err := (&websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		}).Upgrade(ctx.Writer, ctx.Request, nil)
		if err != nil {
			fmt.Println("websocket err:", err)
			return
		}
		// 开启协程处理连接
		go func() {
			defer conn.Close()
			for {
				// 获取客户端消息
				messageType, message, err2 := conn.ReadMessage()
				if err2 != nil {
					// 判断是不是超时
					if netErr, ok := err2.(net.Error); ok {
						if netErr.Timeout() {
							fmt.Printf("ReadMessage timeout remote:%v\n", conn.RemoteAddr())
						}
					}
					// 其他错误
					fmt.Println("ReadMessage error:", err2)
					return
				}
				// 打印消息
				fmt.Println("message type:", messageType)
				fmt.Println("message:", string(message))
				// 写入数据
				var msg = []byte("我是server")
				err = conn.WriteMessage(1, msg)
				if err != nil {
					fmt.Println("write message error:", err)
				}
			}
		}()
	})
	server.Run(":8080")
}
