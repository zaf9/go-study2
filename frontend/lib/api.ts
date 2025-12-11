import axios from "axios";
import { message } from "antd";
import { API_BASE_URL, REQUEST_TIMEOUT } from "./constants";
import { clearTokens, getAccessToken, refreshAccessToken } from "./auth";

const isBrowser = typeof window !== "undefined";

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
