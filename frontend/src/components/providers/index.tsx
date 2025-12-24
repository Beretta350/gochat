"use client";

import { ReduxProvider } from "./redux-provider";
import { TooltipProvider } from "@/components/ui/tooltip";
import { ToastProvider, ToastViewport } from "@/components/ui/toast";

interface ProvidersProps {
  children: React.ReactNode;
}

export function Providers({ children }: ProvidersProps) {
  return (
    <ReduxProvider>
      <ToastProvider>
        <TooltipProvider delayDuration={0}>
          {children}
          <ToastViewport />
        </TooltipProvider>
      </ToastProvider>
    </ReduxProvider>
  );
}

