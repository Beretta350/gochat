"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import { motion } from "framer-motion";
import { useAuth } from "@/hooks";

interface AuthGuardProps {
  children: React.ReactNode;
  requireAuth?: boolean;
}

export function AuthGuard({ children, requireAuth = true }: AuthGuardProps) {
  const router = useRouter();
  const { isAuthenticated, isLoading } = useAuth();
  const [showContent, setShowContent] = useState(!requireAuth);

  useEffect(() => {
    if (isLoading) return;

    if (requireAuth && !isAuthenticated) {
      router.push("/login");
    } else if (!requireAuth && isAuthenticated) {
      router.push("/chat");
    } else {
      setShowContent(true);
    }
  }, [isAuthenticated, isLoading, requireAuth, router]);

  // For public pages (requireAuth=false), show content immediately
  if (!requireAuth) {
    return <>{children}</>;
  }

  // For protected pages, show loading while checking auth
  if (isLoading || !showContent) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-background">
        <motion.div
          animate={{ rotate: 360 }}
          transition={{ duration: 1, repeat: Infinity, ease: "linear" }}
          className="w-8 h-8 border-2 border-primary border-t-transparent rounded-full"
        />
      </div>
    );
  }

  if (requireAuth && !isAuthenticated) {
    return null;
  }

  return <>{children}</>;
}
