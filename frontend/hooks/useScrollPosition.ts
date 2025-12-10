'use client';

import { useEffect, useState } from "react";

export default function useScrollPosition() {
  const [scrollY, setScrollY] = useState(0);

  useEffect(() => {
    const handler = () => setScrollY(window.scrollY || 0);
    handler();
    window.addEventListener("scroll", handler, { passive: true });
    return () => {
      window.removeEventListener("scroll", handler);
    };
  }, []);

  return { scrollY };
}

