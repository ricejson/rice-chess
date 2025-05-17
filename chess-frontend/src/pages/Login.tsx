import {Button, Card, Flex, Form, Input} from 'antd';
import { LockOutlined, UserOutlined } from '@ant-design/icons';
import axios from "axios";

const cardStyle: React.CSSProperties = {
    width: 620,
};

const imgStyle: React.CSSProperties = {
    display: 'block',
    width: 273,
};

const LoginPage: React.FC = () => {
    const onFinish = (values: any) => {
        // 发送请求到后端
        axios.post(
            'http://localhost:8081/user/login',
            values,
        ).then((resp) => {
            console.log(resp)
        }).catch((error) => {
            console.log(error)
        })
    };
    return (
        <Card hoverable style={cardStyle} styles={{ body: { padding: 0, overflow: 'hidden' } }}>
            <Flex justify="space-between">
                <img
                    alt="avatar"
                    src="https://zos.alipayobjects.com/rmsportal/jkjgkEfvpUPVyRjUImniVslZfWPnJuuZ.png"
                    style={imgStyle}
                />
                <Flex vertical align="flex-end" justify="space-between" style={{ padding: 32 }}>
                    <Form
                        name="login"
                        initialValues={{ remember: true }}
                        style={{ maxWidth: 360 }}
                        onFinish={onFinish}
                    >
                        <Form.Item
                            name="username"
                            rules={[{ required: true, message: '请输入用户名' }]}
                        >
                            <Input prefix={<UserOutlined />} placeholder="Username" />
                        </Form.Item>
                        <Form.Item
                            name="password"
                            rules={[{ required: true, message: '请输入密码' }]}
                        >
                            <Input prefix={<LockOutlined />} type="password" placeholder="Password" />
                        </Form.Item>
                        <Form.Item>
                            <Button block type="primary" htmlType="submit">
                                登录
                            </Button>
                        </Form.Item>
                    </Form>
                </Flex>
            </Flex>
        </Card>
    );
}

export default LoginPage;