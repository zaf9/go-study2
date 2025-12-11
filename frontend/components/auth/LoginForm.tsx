"use client";

import { useState } from "react";
import { Button, Checkbox, Form, Input, Typography, message } from "antd";
import { useRouter } from "next/navigation";
import useAuth from "@/hooks/useAuth";

interface LoginFormValues {
  username: string;
  password: string;
  remember: boolean;
}

const { Title, Text } = Typography;

const formLayout = {
  labelCol: { span: 24 },
  wrapperCol: { span: 24 },
};

export default function LoginForm() {
  const router = useRouter();
  const { login } = useAuth();
  const [submitting, setSubmitting] = useState(false);

  const handleFinish = async (values: LoginFormValues) => {
    setSubmitting(true);
    try {
      const profile = await login(
        values.username,
        values.password,
        values.remember,
      );
      message.success("登录成功");
      if (profile.mustChangePassword) {
        router.push("/change-password");
        return;
      }
      router.push("/topics");
    } catch (error) {
      const reason =
        error instanceof Error ? error.message : "登录失败，请重试";
      message.error(reason);
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <div className="max-w-md w-full p-8 bg-white rounded-lg shadow">
      <Title level={3}>欢迎回来</Title>
      <Text type="secondary">使用账号登录以浏览学习主题</Text>
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
              pattern: /^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[^A-Za-z0-9]).{8,}$/,
              message: "需包含大写、小写、数字和特殊字符",
            },
          ]}
        >
          <Input.Password placeholder="请输入密码" />
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
            登录
          </Button>
        </Form.Item>
      </Form>
    </div>
  );
}
