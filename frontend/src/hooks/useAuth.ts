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
import {
  useLoginMutation,
  useRegisterMutation,
  useLogoutMutation,
  useLazyGetMeQuery,
} from "@/store/api/authApi";
import type { LoginRequest, RegisterRequest } from "@/types";

export function useAuth() {
  const dispatch = useAppDispatch();
  const router = useRouter();
  const [mounted, setMounted] = useState(false);
  const [isValidating, setIsValidating] = useState(true);

  const authState = useAppSelector((state) => state.auth);
  const user = authState?.user ?? null;
  const isAuthenticated = authState?.isAuthenticated ?? false;

  const [loginMutation, { isLoading: isLoginLoading }] = useLoginMutation();
  const [registerMutation, { isLoading: isRegisterLoading }] =
    useRegisterMutation();
  const [logoutMutation] = useLogoutMutation();
  const [getMe] = useLazyGetMeQuery();

  // Validate session on mount by calling /me
  // If cookies are valid, backend will authenticate
  useEffect(() => {
    const validateSession = async () => {
      setMounted(true);

      if (typeof window === "undefined") {
        setIsValidating(false);
        dispatch(setLoading(false));
        return;
      }

      try {
        // Try to get current user (cookies are sent automatically)
        const result = await getMe().unwrap();

        // Session is valid - set user in store
        dispatch(
          setCredentials({
            user: result,
          })
        );
      } catch {
        // Session is invalid or no cookies - user not authenticated
        dispatch(logoutAction());
      } finally {
        setIsValidating(false);
        dispatch(setLoading(false));
      }
    };

    validateSession();
  }, [dispatch, getMe]);

  // isLoading is true until mounted and validated on client
  const isLoading = !mounted || isValidating;

  const login = useCallback(
    async (credentials: LoginRequest) => {
      try {
        const response = await loginMutation(credentials).unwrap();

        // Cookies are set automatically by the backend
        // Just update the store with user info
        dispatch(
          setCredentials({
            user: response.user,
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

        // Cookies are set automatically by the backend
        // Just update the store with user info
        dispatch(
          setCredentials({
            user: response.user,
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

  const logout = useCallback(async () => {
    try {
      // Call backend to clear cookies
      await logoutMutation().unwrap();
    } catch {
      // Even if API call fails, clear local state
    }
    dispatch(logoutAction());
    router.push("/login");
  }, [dispatch, router, logoutMutation]);

  return {
    user,
    isAuthenticated,
    isLoading,
    login,
    register,
    logout,
    isLoginLoading,
    isRegisterLoading,
  };
}
