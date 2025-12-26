import { redirect } from "next/navigation";

/**
 * 根页面
 * 未登录用户会被 AuthContext 重定向到登录页面
 * 已登录用户重定向到 Dashboard 首页
 */
export default function Home() {
  redirect("/dashboard");
}
