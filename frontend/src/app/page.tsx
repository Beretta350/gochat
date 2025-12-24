"use client";

import { motion } from "framer-motion";
import Image from "next/image";
import Link from "next/link";
import { Button } from "@/components/ui/button";
import { MessageCircle, Shield, Zap, Users } from "lucide-react";

export default function HomePage() {
  return (
    <div className="min-h-screen bg-gradient-to-br from-background via-background-secondary to-background">
      {/* Hero Section */}
      <div className="relative overflow-hidden">
        {/* Background effects */}
        <div className="absolute inset-0 overflow-hidden">
          <div className="absolute -top-40 -right-40 w-80 h-80 bg-primary/20 rounded-full blur-3xl" />
          <div className="absolute top-1/2 -left-40 w-80 h-80 bg-accent/10 rounded-full blur-3xl" />
        </div>

        <div className="relative max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 pt-20 pb-24">
          {/* Nav */}
          <nav className="flex items-center justify-between mb-16">
            <div className="flex items-center gap-3">
              <Image
                src="/gochat.svg"
                alt="GoChat"
                width={48}
                height={48}
                className="w-12 h-12"
              />
              <span className="text-2xl font-bold gradient-text">GoChat</span>
            </div>
            <div className="flex items-center gap-4">
              <Link href="/login">
                <Button variant="ghost">Login</Button>
              </Link>
              <Link href="/register">
                <Button>Get Started</Button>
              </Link>
            </div>
          </nav>

          {/* Hero Content */}
          <div className="text-center max-w-4xl mx-auto">
            <motion.h1
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ duration: 0.5 }}
              className="text-5xl md:text-7xl font-bold tracking-tight mb-6"
            >
              Chat in{" "}
              <span className="gradient-text">Real-Time</span>
            </motion.h1>
            <motion.p
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ duration: 0.5, delay: 0.1 }}
              className="text-xl text-foreground-muted mb-10 max-w-2xl mx-auto"
            >
              Connect instantly with friends and colleagues. Fast, secure, and
              beautifully designed for modern communication.
            </motion.p>
            <motion.div
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ duration: 0.5, delay: 0.2 }}
              className="flex flex-col sm:flex-row gap-4 justify-center"
            >
              <Link href="/register">
                <Button size="lg" className="w-full sm:w-auto text-lg px-8">
                  Start Chatting Free
                </Button>
              </Link>
              <Link href="/login">
                <Button
                  size="lg"
                  variant="outline"
                  className="w-full sm:w-auto text-lg px-8"
                >
                  Sign In
                </Button>
              </Link>
            </motion.div>
          </div>

          {/* Features */}
          <motion.div
            initial={{ opacity: 0, y: 40 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.6, delay: 0.4 }}
            className="mt-24 grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6"
          >
            <FeatureCard
              icon={<Zap className="w-6 h-6" />}
              title="Lightning Fast"
              description="Real-time messaging powered by WebSocket technology"
            />
            <FeatureCard
              icon={<Shield className="w-6 h-6" />}
              title="Secure"
              description="End-to-end encryption and JWT authentication"
            />
            <FeatureCard
              icon={<Users className="w-6 h-6" />}
              title="Groups"
              description="Create group conversations with unlimited members"
            />
            <FeatureCard
              icon={<MessageCircle className="w-6 h-6" />}
              title="Multi-device"
              description="Stay connected across all your devices"
            />
          </motion.div>
        </div>
      </div>

      {/* Footer */}
      <footer className="border-t border-border py-8">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex flex-col sm:flex-row items-center justify-between gap-4">
            <div className="flex items-center gap-2">
              <Image
                src="/gochat.svg"
                alt="GoChat"
                width={24}
                height={24}
                className="w-6 h-6"
              />
              <span className="text-sm text-muted-foreground">
                © 2024 GoChat. Built with Go & Next.js
              </span>
            </div>
            <div className="flex items-center gap-6">
              <a
                href="https://github.com/Beretta350/gochat"
                target="_blank"
                rel="noopener noreferrer"
                className="text-sm text-muted-foreground hover:text-foreground transition-colors"
              >
                GitHub
              </a>
            </div>
          </div>
        </div>
      </footer>
    </div>
  );
}

function FeatureCard({
  icon,
  title,
  description,
}: {
  icon: React.ReactNode;
  title: string;
  description: string;
}) {
  return (
    <div className="p-6 rounded-xl bg-card/50 border border-border hover:border-primary/50 transition-colors">
      <div className="w-12 h-12 rounded-lg bg-primary/10 flex items-center justify-center text-primary mb-4">
        {icon}
      </div>
      <h3 className="text-lg font-semibold mb-2">{title}</h3>
      <p className="text-sm text-muted-foreground">{description}</p>
    </div>
  );
}

