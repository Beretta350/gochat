import type { Config } from "tailwindcss";

const config: Config = {
  darkMode: ["class"],
  content: [
    "./src/pages/**/*.{js,ts,jsx,tsx,mdx}",
    "./src/components/**/*.{js,ts,jsx,tsx,mdx}",
    "./src/app/**/*.{js,ts,jsx,tsx,mdx}",
  ],
  theme: {
    extend: {
      colors: {
        // GoChat Brand Colors
        background: {
          DEFAULT: "#07182a",
          secondary: "#0a2137",
          tertiary: "#0d2a45",
        },
        foreground: {
          DEFAULT: "#ffffff",
          muted: "#b8fce5",
        },
        primary: {
          DEFAULT: "#11e3de",
          hover: "#21ffe0",
          foreground: "#07182a",
        },
        secondary: {
          DEFAULT: "#00b4c8",
          hover: "#11e3de",
          foreground: "#ffffff",
        },
        accent: {
          DEFAULT: "#21ffe0",
          hover: "#b8fce5",
          foreground: "#07182a",
        },
        muted: {
          DEFAULT: "#0d2a45",
          foreground: "#b8fce5",
        },
        card: {
          DEFAULT: "#0a2137",
          foreground: "#ffffff",
        },
        border: "#1a3a5c",
        input: "#1a3a5c",
        ring: "#11e3de",
        destructive: {
          DEFAULT: "#ef4444",
          foreground: "#ffffff",
        },
        success: {
          DEFAULT: "#21ffe0",
          foreground: "#07182a",
        },
      },
      borderRadius: {
        lg: "0.75rem",
        md: "0.5rem",
        sm: "0.25rem",
      },
      fontFamily: {
        sans: ["var(--font-geist-sans)", "system-ui", "sans-serif"],
        mono: ["var(--font-geist-mono)", "monospace"],
      },
      keyframes: {
        "fade-in": {
          "0%": { opacity: "0", transform: "translateY(10px)" },
          "100%": { opacity: "1", transform: "translateY(0)" },
        },
        "fade-out": {
          "0%": { opacity: "1", transform: "translateY(0)" },
          "100%": { opacity: "0", transform: "translateY(10px)" },
        },
        "slide-in-right": {
          "0%": { transform: "translateX(100%)" },
          "100%": { transform: "translateX(0)" },
        },
        "slide-in-left": {
          "0%": { transform: "translateX(-100%)" },
          "100%": { transform: "translateX(0)" },
        },
        pulse: {
          "0%, 100%": { opacity: "1" },
          "50%": { opacity: "0.5" },
        },
        bounce: {
          "0%, 100%": { transform: "translateY(0)" },
          "50%": { transform: "translateY(-5px)" },
        },
        "typing-dot": {
          "0%, 60%, 100%": { transform: "translateY(0)" },
          "30%": { transform: "translateY(-5px)" },
        },
      },
      animation: {
        "fade-in": "fade-in 0.3s ease-out",
        "fade-out": "fade-out 0.3s ease-out",
        "slide-in-right": "slide-in-right 0.3s ease-out",
        "slide-in-left": "slide-in-left 0.3s ease-out",
        pulse: "pulse 2s cubic-bezier(0.4, 0, 0.6, 1) infinite",
        bounce: "bounce 1s infinite",
        "typing-dot": "typing-dot 1.4s infinite",
      },
    },
  },
  plugins: [require("tailwindcss-animate")],
};

export default config;

