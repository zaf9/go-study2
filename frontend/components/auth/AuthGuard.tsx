"use client";

import type React from "react";
import { useEffect, useState } from "react";
import { usePathname, useRouter } from "next/navigation";
import useAuth from "@/hooks/useAuth";
import Loading from "@/components/common/Loading";

interface AuthGuardProps {
  children: React.ReactNode;
}

export default function AuthGuard({ children }: AuthGuardProps) {
  const router = useRouter();
  const pathname = usePathname();
  const { user, loading } = useAuth();
  const [redirecting, setRedirecting] = useState(false);

  useEffect(() => {
    if (!loading && !user) {
      setRedirecting(true);
      router.replace("/login");
      return;
    }
    if (!loading && user?.mustChangePassword && pathname !== "/change-password") {
      setRedirecting(true);
      router.replace("/change-password");
    }
  }, [loading, user, router, pathname]);

  if (loading || redirecting) {
    return <Loading />;
  }

  if (!user) {
    return null;
  }

  if (user.mustChangePassword && pathname !== "/change-password") {
    return null;
  }

  return <>{children}</>;
}
