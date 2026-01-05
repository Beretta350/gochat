"use client";

import { useEffect, useRef } from "react";
import { motion } from "framer-motion";
import { Avatar, AvatarFallback } from "@/components/ui/avatar";
import { ScrollArea } from "@/components/ui/scroll-area";
import { cn, getInitials } from "@/lib/utils";
import type { Message } from "@/types";

interface MessageListProps {
  messages: Message[];
  currentUserId: string;
}

export function MessageList({ messages, currentUserId }: MessageListProps) {
  const scrollRef = useRef<HTMLDivElement>(null);
  const bottomRef = useRef<HTMLDivElement>(null);

  // Auto-scroll to bottom on new messages
  useEffect(() => {
    bottomRef.current?.scrollIntoView({ behavior: "smooth" });
  }, [messages]);

  const formatMessageTime = (sentAt: string | number) => {
    const date = new Date(typeof sentAt === "number" ? sentAt : sentAt);
    return date.toLocaleTimeString("en-US", {
      hour: "2-digit",
      minute: "2-digit",
    });
  };

  // Group messages by date (messages come pre-sorted from backend)
  const groupedMessages = messages.reduce((groups, message) => {
    const date = new Date(
      typeof message.sent_at === "number" ? message.sent_at : message.sent_at
    ).toLocaleDateString("en-US", {
      weekday: "long",
      year: "numeric",
      month: "long",
      day: "numeric",
    });

    if (!groups[date]) {
      groups[date] = [];
    }
    groups[date].push(message);
    return groups;
  }, {} as Record<string, Message[]>);

  if (messages.length === 0) {
    return (
      <div className="flex-1 flex items-center justify-center">
        <div className="text-center">
          <p className="text-muted-foreground">No messages yet</p>
          <p className="text-muted-foreground text-sm mt-1">
            Send a message to start the conversation
          </p>
        </div>
      </div>
    );
  }

  return (
    <ScrollArea className="flex-1 px-4" ref={scrollRef}>
      <div className="py-4 space-y-6">
        {Object.entries(groupedMessages).map(([date, dateMessages]) => (
          <div key={date}>
            {/* Date separator */}
            <div className="flex items-center justify-center my-4">
              <span className="px-3 py-1 text-xs text-muted-foreground bg-muted rounded-full">
                {date}
              </span>
            </div>

            {/* Messages for this date */}
            <div className="space-y-3">
              {dateMessages.map((message, index) => {
                const isSent = message.sender_id === currentUserId;
                const showAvatar =
                  !isSent &&
                  (index === 0 ||
                    dateMessages[index - 1]?.sender_id !== message.sender_id);

                return (
                  <MessageBubble
                    key={message.id}
                    message={message}
                    isSent={isSent}
                    showAvatar={showAvatar}
                    time={formatMessageTime(message.sent_at)}
                  />
                );
              })}
            </div>
          </div>
        ))}
        <div ref={bottomRef} />
      </div>
    </ScrollArea>
  );
}

interface MessageBubbleProps {
  message: Message;
  isSent: boolean;
  showAvatar: boolean;
  time: string;
}

function MessageBubble({ message, isSent, showAvatar, time }: MessageBubbleProps) {
  return (
    <motion.div
      initial={{ opacity: 0, y: 10, scale: 0.95 }}
      animate={{ opacity: 1, y: 0, scale: 1 }}
      transition={{ duration: 0.2 }}
      className={cn("flex items-end gap-2", isSent && "flex-row-reverse")}
    >
      {/* Avatar for received messages */}
      {!isSent && (
        <div className="w-8 flex-shrink-0">
          {showAvatar && (
            <Avatar className="h-8 w-8">
              <AvatarFallback
                className="text-xs"
              >
                {getInitials(message.sender_username || "U")}
              </AvatarFallback>
            </Avatar>
          )}
        </div>
      )}

      <div
        className={cn(
          "flex flex-col gap-1 max-w-[70%]",
          isSent && "items-end"
        )}
      >
        {/* Sender name for received messages */}
        {!isSent && showAvatar && (
          <span className="text-xs text-muted-foreground ml-1">
            {message.sender_username}
          </span>
        )}

        {/* Message bubble */}
        <div
          className={cn(
            "px-4 py-2 rounded-2xl break-words",
            isSent
              ? "bg-primary text-primary-foreground rounded-br-md"
              : "bg-card text-card-foreground rounded-bl-md"
          )}
        >
          <p className="text-sm whitespace-pre-wrap">{message.content}</p>
        </div>

        {/* Time */}
        <span
          className={cn(
            "text-[10px] text-muted-foreground",
            isSent ? "mr-1" : "ml-1"
          )}
        >
          {time}
        </span>
      </div>
    </motion.div>
  );
}

export function TypingIndicator({ username }: { username?: string }) {
  return (
    <motion.div
      initial={{ opacity: 0, y: 10 }}
      animate={{ opacity: 1, y: 0 }}
      exit={{ opacity: 0, y: 10 }}
      className="flex items-center gap-2 px-4 py-2"
    >
      <div className="flex items-center gap-1 px-4 py-2 bg-card rounded-2xl rounded-bl-md">
        <span className="typing-dot" style={{ animationDelay: "0s" }} />
        <span className="typing-dot" style={{ animationDelay: "0.2s" }} />
        <span className="typing-dot" style={{ animationDelay: "0.4s" }} />
      </div>
      {username && (
        <span className="text-xs text-muted-foreground">
          {username} is typing...
        </span>
      )}
    </motion.div>
  );
}

