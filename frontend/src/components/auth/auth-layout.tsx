"use client";

import Image from "next/image";
import Link from "next/link";
import { m } from "framer-motion";

interface AuthLayoutProps {
  children: React.ReactNode;
}

export function AuthLayout({ children }: AuthLayoutProps) {
  return (
    <div className="min-h-screen bg-gradient-to-br from-background via-background-secondary to-background flex">
      {/* Left side - Branding */}
      <div className="hidden lg:flex lg:w-1/2 relative overflow-hidden">
        {/* Background effects */}
        <div className="absolute inset-0">
          <div className="absolute top-1/4 left-1/4 w-96 h-96 bg-primary/20 rounded-full blur-3xl" />
          <div className="absolute bottom-1/4 right-1/4 w-96 h-96 bg-accent/10 rounded-full blur-3xl" />
        </div>

        <div className="relative z-10 flex flex-col items-center justify-center w-full p-12">
          <m.div
            initial={{ opacity: 0, scale: 0.8 }}
            animate={{ opacity: 1, scale: 1 }}
            transition={{ duration: 0.5 }}
            className="text-center"
          >
            <Image
              src="/gochat.svg"
              alt="GoChat"
              width={180}
              height={180}
              className="mx-auto mb-8"
            />
            <h1 className="text-4xl font-bold gradient-text mb-4">GoChat</h1>
            <p className="text-xl text-foreground-muted max-w-md">
              Connect with anyone, anywhere. Fast, secure, and real-time.
            </p>
          </m.div>

          {/* Floating chat bubbles decoration */}
          <div className="absolute inset-0 pointer-events-none">
            <m.div
              animate={{ y: [0, -10, 0] }}
              transition={{ duration: 3, repeat: Infinity, ease: "easeInOut" }}
              className="absolute top-1/4 left-1/4 w-16 h-16 bg-primary/20 rounded-2xl rounded-bl-sm"
            />
            <m.div
              animate={{ y: [0, 10, 0] }}
              transition={{ duration: 4, repeat: Infinity, ease: "easeInOut" }}
              className="absolute bottom-1/3 right-1/4 w-20 h-12 bg-accent/20 rounded-2xl rounded-br-sm"
            />
            <m.div
              animate={{ y: [0, -15, 0] }}
              transition={{ duration: 3.5, repeat: Infinity, ease: "easeInOut" }}
              className="absolute top-1/2 right-1/3 w-12 h-12 bg-secondary/20 rounded-2xl rounded-bl-sm"
            />
          </div>
        </div>
      </div>

      {/* Right side - Auth form */}
      <div className="w-full lg:w-1/2 flex items-center justify-center p-8">
        <div className="w-full max-w-md">
          {/* Mobile logo */}
          <div className="lg:hidden mb-8 text-center">
            <Link href="/" className="inline-flex items-center gap-3">
              <Image
                src="/gochat.svg"
                alt="GoChat"
                width={48}
                height={48}
                className="w-12 h-12"
              />
              <span className="text-2xl font-bold gradient-text">GoChat</span>
            </Link>
          </div>

          <m.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.4 }}
          >
            {children}
          </m.div>
        </div>
      </div>
    </div>
  );
}

