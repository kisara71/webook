import React, { useState, useEffect } from 'react';
import axios from "@/axios/axios";
import { redirect } from "next/navigation";

function Page() {
    const [isLoading, setLoading] = useState(false)

    useEffect(() => {
        setLoading(true)
        axios.get('/oauth2/github/authurl')
            .then((res) => res.data)
            .then((data) => {
                setLoading(false)
                if(data && data.data) {
                    window.location.href = data.data
                }
            })
    }, [])

    if (isLoading) return <p>Loading...</p>

    return (
        <div>
            <h1>GitHub 登录</h1>
            <p>正在跳转到 GitHub 授权页面...</p>
        </div>
    )
}

export default Page
