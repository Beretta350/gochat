"use client";

import { useEffect, useCallback, useState, useMemo } from "react";
import { m, AnimatePresence } from "framer-motion";
import { MessageCircle, Loader2 } from "lucide-react";
import { Sidebar } from "@/components/chat/sidebar";
import { ChatHeader } from "@/components/chat/chat-header";
import { MessageList } from "@/components/chat/message-list";
import { MessageInput } from "@/components/chat/message-input";
import { AuthGuard } from "@/components/auth";
import { useAppDispatch, useAppSelector } from "@/store";
import {
  useGetConversationsQuery,
  useGetMessagesQuery,
} from "@/store/api/conversationsApi";
import {
  setConversations,
  setMessages,
  addMessage,
  setActiveConversation,
} from "@/store/slices/chatSlice";
import { useAuth, useWebSocket } from "@/hooks";
import { cn } from "@/lib/utils";

function ChatContent() {
  const dispatch = useAppDispatch();
  const { user } = useAuth();
  const { sendMessage, isConnected } = useWebSocket();
  const [mounted, setMounted] = useState(false);

  // Safe selectors with defaults
  const chatState = useAppSelector((state) => state.chat);
  
  const activeConversationId = chatState?.activeConversationId ?? null;
  const conversations = useMemo(() => chatState?.conversations ?? [], [chatState?.conversations]);
  const messagesMap = useMemo(() => chatState?.messages ?? {}, [chatState?.messages]);

  // Set mounted after first render
  useEffect(() => {
    setMounted(true);
  }, []);

  // Fetch conversations
  const { data: conversationsData } = useGetConversationsQuery(undefined, {
    skip: !user || !mounted,
  });

  // Fetch messages for active conversation
  const { data: messagesData } = useGetMessagesQuery(
    { conversationId: activeConversationId || "" },
    { skip: !activeConversationId || !mounted }
  );

  // Sync conversations to Redux
  useEffect(() => {
    if (conversationsData && mounted) {
      dispatch(setConversations(conversationsData));
    }
  }, [conversationsData, dispatch, mounted]);

  // Sync messages to Redux
  useEffect(() => {
    if (messagesData && activeConversationId && mounted) {
      dispatch(
        setMessages({
          conversationId: activeConversationId,
          messages: messagesData.messages,
        })
      );
    }
  }, [messagesData, activeConversationId, dispatch, mounted]);

  // Get active conversation
  const activeConversation = useMemo(() => {
    if (!Array.isArray(conversations)) return undefined;
    return conversations.find((c) => c.id === activeConversationId);
  }, [conversations, activeConversationId]);

  // Get messages for active conversation
  const activeMessages = useMemo(() => {
    if (!activeConversationId || !messagesMap) return [];
    return messagesMap[activeConversationId] ?? [];
  }, [activeConversationId, messagesMap]);

  // Handle send message
  const handleSendMessage = useCallback(
    (content: string) => {
      if (!activeConversationId || !user) return;

      // Send via WebSocket
      sendMessage({
        conversation_id: activeConversationId,
        content,
      });

      // Optimistic update
      dispatch(
        addMessage({
          id: `temp-${Date.now()}`,
          conversation_id: activeConversationId,
          sender_id: user.id,
          sender_username: user.username,
          content,
          type: "text",
          sent_at: new Date().toISOString(),
        })
      );
    },
    [activeConversationId, user, sendMessage, dispatch]
  );

  // Handle back button (mobile)
  const handleBack = useCallback(() => {
    dispatch(setActiveConversation(null));
  }, [dispatch]);

  // Show loading until mounted
  if (!mounted) {
    return (
      <div className="h-screen flex items-center justify-center bg-background">
        <Loader2 className="w-8 h-8 animate-spin text-primary" />
      </div>
    );
  }

  const safeConversations = Array.isArray(conversations) ? conversations : [];

  return (
    <div className="h-screen flex bg-background overflow-hidden">
      {/* Sidebar */}
      <div
        className={cn(
          "w-full lg:w-80 xl:w-96 flex-shrink-0 transition-all duration-300",
          activeConversationId ? "hidden lg:block" : "block"
        )}
      >
        <Sidebar conversations={safeConversations} isConnected={isConnected} />
      </div>

      {/* Chat area */}
      <div
        className={cn(
          "flex-1 flex flex-col min-w-0",
          !activeConversationId ? "hidden lg:flex" : "flex"
        )}
      >
        <AnimatePresence mode="wait">
          {activeConversation ? (
            <m.div
              key={activeConversation.id}
              initial={{ opacity: 0 }}
              animate={{ opacity: 1 }}
              exit={{ opacity: 0 }}
              transition={{ duration: 0.2 }}
              className="flex-1 flex flex-col min-h-0"
            >
              {/* Header */}
              <ChatHeader
                conversation={activeConversation}
                currentUserId={user?.id || ""}
                onBack={handleBack}
              />

              {/* Messages */}
              <MessageList
                messages={activeMessages}
                currentUserId={user?.id || ""}
              />

              {/* Input */}
              <MessageInput
                onSend={handleSendMessage}
                disabled={!isConnected}
                placeholder={
                  isConnected ? "Type a message..." : "Connecting..."
                }
              />
            </m.div>
          ) : (
            <m.div
              key="empty"
              initial={{ opacity: 0 }}
              animate={{ opacity: 1 }}
              exit={{ opacity: 0 }}
              className="flex-1 flex flex-col items-center justify-center bg-background-secondary/30"
            >
              <div className="text-center">
                <div className="w-20 h-20 rounded-full bg-muted flex items-center justify-center mx-auto mb-4">
                  <MessageCircle className="w-10 h-10 text-muted-foreground" />
                </div>
                <h2 className="text-xl font-semibold mb-2">
                  Welcome to GoChat
                </h2>
                <p className="text-muted-foreground max-w-sm">
                  Select a conversation from the sidebar or start a new one to
                  begin chatting.
                </p>
              </div>
            </m.div>
          )}
        </AnimatePresence>
      </div>
    </div>
  );
}

export default function ChatPage() {
  return (
    <AuthGuard requireAuth={true}>
      <ChatContent />
    </AuthGuard>
  );
}
