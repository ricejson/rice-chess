import axios from "axios";

const request = axios.create({
    baseURL: "http://localhost:8081",
    withCredentials: true,
})

// 配置响应拦截器
request.interceptors.response.use((resp) => {
    // 从header中获取token
    const token = resp.headers["x-jwt-token"];
    if (token) {
        // 本地存储
        localStorage.setItem("token", token);
    }
    // 如果接收到401请求直接跳转到登录界面
    if (resp.status === 401) {
        window.location.href = '/user/login';
    }
    return resp;
}, (err) => {
    console.log(err)
    if (err.response.status === 401) {
        window.location.href= "/user/login";
    }
    return err
});

// 配置请求拦截器
request.interceptors.request.use((req) => {
    const token = localStorage.getItem("token");
    req.headers.setAuthorization("Bearer " + token);
    return req;
}, (err) => {
    console.log(err)
});

export default request;