"use client";

import { useCallback, useEffect, useRef, useState } from "react";
import { useAppDispatch, useAppSelector } from "@/store";
import { addMessage, setConnected } from "@/store/slices/chatSlice";
import type { WebSocketMessage, SendMessageRequest } from "@/types";

const WS_URL = process.env.NEXT_PUBLIC_WS_URL || "ws://localhost:8080";

export function useWebSocket() {
  const dispatch = useAppDispatch();
  const { accessToken } = useAppSelector((state) => state.auth);
  const { isConnected } = useAppSelector((state) => state.chat);
  const wsRef = useRef<WebSocket | null>(null);
  const reconnectTimeoutRef = useRef<NodeJS.Timeout | null>(null);
  const [reconnectAttempts, setReconnectAttempts] = useState(0);

  const connect = useCallback(() => {
    if (!accessToken || wsRef.current?.readyState === WebSocket.OPEN) {
      return;
    }

    try {
      const ws = new WebSocket(`${WS_URL}/ws?token=${accessToken}`);

      ws.onopen = () => {
        console.log("WebSocket connected");
        dispatch(setConnected(true));
        setReconnectAttempts(0);
      };

      ws.onmessage = (event) => {
        try {
          const data = JSON.parse(event.data);
          
          // Check if it's an error message
          if (data.error) {
            console.error("WebSocket error:", data.message);
            return;
          }

          // It's a chat message
          const message: WebSocketMessage = data;
          dispatch(
            addMessage({
              id: message.id,
              conversation_id: message.conversation_id,
              sender_id: message.sender_id,
              sender_username: message.sender_username,
              content: message.content,
              type: message.type as "text" | "image" | "file" | "audio",
              sent_at: new Date(message.sent_at).toISOString(),
            })
          );
        } catch (err) {
          console.error("Failed to parse WebSocket message:", err);
        }
      };

      ws.onclose = (event) => {
        console.log("WebSocket closed:", event.code, event.reason);
        dispatch(setConnected(false));
        wsRef.current = null;

        // Reconnect logic with exponential backoff
        if (accessToken && reconnectAttempts < 5) {
          const delay = Math.min(1000 * Math.pow(2, reconnectAttempts), 30000);
          console.log(`Reconnecting in ${delay}ms...`);
          reconnectTimeoutRef.current = setTimeout(() => {
            setReconnectAttempts((prev) => prev + 1);
            connect();
          }, delay);
        }
      };

      ws.onerror = (error) => {
        console.error("WebSocket error:", error);
      };

      wsRef.current = ws;
    } catch (err) {
      console.error("Failed to create WebSocket connection:", err);
    }
  }, [accessToken, dispatch, reconnectAttempts]);

  const disconnect = useCallback(() => {
    if (reconnectTimeoutRef.current) {
      clearTimeout(reconnectTimeoutRef.current);
      reconnectTimeoutRef.current = null;
    }
    if (wsRef.current) {
      wsRef.current.close();
      wsRef.current = null;
    }
    dispatch(setConnected(false));
  }, [dispatch]);

  const sendMessage = useCallback((message: SendMessageRequest) => {
    if (wsRef.current?.readyState === WebSocket.OPEN) {
      wsRef.current.send(JSON.stringify(message));
    } else {
      console.error("WebSocket is not connected");
    }
  }, []);

  useEffect(() => {
    if (accessToken) {
      connect();
    }

    return () => {
      disconnect();
    };
  }, [accessToken, connect, disconnect]);

  return {
    isConnected,
    sendMessage,
    connect,
    disconnect,
  };
}

