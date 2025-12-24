"use client";

import { useCallback, useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import {
  useAppDispatch,
  useAppSelector,
  setCredentials,
  logout as logoutAction,
  setLoading,
} from "@/store";
import { useLoginMutation, useRegisterMutation } from "@/store/api/authApi";
import type { LoginRequest, RegisterRequest } from "@/types";

const TOKEN_KEY = "gochat_tokens";

interface StoredTokens {
  accessToken: string;
  refreshToken: string;
}

export function useAuth() {
  const dispatch = useAppDispatch();
  const router = useRouter();
  const [mounted, setMounted] = useState(false);
  
  const authState = useAppSelector((state) => state.auth);
  const user = authState?.user ?? null;
  const isAuthenticated = authState?.isAuthenticated ?? false;
  const accessToken = authState?.accessToken ?? null;

  const [loginMutation, { isLoading: isLoginLoading }] = useLoginMutation();
  const [registerMutation, { isLoading: isRegisterLoading }] =
    useRegisterMutation();

  // Load tokens from localStorage on mount
  useEffect(() => {
    setMounted(true);
    
    if (typeof window === "undefined") return;
    
    const stored = localStorage.getItem(TOKEN_KEY);
    if (stored) {
      try {
        const tokens: StoredTokens = JSON.parse(stored);
        dispatch(
          setCredentials({
            user: { id: "", email: "", username: "", is_active: true, created_at: "" },
            accessToken: tokens.accessToken,
            refreshToken: tokens.refreshToken,
          })
        );
      } catch {
        localStorage.removeItem(TOKEN_KEY);
      }
    }
    dispatch(setLoading(false));
  }, [dispatch]);
  
  // isLoading is true until mounted on client
  const isLoading = !mounted;

  const login = useCallback(
    async (credentials: LoginRequest) => {
      try {
        const response = await loginMutation(credentials).unwrap();
        
        // Store tokens
        localStorage.setItem(
          TOKEN_KEY,
          JSON.stringify({
            accessToken: response.tokens.access_token,
            refreshToken: response.tokens.refresh_token,
          })
        );

        dispatch(
          setCredentials({
            user: response.user,
            accessToken: response.tokens.access_token,
            refreshToken: response.tokens.refresh_token,
          })
        );

        router.push("/chat");
        return { success: true };
      } catch (error: unknown) {
        const err = error as { data?: { message?: string } };
        return {
          success: false,
          error: err?.data?.message || "Login failed",
        };
      }
    },
    [loginMutation, dispatch, router]
  );

  const register = useCallback(
    async (data: RegisterRequest) => {
      try {
        const response = await registerMutation(data).unwrap();

        // Store tokens
        localStorage.setItem(
          TOKEN_KEY,
          JSON.stringify({
            accessToken: response.tokens.access_token,
            refreshToken: response.tokens.refresh_token,
          })
        );

        dispatch(
          setCredentials({
            user: response.user,
            accessToken: response.tokens.access_token,
            refreshToken: response.tokens.refresh_token,
          })
        );

        router.push("/chat");
        return { success: true };
      } catch (error: unknown) {
        const err = error as { data?: { message?: string } };
        return {
          success: false,
          error: err?.data?.message || "Registration failed",
        };
      }
    },
    [registerMutation, dispatch, router]
  );

  const logout = useCallback(() => {
    localStorage.removeItem(TOKEN_KEY);
    dispatch(logoutAction());
    router.push("/login");
  }, [dispatch, router]);

  return {
    user,
    isAuthenticated,
    isLoading,
    accessToken,
    login,
    register,
    logout,
    isLoginLoading,
    isRegisterLoading,
  };
}

