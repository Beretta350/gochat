"use client";

import { LazyMotion, domAnimation } from "framer-motion";
import { ReduxProvider } from "./redux-provider";
import { TooltipProvider } from "@/components/ui/tooltip";
import { ToastProvider, ToastViewport } from "@/components/ui/toast";

interface ProvidersProps {
  children: React.ReactNode;
}

export function Providers({ children }: ProvidersProps) {
  return (
    <ReduxProvider>
      <LazyMotion features={domAnimation} strict>
        <ToastProvider>
          <TooltipProvider delayDuration={0}>
            {children}
            <ToastViewport />
          </TooltipProvider>
        </ToastProvider>
      </LazyMotion>
    </ReduxProvider>
  );
}

