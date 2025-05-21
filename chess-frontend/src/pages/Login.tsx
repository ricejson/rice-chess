import LoginPage, { Username, Password, TitleSignup, TitleLogin, Submit, Title, Logo } from '@react-login-page/page8';
import LoginLogo from 'react-login-page/logo';
import {useState} from "react";
import axios from "axios";
import {message} from "antd";

const styles = { height: 690 };
const Login: React.FC = () => {
    const [messageApi, contextHolder] = message.useMessage();
    // 注册数据
    const [signUpData, setSignUpData] = useState({
        username: '',
        password: '',
    })
    // 登录数据
    const [loginData, setLoginData] = useState({
        username: '',
        password: '',
    })
    // 处理注册输入框变化事件
    const handleSignUpInputChange = (key: string, value: string) => {
        setSignUpData((prev) => ({...prev, [key]: value}))
    }
    // 处理登录输入框变化事件
    const handleLoginInputChange = (key: string, value: string) => {
        setLoginData((prev) => ({...prev, [key]: value}))
    }
    // 处理注册事件
    const handleRegister = () => {
        // 发送请求到后端
        axios.post(
            'http://localhost:8081/user/register',
            signUpData,
        ).then((resp) => {
            console.log(resp)
            if (resp?.data && resp.data.code === 200) {
                // 注册成功
                messageApi.success('注册成功！')
            } else {
                // 注册失败
                messageApi.error(resp?.data?.msg ?? '注册失败！')
            }
        }).catch((error) => {
            console.log(error)
        })
    }
    // 处理登录事件
    const handleLogin = () => {
        // 发送请求到后端
        axios.post(
            'http://localhost:8081/user/login',
            loginData,
        ).then((resp) => {
            console.log(resp)
            if (resp.data && resp.data.code === 200) {
                // 登录成功
                messageApi.success('登录成功！')
                location.assign('/')
            } else {
                // 登录失败
                messageApi.error('登录失败！')
            }
        }).catch((error) => {
            console.log(error)
        })
    }
    return (
        <div style={styles}>
            {contextHolder}
            <LoginPage>
                <Title />
                <TitleSignup>注册</TitleSignup>
                <TitleLogin>登录</TitleLogin>
                <Logo>
                    <LoginLogo />
                </Logo>
                <Username onChange={(e) => {handleLoginInputChange("username", e.target.value)}} label="用户名" placeholder="请输入用户名" name="loginUsername" />
                <Password onChange={(e) => {handleLoginInputChange("password", e.target.value)}} label="密码" placeholder="请输入密码" name="loginUserPassword" />
                <Submit keyname="submit" onClick={handleLogin}>提交</Submit>
                <Submit keyname="reset">重置</Submit>

                <Username panel="signup" onChange={(e) => {handleSignUpInputChange("username", e.target.value)}} label="用户名" placeholder="请输入用户名" keyname="e-mail" />
                <Password panel="signup" onChange={(e) => {handleSignUpInputChange("password", e.target.value)}} label="密码" placeholder="请输入密码" keyname="password" />
                <Password panel="signup" visible={false} keyname="confirm-password" />
                <Submit panel="signup" keyname="signup-submit" onClick={handleRegister}>
                    注册
                </Submit>
                <Submit panel="signup" keyname="signup-reset">
                    重置
                </Submit>
            </LoginPage>
        </div>
    )
}

export default Login;