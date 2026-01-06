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

// Helper to validate UUID
const isValidUUID = (str: string) => {
  const uuidRegex = /^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$/i;
  return uuidRegex.test(str);
};

// Helper to validate email
const isValidEmail = (str: string) => {
  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
  return emailRegex.test(str);
};

const directSchema = z.object({
  participant: z.string().min(1, "User ID or email is required").refine(
    (val) => isValidUUID(val) || isValidEmail(val),
    "Must be a valid UUID or email address"
  ),
});

const groupSchema = z.object({
  name: z.string().min(1, "Group name is required").max(50),
  participants: z.string().min(1, "At least one participant is required"),
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
      // Determine if input is UUID or email
      const isUUID = isValidUUID(data.participant);
      const result = await createConversation(
        isUUID
          ? { participant_id: data.participant }
          : { participant_email: data.participant }
      ).unwrap();

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
      // Parse comma-separated values (can be UUIDs or emails)
      const values = data.participants
        .split(",")
        .map((val) => val.trim())
        .filter((val) => val.length > 0);

      if (values.length === 0) {
        setError("At least one participant is required");
        return;
      }

      // Separate UUIDs and emails
      const uuids: string[] = [];
      const emails: string[] = [];
      for (const val of values) {
        if (isValidUUID(val)) {
          uuids.push(val);
        } else if (isValidEmail(val)) {
          emails.push(val);
        } else {
          setError(`Invalid value: "${val}" is not a valid UUID or email`);
          return;
        }
      }

      // Build request - prefer emails if all are emails, otherwise use IDs
      const requestBody: { name: string; participant_ids?: string[]; participant_emails?: string[] } = {
        name: data.name,
      };

      if (emails.length > 0 && uuids.length === 0) {
        requestBody.participant_emails = emails;
      } else if (uuids.length > 0 && emails.length === 0) {
        requestBody.participant_ids = uuids;
      } else {
        // Mixed - send IDs (backend will need both handled, or we can just use IDs)
        // For simplicity, if mixed, show error
        setError("Please use either all UUIDs or all emails, not mixed");
        return;
      }

      const result = await createConversation(requestBody).unwrap();

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
              <Label htmlFor="participant">User ID or Email</Label>
              <Input
                id="participant"
                placeholder="Enter UUID or email address"
                {...directForm.register("participant")}
              />
              {directForm.formState.errors.participant && (
                <p className="text-sm text-destructive">
                  {directForm.formState.errors.participant.message}
                </p>
              )}
              <p className="text-xs text-muted-foreground">
                Enter the UUID or email of the user you want to chat with.
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
              <Label htmlFor="participants">Participants</Label>
              <Input
                id="participants"
                placeholder="Enter UUIDs or emails separated by commas"
                {...groupForm.register("participants")}
              />
              {groupForm.formState.errors.participants && (
                <p className="text-sm text-destructive">
                  {groupForm.formState.errors.participants.message}
                </p>
              )}
              <p className="text-xs text-muted-foreground">
                Enter user UUIDs or email addresses separated by commas.
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


