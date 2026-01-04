"use client";

import { useState } from "react";
import Image from "next/image";
import { motion } from "framer-motion";
import {
  MessageSquarePlus,
  Search,
  Settings,
  LogOut,
  Wifi,
  WifiOff,
} from "lucide-react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Avatar, AvatarFallback } from "@/components/ui/avatar";
import { Separator } from "@/components/ui/separator";
import {
  Tooltip,
  TooltipContent,
  TooltipTrigger,
} from "@/components/ui/tooltip";
import { ConversationList } from "./conversation-list";
import { NewConversationDialog } from "./new-conversation-dialog";
import { cn, getInitials, generateAvatarColor } from "@/lib/utils";
import { useAuth } from "@/hooks";
import type { Conversation } from "@/types";

interface SidebarProps {
  conversations: Conversation[];
  isConnected: boolean;
}

export function Sidebar({ conversations, isConnected }: SidebarProps) {
  const { user, logout } = useAuth();
  const [searchQuery, setSearchQuery] = useState("");
  const [showNewConversation, setShowNewConversation] = useState(false);

  const filteredConversations = conversations.filter((conv) => {
    if (!searchQuery) return true;
    const query = searchQuery.toLowerCase();

    if (conv.type === "group" && conv.name) {
      return conv.name.toLowerCase().includes(query);
    }

    return conv.participants.some((p) =>
      p.username.toLowerCase().includes(query)
    );
  });

  return (
    <>
      <div className="w-full h-full flex flex-col bg-background-secondary border-r border-border">
        {/* Header */}
        <div className="p-4 flex items-center justify-between">
          <div className="flex items-center gap-3">
            <Image
              src="/gochat.svg"
              alt="GoChat"
              width={32}
              height={32}
              className="w-8 h-8"
            />
            <span className="font-bold text-lg gradient-text">GoChat</span>
          </div>

          <div className="flex items-center gap-1">
            {/* Connection status */}
            <Tooltip>
              <TooltipTrigger asChild>
                <div
                  className={cn(
                    "p-2 rounded-lg",
                    isConnected ? "text-success" : "text-destructive"
                  )}
                >
                  {isConnected ? (
                    <Wifi className="w-4 h-4" />
                  ) : (
                    <WifiOff className="w-4 h-4" />
                  )}
                </div>
              </TooltipTrigger>
              <TooltipContent>
                {isConnected ? "Connected" : "Disconnected"}
              </TooltipContent>
            </Tooltip>

            {/* New conversation */}
            <Tooltip>
              <TooltipTrigger asChild>
                <Button
                  variant="ghost"
                  size="icon"
                  onClick={() => setShowNewConversation(true)}
                >
                  <MessageSquarePlus className="w-5 h-5" />
                </Button>
              </TooltipTrigger>
              <TooltipContent>New conversation</TooltipContent>
            </Tooltip>
          </div>
        </div>

        {/* Search */}
        <div className="px-4 pb-3">
          <div className="relative">
            <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-muted-foreground" />
            <Input
              placeholder="Search conversations..."
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              className="pl-9 bg-muted border-0"
            />
          </div>
        </div>

        <Separator />

        {/* Conversations list */}
        <ConversationList
          conversations={filteredConversations}
          currentUserId={user?.id || ""}
        />

        <Separator />

        {/* User section */}
        <div className="p-3">
          <div className="flex items-center justify-between gap-2">
            <div className="flex items-center gap-3 min-w-0 flex-1">
              <Avatar className="h-9 w-9 flex-shrink-0">
                <AvatarFallback
                  className={cn(
                    "text-xs",
                    generateAvatarColor(user?.username || "U")
                  )}
                >
                  {getInitials(user?.username || "User")}
                </AvatarFallback>
              </Avatar>
              <div className="min-w-0 flex-1">
                <p className="text-sm font-medium truncate leading-tight">
                  {user?.username}
                </p>
                <p className="text-xs text-muted-foreground truncate leading-tight">
                  {user?.email}
                </p>
              </div>
            </div>

            <div className="flex items-center gap-1 flex-shrink-0">
              <Tooltip>
                <TooltipTrigger asChild>
                  <Button variant="ghost" size="icon" className="h-8 w-8">
                    <Settings className="w-4 h-4" />
                  </Button>
                </TooltipTrigger>
                <TooltipContent>Settings</TooltipContent>
              </Tooltip>

              <Tooltip>
                <TooltipTrigger asChild>
                  <Button
                    variant="ghost"
                    size="icon"
                    className="h-8 w-8 text-destructive hover:text-destructive"
                    onClick={logout}
                  >
                    <LogOut className="w-4 h-4" />
                  </Button>
                </TooltipTrigger>
                <TooltipContent>Logout</TooltipContent>
              </Tooltip>
            </div>
          </div>
        </div>
      </div>

      {/* New conversation dialog */}
      <NewConversationDialog
        open={showNewConversation}
        onOpenChange={setShowNewConversation}
      />
    </>
  );
}

