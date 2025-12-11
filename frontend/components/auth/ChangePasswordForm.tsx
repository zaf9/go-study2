import { useState } from "react";
import { Button, Form, Input, Typography, message } from "antd";
import { useRouter } from "next/navigation";
import useAuth from "@/hooks/useAuth";

interface FormValues {
  oldPassword: string;
  newPassword: string;
  confirmPassword: string;
}

const { Title, Text } = Typography;

const formLayout = {
  labelCol: { span: 24 },
  wrapperCol: { span: 24 },
};

export default function ChangePasswordForm() {
  const router = useRouter();
  const { changePassword } = useAuth();
  const [submitting, setSubmitting] = useState(false);

  const handleFinish = async (values: FormValues) => {
    setSubmitting(true);
    try {
      await changePassword(values.oldPassword, values.newPassword);
      message.success("密码修改成功，请重新登录");
      router.replace("/login");
    } catch (error) {
      const reason =
        error instanceof Error ? error.message : "修改密码失败，请重试";
      message.error(reason);
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <div className="max-w-md w-full p-8 bg-white rounded-lg shadow">
      <Title level={3}>修改密码</Title>
      <Text type="secondary">请填写旧密码和符合策略的新密码</Text>
      <Form
        {...formLayout}
        layout="vertical"
        className="mt-6"
        onFinish={handleFinish}
        requiredMark={false}
      >
        <Form.Item
          label="旧密码"
          name="oldPassword"
          rules={[{ required: true, message: "请输入旧密码" }]}
        >
          <Input.Password placeholder="请输入旧密码" />
        </Form.Item>
        <Form.Item
          label="新密码"
          name="newPassword"
          rules={[
            { required: true, message: "请输入新密码" },
            { min: 8, message: "密码至少 8 位" },
            {
              pattern: /^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[^A-Za-z0-9]).{8,}$/,
              message: "需包含大写、小写、数字和特殊字符",
            },
          ]}
        >
          <Input.Password placeholder="请输入新密码" />
        </Form.Item>
        <Form.Item
          label="确认新密码"
          name="confirmPassword"
          dependencies={["newPassword"]}
          rules={[
            { required: true, message: "请再次输入新密码" },
            ({ getFieldValue }) => ({
              validator(_, value) {
                if (!value || getFieldValue("newPassword") === value) {
                  return Promise.resolve();
                }
                return Promise.reject(new Error("两次输入的密码不一致"));
              },
            }),
          ]}
        >
          <Input.Password placeholder="请再次输入新密码" />
        </Form.Item>
        <Form.Item>
          <Button
            type="primary"
            htmlType="submit"
            block
            size="large"
            loading={submitting}
          >
            保存新密码
          </Button>
        </Form.Item>
      </Form>
    </div>
  );
}

