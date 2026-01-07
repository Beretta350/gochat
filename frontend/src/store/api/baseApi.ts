import {
  fetchBaseQuery,
  type BaseQueryFn,
  type FetchArgs,
  type FetchBaseQueryError,
} from "@reduxjs/toolkit/query/react";
import { logout } from "../slices/authSlice";

// API URL - use env var in production, relative in dev
const apiUrl = process.env.NEXT_PUBLIC_API_URL || "";
const rawBaseQuery = fetchBaseQuery({
  baseUrl: `${apiUrl}/api/v1`,
  credentials: "include", // Envia cookies automaticamente
  prepareHeaders: (headers) => {
    headers.set("Content-Type", "application/json");
    return headers;
  },
});

// BaseQuery with automatic token refresh on 401
export const baseQuery: BaseQueryFn<
  string | FetchArgs,
  unknown,
  FetchBaseQueryError
> = async (args, api, extraOptions) => {
  let result = await rawBaseQuery(args, api, extraOptions);

  // If we get a 401, try to refresh the token
  if (result.error && result.error.status === 401) {
    // Try to refresh the token (cookie is sent automatically)
    const refreshResult = await rawBaseQuery(
      {
        url: "/auth/refresh",
        method: "POST",
      },
      api,
      extraOptions
    );

    if (refreshResult.data) {
      // Refresh succeeded - new cookies are set automatically by backend
      // Retry the original request
      result = await rawBaseQuery(args, api, extraOptions);
    } else {
      // Refresh failed - logout user
      api.dispatch(logout());
    }
  }

  return result;
};
