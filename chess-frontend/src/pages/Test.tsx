function Test() {
    const testWebSocket = () => {
        let ws = new WebSocket("ws://localhost:8080/test/ws")
        // 注册连接成功事件
        ws.onopen = () => {
            console.log("连接建立成功...")
            // 发送消息
            ws.send("我是客户端...")
        }
        // 注册接收消息事件
        ws.onmessage = (e) => {
            console.log("接收到消息:" + e.data)
        }
        // 注册错误事件
        ws.onerror = (e) => {
            console.log("连接发生错误")
        }
        // 注册连接断开事件
        ws.onclose = () => {
            console.log("连接断开成功...")
        }
    }
    return (
        <>
            <button onClick={testWebSocket}>
                测试WebSocket
            </button>
        </>
    )
}

export default Test