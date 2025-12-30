/**
 * 章节状态计算工具测试
 * 测试章节状态的计算逻辑
 */

import {
  getChapterStatus,
  getChapterStatusText,
  getChapterStatusIcon,
  getChapterStatusColor,
  ChapterStatus,
  ChapterProgressData,
} from '../../lib/chapter-status';

describe('chapter-status 工具函数', () => {
  describe('getChapterStatus', () => {
    it('当进度数据为 null 时应返回 "not_started"', () => {
      const status = getChapterStatus(null);
      expect(status).toBe('not_started');
    });

    it('当进度数据为 undefined 时应返回 "not_started"', () => {
      const status = getChapterStatus(undefined);
      expect(status).toBe('not_started');
    });

    it('当 completed_at 有值时应返回 "completed"', () => {
      const progressData: ChapterProgressData = {
        status: 'in_progress',
        completed_at: '2024-01-01T00:00:00Z',
      };
      const status = getChapterStatus(progressData);
      expect(status).toBe('completed');
    });

    it('当状态为 "completed" 时应返回 "completed"', () => {
      const progressData: ChapterProgressData = {
        status: 'completed',
      };
      const status = getChapterStatus(progressData);
      expect(status).toBe('completed');
    });

    it('当 first_visit_at 有值但未完成时应返回 "in_progress"', () => {
      const progressData: ChapterProgressData = {
        status: 'not_started',
        first_visit_at: '2024-01-01T00:00:00Z',
      };
      const status = getChapterStatus(progressData);
      expect(status).toBe('in_progress');
    });

    it('当状态为 "in_progress" 时应返回 "in_progress"', () => {
      const progressData: ChapterProgressData = {
        status: 'in_progress',
        first_visit_at: '2024-01-01T00:00:00Z',
      };
      const status = getChapterStatus(progressData);
      expect(status).toBe('in_progress');
    });

    it('当状态为 "not_started" 且无访问记录时应返回 "not_started"', () => {
      const progressData: ChapterProgressData = {
        status: 'not_started',
      };
      const status = getChapterStatus(progressData);
      expect(status).toBe('not_started');
    });

    it('completed_at 优先级高于 status 字段', () => {
      const progressData: ChapterProgressData = {
        status: 'in_progress',
        completed_at: '2024-01-01T00:00:00Z',
      };
      const status = getChapterStatus(progressData);
      expect(status).toBe('completed');
    });
  });

  describe('getChapterStatusText', () => {
    it('应返回 "未开始" 对于 not_started 状态', () => {
      const text = getChapterStatusText('not_started');
      expect(text).toBe('未开始');
    });

    it('应返回 "学习中" 对于 in_progress 状态', () => {
      const text = getChapterStatusText('in_progress');
      expect(text).toBe('学习中');
    });

    it('应返回 "已完成" 对于 completed 状态', () => {
      const text = getChapterStatusText('completed');
      expect(text).toBe('已完成');
    });
  });

  describe('getChapterStatusIcon', () => {
    it('应返回 🔘 对于 not_started 状态', () => {
      const icon = getChapterStatusIcon('not_started');
      expect(icon).toBe('🔘');
    });

    it('应返回 📖 对于 in_progress 状态', () => {
      const icon = getChapterStatusIcon('in_progress');
      expect(icon).toBe('📖');
    });

    it('应返回 ✅ 对于 completed 状态', () => {
      const icon = getChapterStatusIcon('completed');
      expect(icon).toBe('✅');
    });
  });

  describe('getChapterStatusColor', () => {
    it('应返回 "default" 对于 not_started 状态', () => {
      const color = getChapterStatusColor('not_started');
      expect(color).toBe('default');
    });

    it('应返回 "processing" 对于 in_progress 状态', () => {
      const color = getChapterStatusColor('in_progress');
      expect(color).toBe('processing');
    });

    it('应返回 "success" 对于 completed 状态', () => {
      const color = getChapterStatusColor('completed');
      expect(color).toBe('success');
    });
  });
});
