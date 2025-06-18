import dynamic from 'next/dynamic';
import { Button, Form, Input, message } from "antd";
import { useEffect, useState } from "react";
import axios from "@/axios/axios";
import { useRouter, useSearchParams } from "next/navigation";
import { ProLayout } from "@ant-design/pro-components";
import type { ArticleItem, ArticleCreateRequest, ArticleUpdateRequest } from './model';

const WangEditor = dynamic(
    () => import('../../components/editor'),
    { ssr: false },
);

const ArticleEdit: React.FC = () => {
    const [form] = Form.useForm();
    const [html, setHtml] = useState<string>('');
    const [loading, setLoading] = useState(false);
    const router = useRouter();
    const searchParams = useSearchParams();
    const articleId = searchParams?.get("id");

    useEffect(() => {
        const fetchArticle = async () => {
            if (!articleId) return;

            try {
                const res = await axios.get(`/articles/detail/${articleId}`);
                if (res.data.code === 0) {
                    const article = res.data.data;
                    form.setFieldsValue({
                        title: article.title,
                        abstract: article.abstract,
                    });
                    setHtml(article.content);
                } else {
                    message.error(res.data.msg || '获取文章失败');
                }
            } catch (err) {
                message.error('获取文章失败，请稍后重试');
            }
        };

        fetchArticle();
    }, [articleId, form]);

    const handleSubmit = async (values: ArticleCreateRequest | ArticleUpdateRequest) => {
        try {
            setLoading(true);
            const requestData = {
                ...values,
                content: html,
                ...(articleId && { id: parseInt(articleId) }),
            };

            const res = await axios.post("/articles/edit", requestData);
            
            if (res.data.code === 0) {
                message.success('保存成功');
                router.push('/articles/list');
            } else {
                message.error(res.data.msg || '保存失败');
            }
        } catch (err) {
            message.error('保存失败，请稍后重试');
        } finally {
            setLoading(false);
        }
    };

    const handlePublish = async () => {
        try {
            setLoading(true);
            const values = form.getFieldsValue();
            const requestData = {
                ...values,
                content: html,
                ...(articleId && { id: parseInt(articleId) }),
            };

            const res = await axios.post("/articles/publish", requestData);
            
            if (res.data.code === 0) {
                message.success('发表成功');
                router.push(`/articles/view?id=${res.data.data}`);
            } else {
                message.error(res.data.msg || '发表失败');
            }
        } catch (err) {
            message.error('发表失败，请稍后重试');
        } finally {
            setLoading(false);
        }
    };

    return (
        <ProLayout 
            title="创作中心"
            className="min-h-screen bg-gray-50"
        >
            <div className="max-w-4xl mx-auto px-4 py-8">
                <div className="bg-white rounded-lg shadow-sm p-6">
                    <Form
                        form={form}
                        onFinish={handleSubmit}
                        layout="vertical"
                        className="space-y-6"
                    >
                        <Form.Item
                            name="title"
                            rules={[{ required: true, message: '请输入标题' }]}
                        >
                            <Input 
                                placeholder="请输入标题"
                                size="large"
                                className="text-xl h-12"
                            />
                        </Form.Item>

                        <Form.Item
                            name="abstract"
                            rules={[{ required: true, message: '请输入摘要' }]}
                        >
                            <Input.TextArea
                                placeholder="请输入摘要"
                                autoSize={{ minRows: 2, maxRows: 4 }}
                                className="text-base"
                            />
                        </Form.Item>

                        <div className="border rounded-lg p-4 bg-white">
                            <WangEditor html={html} setHtmlFn={setHtml} />
                        </div>

                        <div className="flex justify-end space-x-4 pt-4 border-t">
                            <Button 
                                type="default" 
                                onClick={handlePublish}
                                loading={loading}
                                className="h-10 px-6"
                            >
                                发表
                            </Button>
                            <Button 
                                type="primary" 
                                htmlType="submit"
                                loading={loading}
                                className="h-10 px-6"
                            >
                                保存
                            </Button>
                        </div>
                    </Form>
                </div>
            </div>
        </ProLayout>
    );
};

export default ArticleEdit;