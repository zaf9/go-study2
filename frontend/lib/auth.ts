import axios from "axios";
import { ACCESS_TOKEN_KEY, API_BASE_URL, API_PATHS, REMEMBER_ME_KEY, REQUEST_TIMEOUT } from "./constants";

let accessTokenInMemory: string | null = null;
let rememberInMemory = false;

const isBrowser = typeof window !== "undefined";

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

export async function refreshAccessToken(): Promise<string | null> {
  try {
    const resp = await axios.post(
      `${API_BASE_URL}${API_PATHS.refresh}`,
      {},
      {
        withCredentials: true,
        timeout: REQUEST_TIMEOUT,
      }
    );
    const token = resp?.data?.data?.accessToken as string | undefined;
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

