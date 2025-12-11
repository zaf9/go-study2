"use client";

import { useState } from "react";
import { Button, Checkbox, Form, Input, Typography, message } from "antd";
import { useRouter } from "next/navigation";
import useAuth from "@/hooks/useAuth";

interface RegisterFormValues {
  username: string;
  password: string;
  confirm: string;
  remember: boolean;
}

const { Title, Text } = Typography;

const formLayout = {
  labelCol: { span: 24 },
  wrapperCol: { span: 24 },
};

export default function RegisterForm() {
  const router = useRouter();
  const { register } = useAuth();
  const [submitting, setSubmitting] = useState(false);

  const handleFinish = async (values: RegisterFormValues) => {
    setSubmitting(true);
    try {
      await register(values.username, values.password, values.remember);
      message.success("注册成功，已自动登录");
      router.push("/topics");
    } catch (error) {
      const reason =
        error instanceof Error ? error.message : "注册失败，请重试";
      message.error(reason);
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <div className="max-w-md w-full p-8 bg-white rounded-lg shadow">
      <Title level={3}>创建账号</Title>
      <Text type="secondary">注册后即可浏览全部学习主题</Text>
      <Form
        {...formLayout}
        layout="vertical"
        className="mt-6"
        initialValues={{ remember: true }}
        onFinish={handleFinish}
        requiredMark={false}
      >
        <Form.Item
          label="用户名"
          name="username"
          rules={[
            { required: true, message: "请输入用户名" },
            { min: 3, message: "用户名至少 3 个字符" },
          ]}
        >
          <Input placeholder="请输入用户名" allowClear />
        </Form.Item>
        <Form.Item
          label="密码"
          name="password"
          rules={[
            { required: true, message: "请输入密码" },
            { min: 8, message: "密码至少 8 位" },
            {
              pattern: /^(?=.*[a-z])(?=.*[A-Z])(?=.*\d).{8,}$/,
              message: "需包含大小写字母与数字",
            },
          ]}
        >
          <Input.Password placeholder="请输入密码" />
        </Form.Item>
        <Form.Item
          label="确认密码"
          name="confirm"
          dependencies={["password"]}
          rules={[
            { required: true, message: "请再次输入密码" },
            ({ getFieldValue }) => ({
              validator(_, value) {
                if (!value || getFieldValue("password") === value) {
                  return Promise.resolve();
                }
                return Promise.reject(new Error("两次输入的密码不一致"));
              },
            }),
          ]}
        >
          <Input.Password placeholder="请确认密码" />
        </Form.Item>
        <Form.Item name="remember" valuePropName="checked">
          <Checkbox>记住我</Checkbox>
        </Form.Item>
        <Form.Item>
          <Button
            type="primary"
            htmlType="submit"
            block
            size="large"
            loading={submitting}
          >
            注册并登录
          </Button>
        </Form.Item>
      </Form>
    </div>
  );
}
