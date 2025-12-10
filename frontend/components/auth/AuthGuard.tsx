'use client';

import type React from "react";
import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import useAuth from "@/hooks/useAuth";
import Loading from "@/components/common/Loading";

interface AuthGuardProps {
  children: React.ReactNode;
}

export default function AuthGuard({ children }: AuthGuardProps) {
  const router = useRouter();
  const { user, loading } = useAuth();
  const [redirecting, setRedirecting] = useState(false);

  useEffect(() => {
    if (!loading && !user) {
      setRedirecting(true);
      router.replace("/login");
    }
  }, [loading, user, router]);

  if (loading || (!user && redirecting)) {
    return <Loading />;
  }

  if (!user) {
    return null;
  }

  return <>{children}</>;
}


