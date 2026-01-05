"use client";

import { useState, useEffect } from "react";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import { m, AnimatePresence } from "framer-motion";
import { Loader2, Mail, Lock, Eye, EyeOff, CheckCircle } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { AuthLayout } from "@/components/auth";
import { useAuth } from "@/hooks";

const loginSchema = z.object({
  email: z.string().email("Invalid email address"),
  password: z.string().min(1, "Password is required"),
});

type LoginForm = z.infer<typeof loginSchema>;

export default function LoginPage() {
  const router = useRouter();
  const { login, isLoginLoading, isAuthenticated, isLoading } = useAuth();
  const [showPassword, setShowPassword] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [checkingAuth, setCheckingAuth] = useState(true);
  const [authSuccess, setAuthSuccess] = useState(false);
  const [isSubmitting, setIsSubmitting] = useState(false);

  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<LoginForm>({
    resolver: zodResolver(loginSchema),
  });

  // Check authentication on page load
  useEffect(() => {
    if (!isLoading) {
      if (isAuthenticated) {
        // Already authenticated - show success and redirect
        setAuthSuccess(true);
        setTimeout(() => {
          router.push("/chat");
        }, 1000);
      } else {
        // Not authenticated - hide popup after 1.5s
        setTimeout(() => {
          setCheckingAuth(false);
        }, 1500);
      }
    }
  }, [isAuthenticated, isLoading, router]);

  const onSubmit = async (data: LoginForm) => {
    setError(null);
    setIsSubmitting(true);
    
    const result = await login(data);
    
    if (result.success) {
      setAuthSuccess(true);
      setTimeout(() => {
        router.push("/chat");
      }, 1000);
    } else {
      setIsSubmitting(false);
      setError(result.error || "Login failed");
    }
  };

  // Show checking auth popup
  if (checkingAuth || authSuccess) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-background">
        <m.div
          initial={{ scale: 0.9, opacity: 0 }}
          animate={{ scale: 1, opacity: 1 }}
          className="bg-card border border-border rounded-2xl p-8 shadow-xl flex flex-col items-center gap-4"
        >
          {authSuccess ? (
            <>
              <CheckCircle className="w-12 h-12 text-green-500" />
              <p className="text-lg font-medium">Welcome back!</p>
              <p className="text-sm text-muted-foreground">Redirecting to chat...</p>
            </>
          ) : (
            <>
              <Loader2 className="w-12 h-12 animate-spin text-primary" />
              <p className="text-lg font-medium">Checking session...</p>
              <p className="text-sm text-muted-foreground">Please wait</p>
            </>
          )}
        </m.div>
      </div>
    );
  }

  return (
    <AuthLayout>
      {/* Submitting Overlay */}
      <AnimatePresence>
        {isSubmitting && (
          <m.div
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
            className="fixed inset-0 z-50 flex items-center justify-center bg-background/80 backdrop-blur-sm"
          >
            <m.div
              initial={{ scale: 0.9, opacity: 0 }}
              animate={{ scale: 1, opacity: 1 }}
              exit={{ scale: 0.9, opacity: 0 }}
              className="bg-card border border-border rounded-2xl p-8 shadow-xl flex flex-col items-center gap-4"
            >
              <Loader2 className="w-12 h-12 animate-spin text-primary" />
              <p className="text-lg font-medium">Signing in...</p>
              <p className="text-sm text-muted-foreground">Please wait</p>
            </m.div>
          </m.div>
        )}
      </AnimatePresence>

      <div className="space-y-6">
        <div className="space-y-2 text-center">
          <h1 className="text-3xl font-bold">Welcome back</h1>
          <p className="text-muted-foreground">
            Enter your credentials to access your account
          </p>
        </div>

        <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
          {error && (
            <m.div
              initial={{ opacity: 0, y: -10 }}
              animate={{ opacity: 1, y: 0 }}
              className="p-3 rounded-lg bg-destructive/10 border border-destructive/20 text-destructive text-sm"
            >
              {error}
            </m.div>
          )}

          <div className="space-y-2">
            <Label htmlFor="email">Email</Label>
            <div className="relative">
              <Mail className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-muted-foreground" />
              <Input
                id="email"
                type="email"
                placeholder="you@example.com"
                className="pl-10"
                {...register("email")}
              />
            </div>
            {errors.email && (
              <p className="text-sm text-destructive">{errors.email.message}</p>
            )}
          </div>

          <div className="space-y-2">
            <Label htmlFor="password">Password</Label>
            <div className="relative">
              <Lock className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-muted-foreground" />
              <Input
                id="password"
                type={showPassword ? "text" : "password"}
                placeholder="••••••••"
                className="pl-10 pr-10"
                {...register("password")}
              />
              <button
                type="button"
                onClick={() => setShowPassword(!showPassword)}
                className="absolute right-3 top-1/2 -translate-y-1/2 text-muted-foreground hover:text-foreground transition-colors"
              >
                {showPassword ? (
                  <EyeOff className="w-4 h-4" />
                ) : (
                  <Eye className="w-4 h-4" />
                )}
              </button>
            </div>
            {errors.password && (
              <p className="text-sm text-destructive">{errors.password.message}</p>
            )}
          </div>

          <Button
            type="submit"
            className="w-full"
            size="lg"
            disabled={isLoginLoading}
          >
            {isLoginLoading ? (
              <>
                <Loader2 className="w-4 h-4 mr-2 animate-spin" />
                Signing in...
              </>
            ) : (
              "Sign in"
            )}
          </Button>
        </form>

        <div className="text-center text-sm">
          <span className="text-muted-foreground">Don&apos;t have an account? </span>
          <Link
            href="/register"
            className="text-primary hover:text-primary-hover transition-colors font-medium"
          >
            Sign up
          </Link>
        </div>
      </div>
    </AuthLayout>
  );
}
