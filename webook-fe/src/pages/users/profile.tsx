import { ProDescriptions } from '@ant-design/pro-components';
import React, { useState, useEffect } from 'react';
import { Button, Spin } from 'antd';
import axios from "@/axios/axios";
import styles from './profile.module.css';

function Page() {
    let p: Profile = {Email: "", Phone: "", Nickname: "", Birthday:"", AboutMe: ""}
    const [data, setData] = useState<Profile>(p)
    const [isLoading, setLoading] = useState(false)

    useEffect(() => {
        setLoading(true)
        axios.get('/users/profile')
            .then((res) => res.data)
            .then((data) => {
                setData(data)
                setLoading(false)
            })
    }, [])

    if (isLoading) return (
        <div className={styles.loadingContainer}>
            <Spin size="large" />
        </div>
    )
    if (!data) return <p>No profile data</p>

    return (
        <div className={styles.container}>
            <div className={styles.profileCard}>
                <h1 className={styles.title}>个人信息</h1>
                <ProDescriptions
                    column={1}
                    className={styles.descriptions}
                >
                    <ProDescriptions.Item label="昵称" valueType="text" className={styles.item}>
                        {data.Nickname}
                    </ProDescriptions.Item>
                    <ProDescriptions.Item
                        valueType="text"
                        label="邮箱"
                        className={styles.item}
                    >
                        {data.Email}
                    </ProDescriptions.Item>
                    <ProDescriptions.Item
                        valueType="text"
                        label="手机"
                        className={styles.item}
                    >
                        {data.Phone}
                    </ProDescriptions.Item>
                    <ProDescriptions.Item 
                        label="生日" 
                        valueType="date"
                        className={styles.item}
                    >
                        {data.Birthday}
                    </ProDescriptions.Item>
                    <ProDescriptions.Item
                        valueType="text"
                        label="关于我"
                        className={styles.item}
                    >
                        {data.AboutMe}
                    </ProDescriptions.Item>
                    <ProDescriptions.Item className={styles.buttonContainer}>
                        <Button 
                            href={"/users/edit"} 
                            type="primary"
                            className={styles.editButton}
                        >
                            修改资料
                        </Button>
                    </ProDescriptions.Item>
                </ProDescriptions>
            </div>
        </div>
    )
}

export default Page