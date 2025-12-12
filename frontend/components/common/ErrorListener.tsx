"use client";

import React, { useEffect, useState } from "react";
import { Alert } from "antd";

const AUTO_HIDE_MS = 5000;

export default function ErrorListener() {
  const [visible, setVisible] = useState(false);
  const [message, setMessage] = useState<string | undefined>(undefined);

  useEffect(() => {
    function handler(e: Event) {
      const detail = (e as CustomEvent)?.detail as any;
      const msg = detail?.message || "发生错误";
      setMessage(msg);
      setVisible(true);
      window.setTimeout(() => setVisible(false), AUTO_HIDE_MS);
    }
    if (typeof window !== "undefined" && window.addEventListener) {
      window.addEventListener("app:error", handler as EventListener);
    }
    return () => {
      if (typeof window !== "undefined" && window.removeEventListener) {
        window.removeEventListener("app:error", handler as EventListener);
      }
    };
  }, []);

  if (!visible || !message) return null;

  return (
    <div style={{ position: "fixed", top: 16, left: "50%", transform: "translateX(-50%)", zIndex: 9999, minWidth: 320 }}>
      <Alert type="error" showIcon message={message} closable onClose={() => setVisible(false)} />
    </div>
  );
}
