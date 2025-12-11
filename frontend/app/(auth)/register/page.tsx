"use client";

import RegisterForm from "@/components/auth/RegisterForm";
import AdminGuard from "@/components/auth/AdminGuard";

export default function RegisterPage() {
  return (
    <AdminGuard>
      <div className="min-h-screen flex items-center justify-center bg-gray-50 px-4">
        <RegisterForm />
      </div>
    </AdminGuard>
  );
}
