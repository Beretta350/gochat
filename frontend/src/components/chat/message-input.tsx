"use client";

import { useState, useRef, useCallback } from "react";
import { Send, Smile, Paperclip } from "lucide-react";
import { m } from "framer-motion";
import { Button } from "@/components/ui/button";
import { cn } from "@/lib/utils";

interface MessageInputProps {
  onSend: (message: string) => void;
  disabled?: boolean;
  placeholder?: string;
}

export function MessageInput({
  onSend,
  disabled = false,
  placeholder = "Type a message...",
}: MessageInputProps) {
  const [message, setMessage] = useState("");
  const textareaRef = useRef<HTMLTextAreaElement>(null);

  const handleSubmit = useCallback(() => {
    const trimmedMessage = message.trim();
    if (trimmedMessage && !disabled) {
      onSend(trimmedMessage);
      setMessage("");
      
      // Reset textarea height
      if (textareaRef.current) {
        textareaRef.current.style.height = "auto";
      }
    }
  }, [message, disabled, onSend]);

  const handleKeyDown = (e: React.KeyboardEvent<HTMLTextAreaElement>) => {
    if (e.key === "Enter" && !e.shiftKey) {
      e.preventDefault();
      handleSubmit();
    }
  };

  const handleInput = (e: React.ChangeEvent<HTMLTextAreaElement>) => {
    setMessage(e.target.value);
    
    // Auto-resize textarea
    const textarea = textareaRef.current;
    if (textarea) {
      textarea.style.height = "auto";
      textarea.style.height = `${Math.min(textarea.scrollHeight, 120)}px`;
    }
  };

  const canSend = message.trim().length > 0 && !disabled;

  return (
    <div className="p-4 border-t border-border bg-background-secondary/50">
      <div className="flex items-center gap-3">
        {/* Attachment button */}
        <Button
          type="button"
          variant="ghost"
          size="icon"
          className="flex-shrink-0 h-10 w-10 text-muted-foreground hover:text-foreground"
          disabled={disabled}
        >
          <Paperclip className="w-5 h-5" />
        </Button>

        {/* Input container */}
        <div className="flex-1 flex items-center relative">
          <textarea
            ref={textareaRef}
            value={message}
            onChange={handleInput}
            onKeyDown={handleKeyDown}
            placeholder={placeholder}
            disabled={disabled}
            rows={1}
            className={cn(
              "w-full resize-none bg-muted rounded-2xl px-4 py-3 pr-12 text-base leading-6",
              "placeholder:text-muted-foreground",
              "focus:outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2 focus:ring-offset-background",
              "disabled:opacity-50 disabled:cursor-not-allowed",
              "scrollbar-hide"
            )}
            style={{ maxHeight: "120px", minHeight: "48px" }}
          />
          
          {/* Emoji button */}
          <Button
            type="button"
            variant="ghost"
            size="icon"
            className="absolute right-2 top-1/2 -translate-y-1/2 h-8 w-8 text-muted-foreground hover:text-foreground"
            disabled={disabled}
          >
            <Smile className="w-5 h-5" />
          </Button>
        </div>

        {/* Send button */}
        <m.div
          initial={false}
          animate={{
            scale: canSend ? 1 : 0.9,
            opacity: canSend ? 1 : 0.5,
          }}
          transition={{ duration: 0.15 }}
          className="flex-shrink-0"
        >
          <Button
            type="button"
            size="icon"
            onClick={handleSubmit}
            disabled={!canSend}
            className="rounded-full w-10 h-10"
          >
            <Send className="w-5 h-5" />
          </Button>
        </m.div>
      </div>
    </div>
  );
}

