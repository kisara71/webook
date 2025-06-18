import React from 'react';
import { Button, Form, Input } from 'antd';
import axios from "@/axios/axios";
import Link from "next/link";
import router from "next/router";
import styles from './signup.module.css';

const onFinish = (values: any) => {
    axios.post("/users/signup", values)
        .then((res) => {
            if(res.status != 200) {
                alert(res.statusText);
                return
            }
            if(typeof res.data == 'string') {
                alert(res.data);
            } else {
                const msg = res.data?.msg || JSON.stringify(res.data)
                alert(msg);
                if(res.data.code == 0) {
                    router.push('/users/login')
                }
            }

        }).catch((err) => {
            alert(err);
    })
};

const onFinishFailed = (errorInfo: any) => {
    alert("输入有误")
};

const SignupForm: React.FC = () => (
    <div className={styles.container}>
        <div className={styles.formWrapper}>
            <h1 className={styles.title}>注册</h1>
            <Form
                name="basic"
                layout="vertical"
                style={{ maxWidth: 400, width: '100%' }}
                initialValues={{ remember: true }}
                onFinish={onFinish}
                onFinishFailed={onFinishFailed}
                autoComplete="off"
                className={styles.form}
            >
                <Form.Item
                    name="email"
                    rules={[{ required: true, message: '请输入邮箱' }]}
                >
                    <Input placeholder="邮箱" size="large" />
                </Form.Item>

                <Form.Item
                    name="password"
                    rules={[{ required: true, message: '请输入密码' }]}
                >
                    <Input.Password placeholder="密码" size="large" />
                </Form.Item>

                <Form.Item
                    name="confirmPassword"
                    rules={[{ required: true, message: '请确认密码' }]}
                >
                    <Input.Password placeholder="确认密码" size="large" />
                </Form.Item>

                <Form.Item>
                    <Button type="primary" htmlType="submit" block size="large" className={styles.submitButton}>
                        注册
                    </Button>
                </Form.Item>

                <div className={styles.links}>
                    <Link href={"/users/login"} className={styles.link}>
                        已有账号？立即登录
                    </Link>
                </div>
            </Form>
        </div>
    </div>
);

export default SignupForm;