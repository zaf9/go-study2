import { Alert } from "antd";

interface ErrorMessageProps {
  message?: string;
  description?: string;
}

// 简单错误提示组件
const ErrorMessage = ({
  message = "发生错误",
  description,
}: ErrorMessageProps) => (
  <Alert type="error" showIcon message={message} description={description} />
);

export default ErrorMessage;
