import axios from "axios";
import { API_BASE_URL, REQUEST_TIMEOUT } from "./constants";
import type { AxiosResponse } from "axios";

export const authClient = axios.create({
  baseURL: API_BASE_URL,
  timeout: REQUEST_TIMEOUT,
  withCredentials: true,
});

// 请求拦截：标准化路径（去除尾部斜杠）
// 某些测试环境（mock axios）可能没有 interceptors 对象，做健壮性检查
if (authClient.interceptors && authClient.interceptors.request) {
  authClient.interceptors.request.use((config) => {
    if (config && typeof config.url === "string") {
      // 不移除根路径的第一个 /
      config.url = config.url.replace(/\/+$|(?<!^)\/$/, "");
    }
    return config;
  });
}

// 响应错误拦截：映射为用户友好的错误信息
if (authClient.interceptors && authClient.interceptors.response) {
  // Success handler: API 仓库约定的 code != 20000 视为错误
  authClient.interceptors.response.use(
    (resp: AxiosResponse) => {
      try {
        const payload = resp?.data as { code?: number; message?: string } | undefined;
        if (payload && typeof payload.code === "number" && payload.code !== 20000) {
          const msg = payload.message || "请求失败";
          // 在浏览器环境广播错误事件，页面组件可监听并展示友好 UI
          if (typeof window !== "undefined" && window.dispatchEvent) {
            window.dispatchEvent(
              new CustomEvent("app:error", { detail: { message: msg } }),
            );
          }
          return Promise.reject(new Error(msg));
        }
      } catch {
        // ignore and continue
      }
      return resp;
    },
    (error) => {
      // Axios 网络或 HTTP 错误
      const aerr = error as { response?: { status?: number; data?: { message?: string } }; message?: string };
      const status = aerr?.response?.status;
      let message = aerr?.message || "网络错误，请稍后重试";
      if (status === 401) {
        message = "用户名或密码错误";
      } else if (status === 400) {
        message = aerr.response?.data?.message || "请求参数无效";
      } else if (status === 404) {
        message = "服务未找到或不可用";
      } else if (status && status >= 500) {
        message = "服务器繁忙，请稍后再试";
      }
      if (typeof window !== "undefined" && window.dispatchEvent) {
        window.dispatchEvent(new CustomEvent("app:error", { detail: { message } }));
      }
      return Promise.reject(new Error(message));
    },
  );
}

export default authClient;
