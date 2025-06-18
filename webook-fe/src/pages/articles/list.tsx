'use client';
import { EditOutlined } from '@ant-design/icons';
import { ProLayout, ProList } from '@ant-design/pro-components';
import { Button, Tag, message } from 'antd';
import React, { useEffect, useState } from 'react';
import axios from "@/axios/axios";
import { useRouter } from "next/router";
import type { ArticleItem } from './model';

const IconText = ({ icon, text, onClick }: { icon: any; text: string; onClick: () => void }) => (
    <Button 
        onClick={onClick} 
        type="default" 
        className="flex items-center hover:bg-gray-50 transition-colors"
    >
        {React.createElement(icon, { style: { marginRight: 8 } })}
        {text}
    </Button>
);

const ArticleList: React.FC = () => {
    const [data, setData] = useState<ArticleItem[]>([]);
    const [loading, setLoading] = useState(true);
    const router = useRouter();

    useEffect(() => {
        const fetchArticles = async () => {
            try {
                const res = await axios.post('/articles/list', {
                    offset: 0,
                    limit: 100,
                });
                
                if (res.data.code === 0) {
                    setData(res.data.data);
                } else {
                    message.error(res.data.msg || '获取文章列表失败');
                }
            } catch (err) {
                message.error('获取文章列表失败，请稍后重试');
            } finally {
                setLoading(false);
            }
        };

        fetchArticles();
    }, []);

    const getStatusTag = (status: number) => {
        switch (status) {
            case 1:
                return <Tag color="default" className="px-2 py-1">未发表</Tag>;
            case 2:
                return <Tag color="success" className="px-2 py-1">已发表</Tag>;
            case 3:
                return <Tag color="warning" className="px-2 py-1">仅自己可见</Tag>;
            default:
                return null;
        }
    };

    return (
        <ProLayout 
            title="创作中心"
            className="min-h-screen bg-gray-50"
        >
            <div className="max-w-6xl mx-auto px-4 py-8">
                <ProList<ArticleItem>
                    toolBarRender={() => [
                        <Button 
                            key="create" 
                            type="primary" 
                            href="/articles/edit"
                            className="h-10 px-6 text-base"
                        >
                            写作
                        </Button>,
                    ]}
                    itemLayout="vertical"
                    rowKey="id"
                    headerTitle="文章列表"
                    loading={loading}
                    dataSource={data}
                    metas={{
                        title: {
                            dataIndex: "title",
                            render: (title) => (
                                <span className="text-xl font-medium text-gray-900">{title}</span>
                            ),
                        },
                        description: {
                            render: (_, record) => getStatusTag(record.status),
                        },
                        actions: {
                            render: (_, record) => [
                                <IconText
                                    icon={EditOutlined}
                                    text="编辑"
                                    onClick={() => router.push(`/articles/edit?id=${record.id}`)}
                                    key="edit"
                                />,
                            ],
                        },
                        content: {
                            render: (_, record) => (
                                <div 
                                    className="prose max-w-none text-gray-600"
                                    dangerouslySetInnerHTML={{ __html: record.abstract }}
                                />
                            ),
                        },
                    }}
                    pagination={{
                        pageSize: 10,
                        className: "mt-4",
                    }}
                    className="bg-white rounded-lg shadow-sm"
                />
            </div>
        </ProLayout>
    );
};

export default ArticleList;
