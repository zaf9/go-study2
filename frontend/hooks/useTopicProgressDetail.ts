import { useEffect, useState } from "react";
import { getTopicProgress } from "@/services/progressService";
import { TopicProgressDetail } from "@/types/learning";

/**
 * 并行获取多个主题的详细进度信息
 * @param topicIds 主题ID数组
 * @returns 主题ID到详细进度信息的映射
 */
export default function useTopicProgressDetail(topicIds: string[]): Record<string, TopicProgressDetail> {
  // SWR不支持并行请求多个key，所以我们需要手动触发
  const [dataMap, setDataMap] = useState<Record<string, TopicProgressDetail>>({});

  useEffect(() => {
    if (topicIds.length === 0) {
      setDataMap({});
      return;
    }

    let cancelled = false;

    const fetchAllTopics = async () => {
      const promises = topicIds.map(async (topicId) => {
        try {
          const data = await getTopicProgress(topicId);
          return { topicId, data };
        } catch (error) {
          console.error(`Failed to fetch progress for ${topicId}:`, error);
          return { topicId, data: null };
        }
      });

      const results = await Promise.all(promises);

      if (!cancelled) {
        const map: Record<string, TopicProgressDetail> = {};
        results.forEach(({ topicId, data }) => {
          if (data) {
            map[topicId] = data;
          }
        });
        setDataMap(map);
      }
    };

    fetchAllTopics();

    return () => {
      cancelled = true;
    };
  }, [topicIds.join(',')]); // 使用join(',')将数组转换为字符串，避免数组引用变化导致的无限循环

  return dataMap;
}
