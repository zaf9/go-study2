'use client';

import { useContext } from "react";
import { AuthContext } from "@/contexts/AuthContext";

export default function useAuth() {
  const ctx = useContext(AuthContext);
  if (!ctx) {
    throw new Error("AuthContext 未初始化");
  }
  return ctx;
}


