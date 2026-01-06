"use client";

import { useMemo, useState } from "react";
import { m } from "framer-motion";
import { ChevronLeft, Phone, Video, MoreVertical, Users } from "lucide-react";
import { Avatar, AvatarFallback } from "@/components/ui/avatar";
import { Button } from "@/components/ui/button";
import { WipDialog } from "@/components/ui/wip-dialog";
import { cn, getInitials } from "@/lib/utils";
import { useAppSelector } from "@/store";
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
  const { onlineUsers } = useAppSelector((state) => state.chat);
  const [wipOpen, setWipOpen] = useState(false);
  const [wipFeature, setWipFeature] = useState("");

  const openWip = (feature: string) => {
    setWipFeature(feature);
    setWipOpen(true);
  };

  const { displayName, subtitle, isGroup, isOtherUserOnline } = useMemo(() => {
    if (conversation.type === "group") {
      return {
        displayName: conversation.name || "Group Chat",
        subtitle: `${conversation.participants.length} members`,
        isGroup: true,
        isOtherUserOnline: false,
      };
    }

    const otherParticipant = conversation.participants.find(
      (p) => p.id !== currentUserId
    );
    const isOnline = otherParticipant
      ? onlineUsers.includes(otherParticipant.id)
      : false;
    return {
      displayName: otherParticipant?.username || "Unknown User",
      subtitle: isOnline ? "Online" : "Offline",
      isGroup: false,
      isOtherUserOnline: isOnline,
    };
  }, [conversation, currentUserId, onlineUsers]);

  return (
    <m.div
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
              className="text-sm font-medium"
            >
              {isGroup ? (
                <Users className="w-4 h-4" />
              ) : (
                getInitials(displayName)
              )}
            </AvatarFallback>
          </Avatar>
          {!isGroup && isOtherUserOnline && (
            <span className="absolute bottom-0 right-0 w-3 h-3 bg-success rounded-full border-2 border-background" />
          )}
        </div>

        {/* Name and status */}
        <div>
          <h2 className="font-semibold text-sm">{displayName}</h2>
          <p
            className={cn(
              "text-xs",
              !isGroup && isOtherUserOnline
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
          onClick={() => openWip("Voice Call")}
        >
          <Phone className="w-5 h-5" />
        </Button>
        <Button
          variant="ghost"
          size="icon"
          className="text-muted-foreground hover:text-foreground"
          onClick={() => openWip("Video Call")}
        >
          <Video className="w-5 h-5" />
        </Button>
        <Button
          variant="ghost"
          size="icon"
          className="text-muted-foreground hover:text-foreground"
          onClick={() => openWip("More Options")}
        >
          <MoreVertical className="w-5 h-5" />
        </Button>
      </div>

      <WipDialog open={wipOpen} onOpenChange={setWipOpen} feature={wipFeature} />
    </m.div>
  );
}

