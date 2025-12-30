/**
 * 错误处理工具
 * 
 * 提供统一的错误处理和用户友好的错误提示
 * 
 * @module lib/error-handler
 */

/* eslint-disable @typescript-eslint/no-explicit-any */

import { toast } from 'react-hot-toast';

/**
 * API错误类型
 */
export interface ApiError {
  code: number;
  message: string;
  details?: any;
}

/**
 * 错误类型枚举
 */
export enum ErrorType {
  NETWORK = 'NETWORK',
  AUTH = 'AUTH',
  VALIDATION = 'VALIDATION',
  NOT_FOUND = 'NOT_FOUND',
  SERVER = 'SERVER',
  UNKNOWN = 'UNKNOWN',
}

/**
 * 错误消息映射
 */
const ERROR_MESSAGES: Record<ErrorType, string> = {
  [ErrorType.NETWORK]: '网络连接失败，请检查网络后重试',
  [ErrorType.AUTH]: '认证失败，请重新登录',
  [ErrorType.VALIDATION]: '输入数据验证失败',
  [ErrorType.NOT_FOUND]: '请求的资源不存在',
  [ErrorType.SERVER]: '服务器错误，请稍后重试',
  [ErrorType.UNKNOWN]: '未知错误，请稍后重试',
};

/**
 * 根据HTTP状态码判断错误类型
 */
export function getErrorType(status?: number): ErrorType {
  if (!status) return ErrorType.NETWORK;

  if (status === 401 || status === 403) return ErrorType.AUTH;
  if (status === 404) return ErrorType.NOT_FOUND;
  if (status >= 400 && status < 500) return ErrorType.VALIDATION;
  if (status >= 500) return ErrorType.SERVER;

  return ErrorType.UNKNOWN;
}

/**
 * 获取用户友好的错误消息
 */
export function getErrorMessage(error: any, defaultMessage?: string): string {
  // API错误
  if (error?.response?.data?.message) {
    return error.response.data.message;
  }

  // 网络错误
  if (error?.message === 'Network Error') {
    return ERROR_MESSAGES[ErrorType.NETWORK];
  }

  // HTTP状态错误
  if (error?.response?.status) {
    const errorType = getErrorType(error.response.status);
    return ERROR_MESSAGES[errorType];
  }

  // 自定义错误消息
  if (error?.message) {
    return error.message;
  }

  return defaultMessage || ERROR_MESSAGES[ErrorType.UNKNOWN];
}

/**
 * 显示错误提示toast
 */
export function showErrorToast(error: any, defaultMessage?: string): void {
  const message = getErrorMessage(error, defaultMessage);
  toast.error(message, {
    duration: 4000,
    position: 'top-center',
  });
}

/**
 * 显示成功提示toast
 */
export function showSuccessToast(message: string): void {
  toast.success(message, {
    duration: 3000,
    position: 'top-center',
  });
}

/**
 * 处理API调用错误
 * 
 * @param error - 错误对象
 * @param options - 处理选项
 * @returns 是否需要重试
 */
export function handleApiError(
  error: any,
  options: {
    silent?: boolean;
    customMessage?: string;
    onRetry?: () => void;
  } = {}
): boolean {
  const { silent = false, customMessage, onRetry } = options;

  if (!silent) {
    showErrorToast(error, customMessage);
  }

  // 网络错误或服务器错误建议重试
  const errorType = getErrorType(error?.response?.status);
  const shouldRetry = errorType === ErrorType.NETWORK || errorType === ErrorType.SERVER;

  if (shouldRetry && onRetry) {
    // 显示重试提示(简化版,避免JSX依赖)
    const retryMessage = '操作失败,点击重试按钮';
    toast.error(retryMessage, {
      duration: 5000,
    });
    // 注: 可在UI层实现更复杂的重试按钮组件
  }

  return shouldRetry;
}

/**
 * 包装异步函数，自动处理错误
 */
export function withErrorHandler<T extends (...args: any[]) => Promise<any>>(
  fn: T,
  options?: {
    silent?: boolean;
    customMessage?: string;
  }
): T {
  return (async (...args: any[]) => {
    try {
      return await fn(...args);
    } catch (error) {
      handleApiError(error, options);
      throw error;
    }
  }) as T;
}
