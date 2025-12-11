"use client";

import type React from "react";
import {
  createContext,
  useCallback,
  useEffect,
  useMemo,
  useState,
} from "react";
import {
  clearTokens,
  changePassword as changePasswordApi,
  fetchProfile,
  getAccessToken,
  loginWithPassword,
  logoutAccount,
  registerAccount,
} from "@/lib/auth";
import { Profile } from "@/types/auth";

export interface AuthContextValue {
  user: Profile | null;
  loading: boolean;
  login: (
    username: string,
    password: string,
    remember: boolean,
  ) => Promise<Profile>;
  register: (
    username: string,
    password: string,
    remember: boolean,
  ) => Promise<Profile>;
  changePassword: (oldPassword: string, newPassword: string) => Promise<void>;
  logout: () => Promise<void>;
  refreshProfile: () => Promise<Profile | null>;
}

export const AuthContext = createContext<AuthContextValue | undefined>(
  undefined,
);

export function AuthProvider({ children }: React.PropsWithChildren) {
  const [user, setUser] = useState<Profile | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const token = getAccessToken();
    if (!token) {
      setLoading(false);
      return;
    }
    fetchProfile()
      .then((profile) => setUser(profile))
      .catch(() => {
        clearTokens();
        setUser(null);
      })
      .finally(() => setLoading(false));
  }, []);

  const login = useCallback(
    async (username: string, password: string, remember: boolean) => {
      setLoading(true);
      try {
        await loginWithPassword({ username, password, remember });
        const profile = await fetchProfile();
        setUser(profile);
        return profile;
      } finally {
        setLoading(false);
      }
    },
    [],
  );

  const register = useCallback(
    async (username: string, password: string, remember: boolean) => {
      setLoading(true);
      try {
        await registerAccount({ username, password, remember });
        const profile = await fetchProfile();
        setUser(profile);
        return profile;
      } finally {
        setLoading(false);
      }
    },
    [],
  );

  const logout = useCallback(async () => {
    try {
      await logoutAccount();
    } finally {
      clearTokens();
      setUser(null);
    }
  }, []);

  const changePassword = useCallback(
    async (oldPassword: string, newPassword: string) => {
      setLoading(true);
      try {
        await changePasswordApi({ oldPassword, newPassword });
        clearTokens();
        setUser(null);
      } finally {
        setLoading(false);
      }
    },
    [],
  );

  const refreshProfile = useCallback(async () => {
    const token = getAccessToken();
    if (!token) {
      setUser(null);
      return null;
    }
    try {
      const profile = await fetchProfile();
      setUser(profile);
      return profile;
    } catch {
      clearTokens();
      setUser(null);
      return null;
    }
  }, []);

  const value = useMemo<AuthContextValue>(
    () => ({
      user,
      loading,
      login,
      register,
      changePassword,
      logout,
      refreshProfile,
    }),
    [user, loading, login, register, changePassword, logout, refreshProfile],
  );

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}
