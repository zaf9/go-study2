import { Spin } from "antd";

// 全局加载状态组件
const Loading = () => (
  <div className="flex items-center justify-center py-12">
    <Spin tip="加载中..." />
  </div>
);

export default Loading;

