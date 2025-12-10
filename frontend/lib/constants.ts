export const API_BASE_URL =
  process.env.NEXT_PUBLIC_API_URL?.replace(/\/+$/, "") || "/api/v1";

export const ACCESS_TOKEN_KEY = "go-study2.access_token";
export const REMEMBER_ME_KEY = "go-study2.remember_me";
export const REQUEST_TIMEOUT = 10000;

export const API_PATHS = {
  login: "/auth/login",
  register: "/auth/register",
  refresh: "/auth/refresh",
  profile: "/auth/profile",
};

