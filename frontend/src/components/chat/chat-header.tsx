"use client";

import { useMemo } from "react";
import { motion } from "framer-motion";
import { ChevronLeft, Phone, Video, MoreVertical, Users } from "lucide-react";
import { Avatar, AvatarFallback } from "@/components/ui/avatar";
import { Button } from "@/components/ui/button";
import { cn, getInitials, generateAvatarColor } from "@/lib/utils";
import type { Conversation } from "@/types";

interface ChatHeaderProps {
  conversation: Conversation;
  currentUserId: string;
  isConnected: boolean;
  onBack?: () => void;
}

export function ChatHeader({
  conversation,
  currentUserId,
  isConnected,
  onBack,
}: ChatHeaderProps) {
  const { displayName, subtitle, isGroup } = useMemo(() => {
    if (conversation.type === "group") {
      return {
        displayName: conversation.name || "Group Chat",
        subtitle: `${conversation.participants.length} members`,
        isGroup: true,
      };
    }

    const otherParticipant = conversation.participants.find(
      (p) => p.id !== currentUserId
    );
    return {
      displayName: otherParticipant?.username || "Unknown User",
      subtitle: isConnected ? "Online" : "Offline",
      isGroup: false,
    };
  }, [conversation, currentUserId, isConnected]);

  return (
    <motion.div
      initial={{ opacity: 0, y: -10 }}
      animate={{ opacity: 1, y: 0 }}
      className="flex items-center justify-between px-4 py-3 border-b border-border bg-background-secondary/50 backdrop-blur-sm"
    >
      <div className="flex items-center gap-3">
        {/* Back button (mobile) */}
        {onBack && (
          <Button
            variant="ghost"
            size="icon"
            onClick={onBack}
            className="lg:hidden -ml-2"
          >
            <ChevronLeft className="w-5 h-5" />
          </Button>
        )}

        {/* Avatar */}
        <div className="relative">
          <Avatar className="h-10 w-10">
            <AvatarFallback
              className={cn(
                "text-sm font-medium",
                generateAvatarColor(displayName)
              )}
            >
              {isGroup ? (
                <Users className="w-4 h-4" />
              ) : (
                getInitials(displayName)
              )}
            </AvatarFallback>
          </Avatar>
          {!isGroup && isConnected && (
            <span className="absolute bottom-0 right-0 w-3 h-3 bg-success rounded-full border-2 border-background" />
          )}
        </div>

        {/* Name and status */}
        <div>
          <h2 className="font-semibold text-sm">{displayName}</h2>
          <p
            className={cn(
              "text-xs",
              !isGroup && isConnected
                ? "text-success"
                : "text-muted-foreground"
            )}
          >
            {subtitle}
          </p>
        </div>
      </div>

      {/* Actions */}
      <div className="flex items-center gap-1">
        <Button
          variant="ghost"
          size="icon"
          className="text-muted-foreground hover:text-foreground"
        >
          <Phone className="w-5 h-5" />
        </Button>
        <Button
          variant="ghost"
          size="icon"
          className="text-muted-foreground hover:text-foreground"
        >
          <Video className="w-5 h-5" />
        </Button>
        <Button
          variant="ghost"
          size="icon"
          className="text-muted-foreground hover:text-foreground"
        >
          <MoreVertical className="w-5 h-5" />
        </Button>
      </div>
    </motion.div>
  );
}

