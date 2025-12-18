import axios from "axios";
import { message } from "antd";
import { API_BASE_URL, REQUEST_TIMEOUT } from "./constants";
import { clearTokens, getAccessToken, refreshAccessToken } from "./auth";

const isBrowser = typeof window !== "undefined";

// 网络超时重试配置
const RETRY_CONFIG = {
  maxRetries: 3,          // 最多重试 3 次
  retryDelay: 1000,       // 初始延迟 1 秒
  maxRetryDelay: 10000,   // 最大延迟 10 秒
  retryableStatuses: [408, 429, 500, 502, 503, 504], // 可重试的 HTTP 状态码
};

// 计算重试延迟（指数退避）
function getRetryDelay(retryCount: number): number {
  const delay = RETRY_CONFIG.retryDelay * Math.pow(2, retryCount - 1);
  return Math.min(delay, RETRY_CONFIG.maxRetryDelay);
}

// 是否应该重试该错误
function shouldRetry(error: unknown, retryCount: number): boolean {
  if (retryCount >= RETRY_CONFIG.maxRetries) {
    return false;
  }

  if (typeof error !== 'object' || error === null) {
    return false;
  }

  // 超时错误
  const err = error as Record<string, unknown>;
  if (err.code === 'ECONNABORTED' || err.message === '网络超时') {
    return true;
  }

  // 可重试的 HTTP 状态码
  if (err?.response && typeof err.response === 'object' && err.response !== null) {
    const response = err.response as Record<string, unknown>;
    if (response.status && RETRY_CONFIG.retryableStatuses.includes(response.status as number)) {
      return true;
    }
  }

  return false;
}

const api = axios.create({
  baseURL: API_BASE_URL,
  timeout: REQUEST_TIMEOUT,
  withCredentials: true,
});

api.interceptors.request.use((config) => {
  const token = getAccessToken();
  if (token) {
    config.headers = config.headers ?? {};
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

let refreshing: Promise<string | null> | null = null;

api.interceptors.response.use(
  (response) => {
    const { data } = response;
    if (
      data &&
      typeof data.code === "number" &&
      data.code !== 20000 &&
      data.code !== 0
    ) {
      if (isBrowser && data.message) {
        message.error(data.message);
      }
      return Promise.reject(new Error(data.message || "请求失败"));
    }
    return data?.data ?? data;
  },
  async (error) => {
    const config = error?.config;
    
    // 网络超时重试机制
    if (config && shouldRetry(error, (config.__retryCount || 0) + 1)) {
      config.__retryCount = (config.__retryCount || 0) + 1;
      const delay = getRetryDelay(config.__retryCount);
      
      if (isBrowser) {
        console.warn(`[API] 请求超时或失败，${delay}ms 后进行第 ${config.__retryCount} 次重试...`);
      }
      
      // 等待后重试
      await new Promise(resolve => setTimeout(resolve, delay));
      return api(config);
    }
    
    const status = error?.response?.status;
    if (status === 401) {
      if (!refreshing) {
        refreshing = refreshAccessToken().finally(() => {
          refreshing = null;
        });
      }
      const newToken = await refreshing;
      if (newToken) {
        const originalRequest = error.config;
        originalRequest.headers = originalRequest.headers ?? {};
        originalRequest.headers.Authorization = `Bearer ${newToken}`;
        return api(originalRequest);
      }
      if (isBrowser) {
        clearTokens();
        window.location.href = "/login";
      }
    }
    if (
      status === 403 &&
      error?.response?.data?.code === 40011 &&
      isBrowser
    ) {
      window.location.href = "/change-password";
    }

    const msg =
      error?.response?.data?.message || error.message || "网络错误，请重试";
    if (isBrowser) {
      message.error(msg);
    }
    return Promise.reject(error);
  },
);

export default api;
