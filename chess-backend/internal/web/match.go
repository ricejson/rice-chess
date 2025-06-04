package web

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net"
	"net/http"
)

type MatchHandler struct {
}

func NewMatchHandler() *MatchHandler {
	return &MatchHandler{}
}

func (mh *MatchHandler) FindMatch(ctx *gin.Context) {
	// 建立 WebSocket 连接
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
	// 异步处理该连接
	go func() {
		defer conn.Close()
		for {
			messageType, message, err2 := conn.ReadMessage()
			if err2 != nil {
				// 判断是不是超时错误
				if netErr, ok := err2.(net.Error); ok {
					if netErr.Timeout() {
						fmt.Printf("ReadMessage timeout remote:%v\n", conn.RemoteAddr())
					}
				}
				// 是其他错误
				fmt.Println("ReadMessage error:", err2)
				return
			}
			// 打印消息
			fmt.Println("message type:", messageType)
			fmt.Println("message:", string(message))
			type MatchReq struct {
				Message string `json:"message"`
			}
			// 序列化消息
			var req MatchReq
			if json.Unmarshal(message, &req) != nil {
				// 反序列化错误
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"code":    500,
					"message": "系统错误！",
				})
				return
			}
			// 根据message字段返回响应
			type MatchResp struct {
				Code    int    `json:"code"`
				Message string `json:"message"`
			}
			if req.Message == "startMatch" {
				// 开始匹配
				var resp = MatchResp{
					Code:    200,
					Message: "startMatch",
				}
				bytes, _ := json.Marshal(resp)
				conn.WriteMessage(1, bytes)
			} else if req.Message == "stopMatch" {
				// 停止匹配
				var resp = MatchResp{
					Code:    200,
					Message: "stopMatch",
				}
				bytes, _ := json.Marshal(resp)
				conn.WriteMessage(1, bytes)
			}
		}
	}()
}

func (mh *MatchHandler) RegisterRoutes(server *gin.Engine) {
	mg := server.Group("/match")
	mg.GET("/findMatch", mh.FindMatch)
}
