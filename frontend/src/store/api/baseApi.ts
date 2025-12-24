import { fetchBaseQuery } from "@reduxjs/toolkit/query/react";
import type { RootState } from "../store";

const API_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";

export const baseQuery = fetchBaseQuery({
  baseUrl: `${API_URL}/api/v1`,
  prepareHeaders: (headers, { getState }) => {
    const token = (getState() as RootState).auth.accessToken;
    if (token) {
      headers.set("Authorization", `Bearer ${token}`);
    }
    headers.set("Content-Type", "application/json");
    return headers;
  },
});

