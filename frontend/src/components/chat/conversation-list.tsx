"use client";

import { useMemo } from "react";
import { motion } from "framer-motion";
import { MessageCircle, Users } from "lucide-react";
import { Avatar, AvatarFallback } from "@/components/ui/avatar";
import { ScrollArea } from "@/components/ui/scroll-area";
import { cn, formatDate, getInitials, generateAvatarColor } from "@/lib/utils";
import { useAppDispatch, useAppSelector } from "@/store";
import { setActiveConversation } from "@/store/slices/chatSlice";
import type { Conversation } from "@/types";

interface ConversationListProps {
  conversations: Conversation[];
  currentUserId: string;
}

export function ConversationList({
  conversations,
  currentUserId,
}: ConversationListProps) {
  const dispatch = useAppDispatch();
  const { activeConversationId, unreadCounts } = useAppSelector(
    (state) => state.chat
  );

  return (
    <ScrollArea className="flex-1">
      <div className="p-2 space-y-1">
        {conversations.length === 0 ? (
          <div className="flex flex-col items-center justify-center py-12 px-4 text-center">
            <div className="w-16 h-16 rounded-full bg-muted flex items-center justify-center mb-4">
              <MessageCircle className="w-8 h-8 text-muted-foreground" />
            </div>
            <p className="text-muted-foreground text-sm">
              No conversations yet
            </p>
            <p className="text-muted-foreground text-xs mt-1">
              Start a new conversation to begin chatting
            </p>
          </div>
        ) : (
          conversations.map((conversation, index) => (
            <ConversationItem
              key={conversation.id}
              conversation={conversation}
              currentUserId={currentUserId}
              isActive={activeConversationId === conversation.id}
              unreadCount={unreadCounts[conversation.id] || 0}
              index={index}
              onClick={() => dispatch(setActiveConversation(conversation.id))}
            />
          ))
        )}
      </div>
    </ScrollArea>
  );
}

interface ConversationItemProps {
  conversation: Conversation;
  currentUserId: string;
  isActive: boolean;
  unreadCount: number;
  index: number;
  onClick: () => void;
}

function ConversationItem({
  conversation,
  currentUserId,
  isActive,
  unreadCount,
  index,
  onClick,
}: ConversationItemProps) {
  const { displayName, isGroup } = useMemo(() => {
    if (conversation.type === "group") {
      return {
        displayName: conversation.name || "Group Chat",
        isGroup: true,
      };
    }

    // For direct conversations, show the other participant's name
    const otherParticipant = conversation.participants.find(
      (p) => p.id !== currentUserId
    );
    return {
      displayName: otherParticipant?.username || "Unknown User",
      isGroup: false,
    };
  }, [conversation, currentUserId]);

  const lastMessagePreview = useMemo(() => {
    if (!conversation.last_message) return "No messages yet";
    const content = conversation.last_message.content;
    return content.length > 40 ? content.slice(0, 40) + "..." : content;
  }, [conversation.last_message]);

  const lastMessageTime = useMemo(() => {
    if (!conversation.last_message) return "";
    const sentAt = conversation.last_message.sent_at;
    if (typeof sentAt === "number") {
      return formatDate(sentAt);
    }
    return formatDate(sentAt);
  }, [conversation.last_message]);

  return (
    <motion.button
      initial={{ opacity: 0, x: -20 }}
      animate={{ opacity: 1, x: 0 }}
      transition={{ duration: 0.2, delay: index * 0.05 }}
      onClick={onClick}
      className={cn(
        "w-full flex items-center gap-3 p-3 rounded-xl transition-all duration-200",
        "hover:bg-muted/50",
        isActive && "bg-muted"
      )}
    >
      <div className="relative">
        <Avatar className="h-12 w-12">
          <AvatarFallback
            className={cn(
              "text-sm font-medium",
              generateAvatarColor(displayName)
            )}
          >
            {isGroup ? (
              <Users className="w-5 h-5" />
            ) : (
              getInitials(displayName)
            )}
          </AvatarFallback>
        </Avatar>
        {!isGroup && (
          <span className="absolute bottom-0 right-0 w-3 h-3 bg-success rounded-full border-2 border-background" />
        )}
      </div>

      <div className="flex-1 min-w-0 text-left">
        <div className="flex items-center justify-between gap-2">
          <span className="font-medium text-sm truncate">{displayName}</span>
          {lastMessageTime && (
            <span className="text-xs text-muted-foreground flex-shrink-0">
              {lastMessageTime}
            </span>
          )}
        </div>
        <div className="flex items-center justify-between gap-2 mt-0.5">
          <p className="text-xs text-muted-foreground truncate">
            {lastMessagePreview}
          </p>
          {unreadCount > 0 && (
            <span className="flex-shrink-0 w-5 h-5 rounded-full bg-primary text-primary-foreground text-xs flex items-center justify-center font-medium">
              {unreadCount > 9 ? "9+" : unreadCount}
            </span>
          )}
        </div>
      </div>
    </motion.button>
  );
}

