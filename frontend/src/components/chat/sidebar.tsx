"use client";

import { useState } from "react";
import Image from "next/image";
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
import {
  Tooltip,
  TooltipContent,
  TooltipTrigger,
} from "@/components/ui/tooltip";
import { WipDialog } from "@/components/ui/wip-dialog";
import { ConversationList } from "./conversation-list";
import { NewConversationDialog } from "./new-conversation-dialog";
import { cn, getInitials } from "@/lib/utils";
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
  const [showWip, setShowWip] = useState(false);

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
              width={36}
              height={36}
              className="w-9 h-9"
            />
            <span className="font-bold text-xl gradient-text">GoChat</span>
          </div>

          <div className="flex items-center gap-2">
            {/* Connection status */}
            <Tooltip>
              <TooltipTrigger asChild>
                <div
                  className={cn(
                    "h-10 w-10 flex items-center justify-center rounded-lg",
                    isConnected ? "text-success" : "text-destructive"
                  )}
                >
                  {isConnected ? (
                    <Wifi className="w-5 h-5" />
                  ) : (
                    <WifiOff className="w-5 h-5" />
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
                  className="h-10 w-10"
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
        <div className="p-4">
          <div className="relative">
            <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-muted-foreground" />
            <Input
              placeholder="Search conversations..."
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              className="pl-10 h-11 bg-muted border-0 text-sm"
            />
          </div>
        </div>

        {/* Conversations list */}
        <ConversationList
          conversations={filteredConversations}
          currentUserId={user?.id || ""}
        />

        {/* User section */}
        <div className="p-4">
          <div className="flex items-center justify-between gap-3">
            <div className="flex items-center gap-3 min-w-0 flex-1">
              <Avatar className="h-11 w-11 flex-shrink-0">
                <AvatarFallback
                  className="text-sm font-medium"
                >
                  {getInitials(user?.username || "User")}
                </AvatarFallback>
              </Avatar>
              <div className="min-w-0 flex-1">
                <p className="text-sm font-semibold truncate">
                  {user?.username}
                </p>
                <p className="text-xs text-muted-foreground truncate mt-0.5">
                  {user?.email}
                </p>
              </div>
            </div>

            <div className="flex items-center gap-2 flex-shrink-0">
              <Tooltip>
                <TooltipTrigger asChild>
                  <Button
                    variant="ghost"
                    size="icon"
                    className="h-10 w-10"
                    onClick={() => setShowWip(true)}
                  >
                    <Settings className="w-5 h-5" />
                  </Button>
                </TooltipTrigger>
                <TooltipContent>Settings</TooltipContent>
              </Tooltip>

              <Tooltip>
                <TooltipTrigger asChild>
                  <Button
                    variant="ghost"
                    size="icon"
                    className="h-10 w-10 text-destructive hover:text-destructive"
                    onClick={logout}
                  >
                    <LogOut className="w-5 h-5" />
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

      {/* Work in Progress dialog */}
      <WipDialog
        open={showWip}
        onOpenChange={setShowWip}
        feature="Settings"
      />
    </>
  );
}

