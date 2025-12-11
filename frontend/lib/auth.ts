import axios, { AxiosResponse } from "axios";
import {
  ACCESS_TOKEN_KEY,
  API_BASE_URL,
  API_PATHS,
  REMEMBER_ME_KEY,
  REQUEST_TIMEOUT,
} from "./constants";
import { ApiResponse } from "@/types/api";
import {
  AuthTokens,
  ChangePasswordRequest,
  LoginRequest,
  Profile,
  RegisterRequest,
} from "@/types/auth";

let accessTokenInMemory: string | null = null;
let rememberInMemory = false;

const isBrowser = typeof window !== "undefined";

const authClient = axios.create({
  baseURL: API_BASE_URL,
  timeout: REQUEST_TIMEOUT,
  withCredentials: true,
});

function readStorage(key: string): string | null {
  if (!isBrowser) return null;
  try {
    return window.localStorage.getItem(key);
  } catch {
    return null;
  }
}

function writeStorage(key: string, value: string | null) {
  if (!isBrowser) return;
  try {
    if (value === null) {
      window.localStorage.removeItem(key);
    } else {
      window.localStorage.setItem(key, value);
    }
  } catch {
    // 忽略存储异常（如无痕模式）
  }
}

export function setRememberMe(enabled: boolean) {
  rememberInMemory = enabled;
  writeStorage(REMEMBER_ME_KEY, enabled ? "1" : "0");
}

export function isRememberMe(): boolean {
  if (rememberInMemory) return true;
  const fromStorage = readStorage(REMEMBER_ME_KEY);
  return fromStorage === "1";
}

export function setAccessToken(token: string | null) {
  accessTokenInMemory = token;
  if (isRememberMe()) {
    writeStorage(ACCESS_TOKEN_KEY, token);
  } else {
    writeStorage(ACCESS_TOKEN_KEY, null);
  }
}

export function getAccessToken(): string | null {
  if (accessTokenInMemory) return accessTokenInMemory;
  if (!isRememberMe()) return null;
  const stored = readStorage(ACCESS_TOKEN_KEY);
  if (stored) {
    accessTokenInMemory = stored;
    return stored;
  }
  return null;
}

export function clearTokens() {
  accessTokenInMemory = null;
  writeStorage(ACCESS_TOKEN_KEY, null);
}

function buildAuthHeaders() {
  const token = getAccessToken();
  if (!token) {
    return {};
  }
  return {
    Authorization: `Bearer ${token}`,
  };
}

function unwrapResponse<T>(resp: AxiosResponse<ApiResponse<T>>): T {
  const payload = resp.data;
  if (!payload) {
    throw new Error("服务返回为空");
  }
  if (typeof payload.code === "number" && payload.code !== 20000) {
    throw new Error(payload.message || "请求失败");
  }
  return payload.data;
}

export async function loginWithPassword(
  request: LoginRequest,
): Promise<AuthTokens> {
  const remember = !!request.remember;
  setRememberMe(remember);
  const resp = await authClient.post<ApiResponse<AuthTokens>>(API_PATHS.login, {
    username: request.username,
    password: request.password,
    remember,
  });
  const tokens = unwrapResponse<AuthTokens>(resp);
  if (tokens?.accessToken) {
    setAccessToken(tokens.accessToken);
  }
  return tokens;
}

export async function registerAccount(
  request: RegisterRequest & { remember?: boolean },
): Promise<AuthTokens> {
  const resp = await authClient.post<ApiResponse<AuthTokens>>(
    API_PATHS.register,
    {
      username: request.username,
      password: request.password,
      remember: !!request.remember,
    },
    { headers: buildAuthHeaders() },
  );
  return unwrapResponse<AuthTokens>(resp);
}

export async function fetchProfile(): Promise<Profile> {
  const resp = await authClient.get<ApiResponse<Profile>>(API_PATHS.profile, {
    headers: buildAuthHeaders(),
  });
  return unwrapResponse<Profile>(resp);
}

export async function logoutAccount(): Promise<void> {
  await authClient.post<ApiResponse<null>>(
    API_PATHS.logout,
    {},
    { headers: buildAuthHeaders() },
  );
  clearTokens();
}

export async function refreshAccessToken(): Promise<string | null> {
  try {
    const resp = await authClient.post<ApiResponse<AuthTokens>>(
      API_PATHS.refresh,
      {},
    );
    const tokens = unwrapResponse<AuthTokens>(resp);
    const token = tokens?.accessToken;
    if (token) {
      setAccessToken(token);
      return token;
    }
    return null;
  } catch {
    clearTokens();
    return null;
  }
}

export async function changePassword(
  request: ChangePasswordRequest,
): Promise<void> {
  const resp = await authClient.post<ApiResponse<null>>(
    API_PATHS.changePassword,
    {
      oldPassword: request.oldPassword,
      newPassword: request.newPassword,
    },
    { headers: buildAuthHeaders() },
  );
  unwrapResponse<null>(resp);
  clearTokens();
}
