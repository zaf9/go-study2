import api from "./api";
import { API_PATHS } from "./constants";
import { ChapterContent, ChapterSummary, TopicSummary } from "@/types/learning";

type TypesTopicContent = {
  concept?: {
    title?: string;
    summary?: string;
    rules?: string[];
    printableOutline?: string[];
  };
  rules?: Array<{ description?: string }>;
  examples?: Array<{
    title?: string;
    code?: string;
    expectedOutput?: string;
    isValid?: boolean;
  }>;
  references?: Array<{
    keyword?: string;
    summary?: string;
  }>;
};

function buildTypesMarkdown(content: TypesTopicContent): string {
  const lines: string[] = [];
  const concept = content.concept ?? {};
  const ruleTexts = [
    ...(concept.rules ?? []),
    ...(content.rules?.map((r) => r.description).filter(Boolean) ?? []),
  ];
  const outline = concept.printableOutline ?? [];
  const examples = content.examples ?? [];
  const references = content.references ?? [];

  if (concept.title) {
    lines.push(`# ${concept.title}`);
  }
  if (concept.summary) {
    lines.push(concept.summary);
  }
  if (outline.length > 0) {
    lines.push(
      ["## 提纲", ...outline.map((item) => `- ${item}`)].join("\n"),
    );
  }
  if (ruleTexts.length > 0) {
    lines.push(["## 规则", ...ruleTexts.map((r) => `- ${r}`)].join("\n"));
  }
  if (examples.length > 0) {
    const exampleBlocks = examples.map((ex) => {
      const code = ex.code ? `\n\`\`\`go\n${ex.code}\n\`\`\`\n` : "";
      const expected = ex.expectedOutput
        ? `\n> 输出：${ex.expectedOutput}`
        : "";
      const validity =
        typeof ex.isValid === "boolean"
          ? `\n> 是否符合规则：${ex.isValid ? "是" : "否"}`
          : "";
      return `### 示例：${ex.title ?? "示例"}${code}${expected}${validity}`;
    });
    lines.push(["## 示例", ...exampleBlocks].join("\n\n"));
  }
  if (references.length > 0) {
    lines.push(
      [
        "## 参考",
        ...references.map(
          (ref) => `- ${ref.keyword ?? ""}${ref.summary ? `：${ref.summary}` : ""}`,
        ),
      ].join("\n"),
    );
  }
  if (lines.length === 0) {
    return "";
  }
  return lines.join("\n\n");
}

export async function fetchTopics(): Promise<TopicSummary[]> {
  const data = await api.get<{
    topics: Array<{ id: string; title: string; description?: string }>;
  }>(API_PATHS.topics);
  const topics = Array.isArray(data?.topics) ? data.topics : null;
  if (!topics) {
    throw new Error("未能解析主题数据，请检查 API 地址配置是否指向后端 /api/v1");
  }
  return topics.map<TopicSummary>((item) => ({
    key: item.id as TopicSummary["key"],
    title: item.title,
    summary: item.description ?? "",
    chapterCount: 0,
  }));
}

export async function fetchChapters(topic: string): Promise<ChapterSummary[]> {
  const data = await api.get<{
    items: Array<{ id: number; title: string; name: string }>;
  }>(API_PATHS.topicMenu(topic));
  const items = data?.items ?? [];
  return items.map<ChapterSummary>((item) => ({
    id: item.name,
    topicKey: topic as ChapterSummary["topicKey"],
    title: item.title,
    order: item.id,
  }));
}

export async function fetchChapterContent(
  topic: string,
  chapter: string,
): Promise<ChapterContent> {
  const data = await api.get<{
    title?: string;
    content?: string | TypesTopicContent;
  }>(
    API_PATHS.chapterContent(topic, chapter),
  );
  const rawContent = data?.content;
  const title =
    data?.title ||
    (typeof rawContent === "object" && rawContent?.concept?.title) ||
    chapter;

  let markdown = "";
  if (typeof rawContent === "string") {
    markdown = rawContent;
  } else if (rawContent && typeof rawContent === "object") {
    markdown = buildTypesMarkdown(rawContent);
  }

  return {
    id: chapter,
    topicKey: topic as ChapterContent["topicKey"],
    title,
    markdown,
  };
}
