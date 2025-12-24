"use client";

import { useState } from "react";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import { Loader2, User, Users } from "lucide-react";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
} from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { useCreateConversationMutation } from "@/store/api/conversationsApi";
import { useAppDispatch } from "@/store";
import { setActiveConversation, addConversation } from "@/store/slices/chatSlice";
import { cn } from "@/lib/utils";

const directSchema = z.object({
  participant_id: z.string().uuid("Invalid user ID"),
});

const groupSchema = z.object({
  name: z.string().min(1, "Group name is required").max(50),
  participant_ids: z.string().min(1, "At least one participant is required"),
});

type DirectForm = z.infer<typeof directSchema>;
type GroupForm = z.infer<typeof groupSchema>;

interface NewConversationDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
}

export function NewConversationDialog({
  open,
  onOpenChange,
}: NewConversationDialogProps) {
  const dispatch = useAppDispatch();
  const [mode, setMode] = useState<"direct" | "group">("direct");
  const [error, setError] = useState<string | null>(null);
  const [createConversation, { isLoading }] = useCreateConversationMutation();

  const directForm = useForm<DirectForm>({
    resolver: zodResolver(directSchema),
  });

  const groupForm = useForm<GroupForm>({
    resolver: zodResolver(groupSchema),
  });

  const handleDirectSubmit = async (data: DirectForm) => {
    setError(null);
    try {
      const result = await createConversation({
        participant_id: data.participant_id,
      }).unwrap();

      dispatch(addConversation(result));
      dispatch(setActiveConversation(result.id));
      onOpenChange(false);
      directForm.reset();
    } catch (err: unknown) {
      const error = err as { data?: { message?: string } };
      setError(error?.data?.message || "Failed to create conversation");
    }
  };

  const handleGroupSubmit = async (data: GroupForm) => {
    setError(null);
    try {
      // Parse comma-separated IDs
      const ids = data.participant_ids
        .split(",")
        .map((id) => id.trim())
        .filter((id) => id.length > 0);

      if (ids.length === 0) {
        setError("At least one participant ID is required");
        return;
      }

      const result = await createConversation({
        name: data.name,
        participant_ids: ids,
      }).unwrap();

      dispatch(addConversation(result));
      dispatch(setActiveConversation(result.id));
      onOpenChange(false);
      groupForm.reset();
    } catch (err: unknown) {
      const error = err as { data?: { message?: string } };
      setError(error?.data?.message || "Failed to create group");
    }
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-md">
        <DialogHeader>
          <DialogTitle>New Conversation</DialogTitle>
          <DialogDescription>
            Start a direct conversation or create a group chat.
          </DialogDescription>
        </DialogHeader>

        {/* Mode selector */}
        <div className="flex gap-2 p-1 bg-muted rounded-lg">
          <button
            onClick={() => {
              setMode("direct");
              setError(null);
            }}
            className={cn(
              "flex-1 flex items-center justify-center gap-2 py-2 px-4 rounded-md text-sm font-medium transition-colors",
              mode === "direct"
                ? "bg-background text-foreground shadow-sm"
                : "text-muted-foreground hover:text-foreground"
            )}
          >
            <User className="w-4 h-4" />
            Direct
          </button>
          <button
            onClick={() => {
              setMode("group");
              setError(null);
            }}
            className={cn(
              "flex-1 flex items-center justify-center gap-2 py-2 px-4 rounded-md text-sm font-medium transition-colors",
              mode === "group"
                ? "bg-background text-foreground shadow-sm"
                : "text-muted-foreground hover:text-foreground"
            )}
          >
            <Users className="w-4 h-4" />
            Group
          </button>
        </div>

        {error && (
          <div className="p-3 rounded-lg bg-destructive/10 border border-destructive/20 text-destructive text-sm">
            {error}
          </div>
        )}

        {/* Direct conversation form */}
        {mode === "direct" && (
          <form
            onSubmit={directForm.handleSubmit(handleDirectSubmit)}
            className="space-y-4"
          >
            <div className="space-y-2">
              <Label htmlFor="participant_id">User ID</Label>
              <Input
                id="participant_id"
                placeholder="Enter user UUID"
                {...directForm.register("participant_id")}
              />
              {directForm.formState.errors.participant_id && (
                <p className="text-sm text-destructive">
                  {directForm.formState.errors.participant_id.message}
                </p>
              )}
              <p className="text-xs text-muted-foreground">
                Enter the UUID of the user you want to chat with.
              </p>
            </div>

            <Button type="submit" className="w-full" disabled={isLoading}>
              {isLoading ? (
                <>
                  <Loader2 className="w-4 h-4 mr-2 animate-spin" />
                  Creating...
                </>
              ) : (
                "Start Conversation"
              )}
            </Button>
          </form>
        )}

        {/* Group conversation form */}
        {mode === "group" && (
          <form
            onSubmit={groupForm.handleSubmit(handleGroupSubmit)}
            className="space-y-4"
          >
            <div className="space-y-2">
              <Label htmlFor="name">Group Name</Label>
              <Input
                id="name"
                placeholder="Enter group name"
                {...groupForm.register("name")}
              />
              {groupForm.formState.errors.name && (
                <p className="text-sm text-destructive">
                  {groupForm.formState.errors.name.message}
                </p>
              )}
            </div>

            <div className="space-y-2">
              <Label htmlFor="participant_ids">Participant IDs</Label>
              <Input
                id="participant_ids"
                placeholder="Enter UUIDs separated by commas"
                {...groupForm.register("participant_ids")}
              />
              {groupForm.formState.errors.participant_ids && (
                <p className="text-sm text-destructive">
                  {groupForm.formState.errors.participant_ids.message}
                </p>
              )}
              <p className="text-xs text-muted-foreground">
                Enter user UUIDs separated by commas.
              </p>
            </div>

            <Button type="submit" className="w-full" disabled={isLoading}>
              {isLoading ? (
                <>
                  <Loader2 className="w-4 h-4 mr-2 animate-spin" />
                  Creating...
                </>
              ) : (
                "Create Group"
              )}
            </Button>
          </form>
        )}
      </DialogContent>
    </Dialog>
  );
}

