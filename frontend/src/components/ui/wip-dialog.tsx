"use client";

import { Construction } from "lucide-react";
import { m } from "framer-motion";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";

interface WipDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  feature?: string;
}

export function WipDialog({ open, onOpenChange, feature }: WipDialogProps) {
  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-md">
        <DialogHeader className="text-center sm:text-center">
          <m.div
            initial={{ scale: 0.8, opacity: 0 }}
            animate={{ scale: 1, opacity: 1 }}
            transition={{ delay: 0.1, type: "spring", stiffness: 200 }}
            className="mx-auto mb-4 flex h-16 w-16 items-center justify-center rounded-full bg-amber-500/10"
          >
            <Construction className="h-8 w-8 text-amber-500" />
          </m.div>
          <DialogTitle className="text-xl">Work in Progress</DialogTitle>
          <DialogDescription className="text-base pt-2">
            {feature ? (
              <>
                <span className="font-medium text-foreground">{feature}</span>{" "}
                is not implemented yet.
              </>
            ) : (
              "This feature is not implemented yet."
            )}
            <br />
            <span className="text-muted-foreground mt-2 block">
              Stay tuned for future updates! 🚧
            </span>
          </DialogDescription>
        </DialogHeader>
      </DialogContent>
    </Dialog>
  );
}

