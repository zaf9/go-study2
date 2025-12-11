import { TopicKey } from "@/types/learning";

// 预定义可导出的主题与章节，需与后端内容保持一致
export const topics: TopicKey[] = [
  "lexical_elements",
  "constants",
  "variables",
  "types",
];

export const topicChapters: Record<TopicKey, string[]> = {
  lexical_elements: [
    "comments",
    "tokens",
    "semicolons",
    "identifiers",
    "keywords",
    "operators",
    "integers",
    "floats",
    "imaginary",
    "runes",
    "strings",
  ],
  constants: [
    "boolean",
    "rune",
    "integer",
    "floating_point",
    "complex",
    "string",
    "expressions",
    "typed_untyped",
    "conversions",
    "builtin_functions",
    "iota",
    "implementation_restrictions",
  ],
  variables: ["storage", "static", "dynamic", "zero"],
  types: [
    "boolean",
    "numeric",
    "string",
    "array",
    "slice",
    "struct",
    "pointer",
    "function",
    "interface_basic",
    "interface_embedded",
    "interface_general",
    "interface_impl",
    "map",
    "channel",
  ],
};

export function buildTopicParams() {
  return topics.map((topic) => ({ topic }));
}

export function buildTopicChapterParams() {
  return topics.flatMap((topic) =>
    (topicChapters[topic] || []).map((chapter) => ({ topic, chapter })),
  );
}
