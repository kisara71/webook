import React, { useState, useEffect } from 'react';
import axios from "@/axios/axios";
import { Button, Modal, QRCode, Typography, message, Spin } from "antd";
import { ProLayout } from "@ant-design/pro-components";
import { EyeOutlined, LikeOutlined, MoneyCollectOutlined, StarOutlined } from "@ant-design/icons";
import { useSearchParams } from "next/navigation";
import type { ArticleDetail, CodeURL } from './model';

export const dynamic = 'force-dynamic';

const ArticleView: React.FC = () => {
    const [article, setArticle] = useState<ArticleDetail | null>(null);
    const [loading, setLoading] = useState(true);
    const [rewardLoading, setRewardLoading] = useState(false);
    const [openQRCode, setOpenQRCode] = useState(false);
    const [codeURL, setCodeURL] = useState('');
    const [rewardId, setRewardId] = useState(0);
    const searchParams = useSearchParams();
    const articleId = searchParams?.get('id') || '1';

    useEffect(() => {
        const fetchArticle = async () => {
            try {
                const res = await axios.get(`/articles/pub/${articleId}`);
                if (res.data.code === 0) {
                    setArticle(res.data.data);
                    } else {
                    message.error(res.data.msg || '获取文章失败');
                }
            } catch (err) {
                message.error('获取文章失败，请稍后重试');
            } finally {
                setLoading(false);
            }
        };

        fetchArticle();
    }, [articleId]);

    const handleLike = async () => {
        if (!article) return;

        try {
            const res = await axios.post('/articles/pub/like', {
                id: parseInt(articleId),
                like: !article.liked
            });

            if (res.data.code === 0) {
                setArticle(prev => {
                    if (!prev) return null;
                    return {
                        ...prev,
                        liked: !prev.liked,
                        likeCnt: prev.liked ? prev.likeCnt - 1 : prev.likeCnt + 1
                    };
                });
            } else {
                message.error(res.data.msg || '操作失败');
                }
        } catch (err) {
            message.error('操作失败，请稍后重试');
    }
    };

    const handleCollect = async () => {
        if (!article || article.collected) return;

        try {
            const res = await axios.post('/articles/pub/collect', {
                id: parseInt(articleId),
            cid: 0,
            });

            if (res.data.code === 0) {
                setArticle(prev => {
                    if (!prev) return null;
                    return {
                        ...prev,
                        collected: true,
                        collectCnt: prev.collectCnt + 1
                    };
                });
            } else {
                message.error(res.data.msg || '收藏失败');
                }
        } catch (err) {
            message.error('收藏失败，请稍后重试');
        }
    };

    const handleReward = async () => {
        try {
            setRewardLoading(true);
            const res = await axios.post<{ code: number; msg: string; data: CodeURL }>('/articles/pub/reward', {
                id: parseInt(articleId),
            amt: 1,
            });

            if (res.data.code === 0) {
                setCodeURL(res.data.data.codeURL);
                setRewardId(res.data.data.rid);
                setOpenQRCode(true);
            } else {
                message.error(res.data.msg || '生成打赏码失败');
            }
        } catch (err) {
            message.error('生成打赏码失败，请稍后重试');
        } finally {
            setRewardLoading(false);
        }
    };

    const handleCloseModal = async () => {
        setOpenQRCode(false);
        if (rewardId > 0) {
            try {
                const res = await axios.post<{ code: number; msg: string; data: string }>('/reward/detail', {
                    rid: rewardId,
                });

                if (res.data.code === 0 && res.data.data === 'RewardStatusPayed') {
                    message.success('打赏成功');
                }
            } catch (err) {
                console.error('查询打赏状态失败:', err);
            }
        }
    };

    if (loading) {
        return (
            <div className="min-h-screen flex items-center justify-center">
                <Spin size="large" tip="加载中..." />
            </div>
        );
    }

    if (!article) {
        return (
            <div className="min-h-screen flex items-center justify-center">
                <Typography.Text type="secondary">文章不存在或已被删除</Typography.Text>
            </div>
        );
    }

    return (
        <ProLayout pure>
            <div className="max-w-4xl mx-auto p-6">
                <article className="bg-white rounded-lg shadow p-8">
                    <Typography.Title level={2} className="mb-6">
                        {article.title}
                </Typography.Title>

                    <div 
                        className="prose max-w-none mb-8"
                        dangerouslySetInnerHTML={{ __html: article.content }}
                    />

                    <div className="flex items-center space-x-4 border-t pt-4">
                        <Button icon={<EyeOutlined />} type="text">
                            {article.readCnt}
                        </Button>
                        <Button 
                            icon={<MoneyCollectOutlined />} 
                            onClick={handleReward}
                            loading={rewardLoading}
                        >
                            打赏一分钱
                        </Button>
                        <Button 
                            icon={<LikeOutlined style={{ color: article.liked ? 'red' : undefined }} />}
                            onClick={handleLike}
                            type={article.liked ? 'primary' : 'default'}
                        >
                            {article.likeCnt}
                        </Button>
                        <Button 
                            icon={<StarOutlined style={{ color: article.collected ? 'red' : undefined }} />}
                            onClick={handleCollect}
                            disabled={article.collected}
                            type={article.collected ? 'primary' : 'default'}
                        >
                            {article.collectCnt}
                        </Button>
                    </div>
                </article>
            </div>

            <Modal 
                title="扫描二维码" 
                open={openQRCode} 
                onCancel={handleCloseModal} 
                onOk={handleCloseModal}
                centered
            >
                <div className="flex flex-col items-center p-4">
                    <QRCode value={codeURL} size={256} />
                    <Typography.Text className="mt-4 text-gray-500">
                        请使用微信扫描二维码完成打赏
                    </Typography.Text>
                </div>
            </Modal>
        </ProLayout>
    );
};

export default ArticleView;