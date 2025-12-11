import { useEffect, useState } from "react";
import { message } from "antd";
import { useRouter } from "next/navigation";
import Loading from "@/components/common/Loading";
import useAuth from "@/hooks/useAuth";

interface AdminGuardProps {
  children: React.ReactNode;
}

export default function AdminGuard({ children }: AdminGuardProps) {
  const router = useRouter();
  const { user, loading } = useAuth();
  const [redirecting, setRedirecting] = useState(false);

  useEffect(() => {
    if (loading) {
      return;
    }
    if (!user) {
      setRedirecting(true);
      router.replace("/login");
      return;
    }
    if (!user.isAdmin) {
      setRedirecting(true);
      message.error("需要管理员权限");
      router.replace("/topics");
    }
  }, [loading, user, router]);

  if (loading || redirecting) {
    return <Loading />;
  }

  if (!user || !user.isAdmin) {
    return null;
  }

  return <>{children}</>;
}

