"use client";

import { Layout, Typography } from "antd";

const { Footer: AntFooter } = Layout;
const { Text } = Typography;

export default function Footer() {
    return (
        <AntFooter className="bg-white px-0 py-4 text-center">
            <div className="mx-auto max-w-screen-2xl px-4 sm:px-6 lg:px-8">
                <Text type="secondary">Go Study 2 · 学习与实践并重</Text>
            </div>
        </AntFooter>
    );
}
