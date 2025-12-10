'use client';

import { useEffect } from "react";
import ReactMarkdown from "react-markdown";
import remarkGfm from "remark-gfm";
import Prism from "prismjs";
import "prismjs/components/prism-go";
import "prismjs/components/prism-typescript";
import "prismjs/components/prism-javascript";
import "prismjs/components/prism-json";
import "prismjs/components/prism-bash";
import "prismjs/components/prism-markdown";
import "prismjs/themes/prism.css";
import { Card, Typography } from "antd";
import { ChapterContent as ChapterContentType } from "@/types/learning";

interface ChapterContentProps {
  content: ChapterContentType;
}

const { Title } = Typography;

export default function ChapterContent({ content }: ChapterContentProps) {
  useEffect(() => {
    Prism.highlightAll();
  }, [content.markdown]);

  return (
    <Card>
      <Title level={3}>{content.title}</Title>
      <ReactMarkdown
        remarkPlugins={[remarkGfm]}
        components={{
          code({ className, children, ...props }) {
            const language = className?.replace("language-", "") || "go";
            const codeText = String(children).trim();
            const html = Prism.highlight(
              codeText,
              Prism.languages[language] || Prism.languages.markup,
              language
            );
            return (
              <pre className={`language-${language}`}>
                <code dangerouslySetInnerHTML={{ __html: html }} {...props} />
              </pre>
            );
          },
        }}
      >
        {content.markdown || "暂无内容"}
      </ReactMarkdown>
    </Card>
  );
}


