/**
 * 测试 chapter-count.ts 工具函数
 * 验证从 static-routes.ts 获取章节数量的准确性
 */

import {
  getChapterCount,
  getAllChapterCounts,
  getTotalChapterCount,
} from '@/lib/chapter-count';
import { topicChapters } from '@/lib/static-routes';

describe('chapter-count utils', () => {
  describe('getChapterCount', () => {
    it('应该返回lexical_elements的正确章节数', () => {
      const count = getChapterCount('lexical_elements');
      expect(count).toBe(topicChapters.lexical_elements.length);
      expect(count).toBe(11); // 根据static-routes.ts的实际数据
    });

    it('应该返回constants的正确章节数', () => {
      const count = getChapterCount('constants');
      expect(count).toBe(topicChapters.constants.length);
      expect(count).toBe(12);
    });

    it('应该返回variables的正确章节数', () => {
      const count = getChapterCount('variables');
      expect(count).toBe(topicChapters.variables.length);
      expect(count).toBe(4);
    });

    it('应该返回types的正确章节数', () => {
      const count = getChapterCount('types');
      expect(count).toBe(topicChapters.types.length);
      expect(count).toBe(14);
    });

    it('应该返回properties的正确章节数', () => {
      const count = getChapterCount('properties');
      expect(count).toBe(topicChapters.properties.length);
      expect(count).toBe(7);
    });
  });

  describe('getAllChapterCounts', () => {
    it('应该返回所有主题的章节数映射', () => {
      const counts = getAllChapterCounts();

      expect(counts).toHaveProperty('lexical_elements');
      expect(counts).toHaveProperty('constants');
      expect(counts).toHaveProperty('variables');
      expect(counts).toHaveProperty('types');
      expect(counts).toHaveProperty('properties');

      expect(counts.lexical_elements).toBe(11);
      expect(counts.constants).toBe(12);
      expect(counts.variables).toBe(4);
      expect(counts.types).toBe(14);
      expect(counts.properties).toBe(7);
    });

    it('返回的数量应该与static-routes.ts一致', () => {
      const counts = getAllChapterCounts();
      
      Object.keys(counts).forEach((topic) => {
        const topicKey = topic as keyof typeof topicChapters;
        expect(counts[topicKey]).toBe(topicChapters[topicKey].length);
      });
    });
  });

  describe('getTotalChapterCount', () => {
    it('应该返回所有章节的总数', () => {
      const total = getTotalChapterCount();
      // 11 + 12 + 4 + 14 + 7 = 48
      expect(total).toBe(48);
    });

    it('总数应该等于各主题章节数之和', () => {
      const total = getTotalChapterCount();
      const sum = Object.values(topicChapters).reduce(
        (acc, chapters) => acc + chapters.length,
        0
      );
      expect(total).toBe(sum);
    });
  });
});
