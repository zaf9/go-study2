"use client";

import { Drawer, Typography, message } from "antd";
import { useState, useEffect } from "react";
import useSWR from "swr";
import { fetchChapterContent } from "@/lib/learning";
import ReactMarkdown from "react-markdown";

const { Title, Paragraph } = Typography;

interface DeepLearningDrawerProps {
  open: boolean;
  onClose: () => void;
  topic: string;
  chapter: string;
  chapterTitle?: string;
}

export default function DeepLearningDrawer({
  open,
  onClose,
  topic,
  chapter,
  chapterTitle,
}: DeepLearningDrawerProps) {
  const [markdown, setMarkdown] = useState<string>("");

  const {
    error,
    isLoading,
  } = useSWR(
    open ? ["chapter-content", topic, chapter] : null,
    () => fetchChapterContent(topic, chapter),
    {
      onSuccess: (data) => {
        setMarkdown(data.markdown);
      },
      onError: (err) => {
        console.error("Failed to load chapter content:", err);
        message.error("加载章节内容失败");
      },
    }
  );

  useEffect(() => {
    if (!open) {
      setMarkdown("");
    }
  }, [open]);

  return (
    <Drawer
      title={
        <div className="flex items-center gap-2">
          <span className="text-base font-semibold">深入学习</span>
          <span className="text-sm text-gray-500">
            - {chapterTitle || chapter}
          </span>
        </div>
      }
      placement="right"
      open={open}
      onClose={onClose}
      width={720}
      styles={{
        body: { padding: "24px" },
      }}
    >
      {isLoading && (
        <div className="flex items-center justify-center py-12">
          <div className="text-center">
            <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-500 mx-auto mb-4"></div>
            <p className="text-gray-500">加载章节内容中...</p>
          </div>
        </div>
      )}

      {error && (
        <div className="text-center py-12">
          <p className="text-red-500 mb-2">加载失败</p>
          <p className="text-gray-500 text-sm">请稍后重试</p>
        </div>
      )}

      {!isLoading && !error && markdown && (
        <div className="prose prose-sm max-w-none">
          <div className="markdown-content">
            <ReactMarkdown
              components={{
                h1: ({ children }) => (
                  <Title level={3} className="mt-4 mb-2">
                    {children}
                  </Title>
                ),
                h2: ({ children }) => (
                  <Title level={4} className="mt-3 mb-2">
                    {children}
                  </Title>
                ),
                h3: ({ children }) => (
                  <Title level={5} className="mt-2 mb-1">
                    {children}
                  </Title>
                ),
                p: ({ children }) => (
                  <Paragraph className="mb-2">{children}</Paragraph>
                ),
                ul: ({ children }) => (
                  <ul className="list-disc pl-6 mb-2 space-y-1">{children}</ul>
                ),
                ol: ({ children }) => (
                  <ol className="list-decimal pl-6 mb-2 space-y-1">{children}</ol>
                ),
                li: ({ children }) => (
                  <li className="text-gray-700">{children}</li>
                ),
                code: ({ className, children, ...props }) => {
                  const match = /language-(\w+)/.exec(className || "");
                  return match ? (
                    <code
                      className="bg-gray-800 text-gray-100 px-2 py-1 rounded text-sm block overflow-x-auto"
                      {...props}
                    >
                      {children}
                    </code>
                  ) : (
                    <code
                      className="bg-gray-100 text-gray-800 px-1 py-0.5 rounded text-sm"
                      {...props}
                    >
                      {children}
                    </code>
                  );
                },
                pre: ({ children }) => (
                  <pre className="bg-gray-800 text-gray-100 p-4 rounded-md overflow-x-auto my-2">
                    {children}
                  </pre>
                ),
                blockquote: ({ children }) => (
                  <blockquote className="border-l-4 border-blue-500 pl-4 italic text-gray-600 my-2">
                    {children}
                  </blockquote>
                ),
                strong: ({ children }) => (
                  <strong className="font-semibold text-gray-900">{children}</strong>
                ),
              }}
            >
              {markdown}
            </ReactMarkdown>
          </div>
        </div>
      )}

      {!isLoading && !error && !markdown && (
        <div className="text-center py-12">
          <p className="text-gray-500">暂无章节内容</p>
        </div>
      )}
    </Drawer>
  );
}
