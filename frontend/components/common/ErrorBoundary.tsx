'use client';
import React from "react";
import { Result, Button } from "antd";

interface ErrorBoundaryState {
  hasError: boolean;
}

// 捕获渲染错误并提供重试入口
class ErrorBoundary extends React.Component<React.PropsWithChildren, ErrorBoundaryState> {
  constructor(props: React.PropsWithChildren) {
    super(props);
    this.state = { hasError: false };
  }

  static getDerivedStateFromError() {
    return { hasError: true };
  }

  componentDidCatch(error: Error, info: React.ErrorInfo) {
    // 记录错误以便后续排查
    // eslint-disable-next-line no-console
    console.error("组件渲染错误", error, info);
  }

  handleRetry = () => {
    this.setState({ hasError: false });
  };

  render() {
    if (this.state.hasError) {
      return (
        <Result
          status="error"
          title="页面加载失败"
          subTitle="请重试或刷新页面。"
          extra={
            <Button type="primary" onClick={this.handleRetry}>
              重试
            </Button>
          }
        />
      );
    }
    return this.props.children;
  }
}

export default ErrorBoundary;

