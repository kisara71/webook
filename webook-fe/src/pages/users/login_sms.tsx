import React, { useState } from 'react';
import { Button, Form, Input } from 'antd';
import axios from "@/axios/axios";
import router from "next/router";
import Link from "next/link";
import styles from './login_sms.module.css';

const onFinish = (values: any) => {
    axios.post("/users/login_sms", values)
        .then((res) => {
            if(res.status != 200) {
                alert(res.statusText);
                return
            }

            if (res.data.code == 0) {
                router.push('/users/profile')
                return;
            }
            alert(res.data.msg)
        }).catch((err) => {
        alert(err);
    })
};

const onFinishFailed = (errorInfo: any) => {
    alert("输入有误")
};

const LoginFormSMS: React.FC = () => {
    const [form] = Form.useForm();
    const [countdown, setCountdown] = useState(0);

    const sendCode = () => {
        const data = form.getFieldValue("phone")
        if (!data) {
            alert("请先输入手机号码");
            return;
        }
        axios.post("/users/login_sms/code/send", {"phone": data} ).then((res) => {
            if(res.status != 200) {
                alert(res.statusText);
                return
            }
            alert(res?.data?.msg || "系统错误，请重试");
            // 开始倒计时
            setCountdown(60);
            const timer = setInterval(() => {
                setCountdown((prev) => {
                    if (prev <= 1) {
                        clearInterval(timer);
                        return 0;
                    }
                    return prev - 1;
                });
            }, 1000);
        }).catch((err) => {
            alert(err);
        })
    }

    return (
        <div className={styles.container}>
            <div className={styles.formWrapper}>
                <h1 className={styles.title}>手机号登录</h1>
                <Form
                    name="basic"
                    layout="vertical"
                    style={{ maxWidth: 400, width: '100%' }}
                    initialValues={{ remember: true }}
                    onFinish={onFinish}
                    onFinishFailed={onFinishFailed}
                    autoComplete="off"
                    form={form}
                    className={styles.form}
                >
                    <Form.Item
                        name="phone"
                        rules={[{ required: true, message: '请输入手机号码' }]}
                    >
                        <Input placeholder="手机号码" size="large" />
                    </Form.Item>

                    <Form.Item
                        name="code"
                        rules={[{ required: true, message: '请输入验证码' }]}
                    >
                        <div className={styles.codeInputWrapper}>
                            <Input placeholder="验证码" size="large" />
                            <Button 
                                type="default" 
                                onClick={sendCode} 
                                disabled={countdown > 0}
                                className={styles.sendCodeButton}
                            >
                                {countdown > 0 ? `${countdown}秒后重试` : '发送验证码'}
                            </Button>
                        </div>
                    </Form.Item>

                    <Form.Item>
                        <Button type="primary" htmlType="submit" block size="large" className={styles.submitButton}>
                            登录/注册
                        </Button>
                    </Form.Item>

                    <div className={styles.links}>
                        <Link href={"/users/login"} className={styles.link}>
                            使用邮箱登录
                        </Link>
                    </div>
                </Form>
            </div>
        </div>
    );
};

export default LoginFormSMS;