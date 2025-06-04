import { useEffect, useState } from 'react';
import { keyframes } from '@emotion/react';
import { Card, Row, Col, Typography, Button, Space, Avatar, Progress, Modal } from 'antd';
import { UserOutlined, TrophyOutlined, FireOutlined, TeamOutlined, RocketOutlined } from '@ant-design/icons';
import { useNavigate } from 'react-router-dom';
import styled from '@emotion/styled';
import request from "../axios/axios";

const { Title, Text } = Typography;

// 炫酷背景样式
const GradientBackground = styled.div`
  min-height: 100vh;
  background: linear-gradient(135deg, #1a1a2e 0%, #16213e 100%);
  padding: 2rem;
`;

// 动态卡片样式（调整最大宽度）
const HoverCard = styled(Card)`
  transition: all 0.3s ease;
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.1);
  max-width: 1200px;
  margin: 0 auto;
  &:hover {
    transform: translateY(-5px);
    box-shadow: 0 8px 24px rgba(0, 255, 255, 0.1);
  }
`;

// 脉冲动画
const pulse = keyframes`
  0% { transform: scale(1); }
  50% { transform: scale(1.05); }
  100% { transform: scale(1); }
`;

const PulseButton = styled(Button)`
  animation: ${pulse} 2s infinite;
  &:hover {
    animation: none;
    transform: scale(1.05);
  }
`;

interface User {
    userId: number,
    username: string,
    score: number,
    totalCount: number,
    winCount: number,
}

const Home: React.FC = () => {
    const [isMatching, setIsMatching] = useState(false);
    const [userData, setUserData] = useState({
        nickname: 'CyberPlayer',
        rating: 2450,
        totalGames: 128,
        wins: 89,
        avatar: 'https://randomuser.me/api/portraits/lego/1.jpg'
    });

    const navigate = useNavigate();
    // 项目初始化回调函数
    useEffect(() => {
        // 1. 获取当前登录用户信息
        request.get('/user/profile')
            .then((resp) => {
                if (resp.data && resp.data.code === 200) {
                    const userInfo: User = resp.data.data
                    setUserData((src) => ({
                        ...src,
                        nickname: userInfo.username,
                        totalGames: userInfo.totalCount,
                        wins: userInfo.winCount,
                        rating: userInfo.score,
                    }))
                } else {
                    console.log('获取用户信息失败！')
                }
            })
            .catch((error) => {
                console.log(error)
            })
    }, [])
    useEffect(() => {
        if (isMatching) {
            const timer = setTimeout(() => {
                Modal.info({
                    title: '匹配成功!',
                    content: '正在进入对战房间...',
                    onOk: () => navigate('/game')
                });
            }, 2000);
            return () => clearTimeout(timer);
        }
    }, [isMatching]);

    const winRate = Math.round((userData.totalGames === 0 ? 0 : (userData.wins / userData.totalGames)) * 100);

    // 建立 WebSocket 连接
    const ws = new WebSocket("ws://localhost:8081/match/findMatch");

    ws.onopen = () => {
        console.log("连接建立成功...");
    }
    ws.onclose = () => {
        console.log("连接已断开...");
    }
    ws.onerror = (e) => {
        console.log("连接发生错误");
    }
    ws.onmessage = (e) => {
        // 响应反序列化
        const resp = JSON.parse(e.data);
        if (resp?.code === 200) {
            if (resp.message === "startMatch") {
                // 服务器返回开始匹配
                // 设置状态为开始匹配
                setIsMatching(true);
            } else {
                // 服务器返回结束匹配
                // 设置状态为匹配中
                setIsMatching(false);
            }
        } else {
            console.log("发生了异常" + resp.message);
        }
    }
    // 处理开始匹配逻辑
    const handleMatch = () => {
        // 判断到底是开始匹配还是取消匹配
        if (!isMatching) {
            // 开始匹配
            // 向后端发送请求
            ws.send(JSON.stringify({
                "message": "startMatch"
            }))
        } else {
            // 取消匹配
            // 向后端发送请求
            ws.send(JSON.stringify({
                "message": "stopMatch"
            }))
        }
    }

    return (
        <GradientBackground>
            <Title level={2} style={{ color: '#fff', textAlign: 'center', marginBottom: '2rem' }}>
                <RocketOutlined /> 五子棋竞技场
            </Title>

            <Row gutter={[24, 24]} justify="center">
                <Col xs={24} xl={18}>
                    <HoverCard>
                        {/* 用户信息主容器 */}
                        <div style={{
                            display: 'flex',
                            alignItems: 'center',
                            gap: 32,
                            padding: 24,
                            background: 'rgba(255, 255, 255, 0.05)',
                            borderRadius: 8
                        }}>
                            {/* 头像和昵称 */}
                            <div style={{
                                display: 'flex',
                                alignItems: 'center',
                                flexShrink: 0,
                                minWidth: 200
                            }}>
                                <Avatar
                                    size={64}
                                    src={userData.avatar}
                                    style={{
                                        border: '2px solid #00f3ff',
                                        marginRight: 16,
                                        boxShadow: '0 0 12px rgba(0, 243, 255, 0.3)'
                                    }}
                                />
                                <Title
                                    level={4}
                                    style={{
                                        color: '#fff',
                                        margin: 0,
                                        whiteSpace: 'nowrap',
                                        fontSize: '1.25rem'
                                    }}
                                >
                                    <UserOutlined style={{ marginRight: 8 }} />
                                    {userData.nickname}
                                </Title>
                            </div>

                            {/* 数据指标区 */}
                            <div style={{
                                display: 'flex',
                                gap: 32,
                                flexWrap: 'wrap',
                                flexGrow: 1,
                                minWidth: 0
                            }}>
                                {/* 天梯积分 */}
                                <div style={{
                                    display: 'flex',
                                    alignItems: 'center',
                                    gap: 8,
                                    padding: '8px 16px',
                                    background: 'rgba(255, 215, 0, 0.1)',
                                    borderRadius: 6,
                                    flexShrink: 0
                                }}>
                                    <TrophyOutlined style={{ color: '#ffd700', fontSize: 20 }} />
                                    <div>
                                        <Text strong style={{ color: '#7ed6df', display: 'block' }}>天梯积分</Text>
                                        <Text strong style={{ color: '#ffd700', fontSize: 18 }}>{userData.rating}</Text>
                                    </div>
                                </div>

                                {/* 总对局 */}
                                <div style={{
                                    display: 'flex',
                                    alignItems: 'center',
                                    gap: 8,
                                    padding: '8px 16px',
                                    background: 'rgba(0, 243, 255, 0.1)',
                                    borderRadius: 6,
                                    flexShrink: 0
                                }}>
                                    <TeamOutlined style={{ color: '#00f3ff', fontSize: 20 }} />
                                    <div>
                                        <Text strong style={{ color: '#7ed6df', display: 'block' }}>总对局</Text>
                                        <Text strong style={{ color: '#00f3ff', fontSize: 18 }}>{userData.totalGames}</Text>
                                    </div>
                                </div>

                                {/* 胜率 */}
                                <div style={{
                                    display: 'flex',
                                    alignItems: 'center',
                                    gap: 8,
                                    padding: '8px 16px',
                                    background: 'rgba(135, 208, 104, 0.1)',
                                    borderRadius: 6,
                                    flexShrink: 0
                                }}>
                                    <FireOutlined style={{ color: '#87d068', fontSize: 20 }} />
                                    <div>
                                        <Text strong style={{ color: '#7ed6df', display: 'block' }}>胜率</Text>
                                        <Text strong style={{ color: '#87d068', fontSize: 18 }}>{winRate}%</Text>
                                    </div>
                                </div>
                            </div>
                        </div>
                    </HoverCard>
                </Col>

                {/* 匹配区域 */}
                <Col xs={24} xl={18}>
                    <HoverCard style={{ height: '100%' }}>
                        <Space direction="vertical" style={{ width: '100%', textAlign: 'center' }}>
                            <Title level={3} style={{ color: '#fff' }}>
                                快速匹配
                            </Title>
                            <Text type="secondary" style={{ color: '#aaa' }}>
                                点击按钮开始寻找对手
                            </Text>

                            <div style={{ margin: '2rem 0' }}>
                                <PulseButton
                                    type="primary"
                                    size="large"
                                    shape="round"
                                    icon={<FireOutlined />}
                                    onClick={handleMatch}
                                    style={{
                                        background: 'linear-gradient(45deg, #ff6b6b, #ff8e53)',
                                        border: 'none',
                                        height: '60px',
                                        width: '200px',
                                        fontSize: '1.2rem'
                                    }}
                                >
                                    {isMatching ? '匹配中（点击取消）' : '开始匹配'}
                                </PulseButton>
                            </div>

                            <Row gutter={16} style={{ marginTop: '2rem' }}>
                                <Col span={12}>
                                    <Card
                                        size="small"
                                        style={{
                                            background: 'rgba(0, 255, 255, 0.1)',
                                            border: '1px solid rgba(0, 255, 255, 0.2)'
                                        }}
                                    >
                                        <Text strong style={{ color: '#00f3ff' }}>在线玩家</Text>
                                        <Title level={3} style={{ color: '#fff', margin: 0 }}>1,234</Title>
                                    </Card>
                                </Col>
                                <Col span={12}>
                                    <Card
                                        size="small"
                                        style={{
                                            background: 'rgba(255, 107, 107, 0.1)',
                                            border: '1px solid rgba(255, 107, 107, 0.2)'
                                        }}
                                    >
                                        <Text strong style={{ color: '#ff6b6b' }}>正在对局</Text>
                                        <Title level={3} style={{ color: '#fff', margin: 0 }}>567</Title>
                                    </Card>
                                </Col>
                            </Row>
                        </Space>
                    </HoverCard>
                </Col>
            </Row>
        </GradientBackground>
    );
};

export default Home;