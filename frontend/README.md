# 🎨 GoChat Frontend

Modern real-time chat interface built with Next.js, React, and TypeScript.

## 🚀 Tech Stack

| Technology | Purpose |
|------------|---------|
| **Next.js 14** | React framework with App Router |
| **TypeScript** | Type safety |
| **Tailwind CSS** | Styling |
| **Redux Toolkit** | State management |
| **RTK Query** | Data fetching & caching |
| **React Hook Form** | Form handling |
| **Zod** | Schema validation |
| **Framer Motion** | Animations |
| **Radix UI** | Accessible components |

## 📁 Project Structure

```
frontend/
├── src/
│   ├── app/                    # Next.js App Router
│   │   ├── (auth)/             # Auth pages (login, register)
│   │   ├── (chat)/             # Chat pages
│   │   ├── layout.tsx          # Root layout
│   │   ├── page.tsx            # Landing page
│   │   └── globals.css         # Global styles
│   ├── components/
│   │   ├── ui/                 # Base UI components
│   │   ├── chat/               # Chat components
│   │   └── providers/          # Context providers
│   ├── hooks/                  # Custom hooks
│   │   ├── useAuth.ts          # Authentication hook
│   │   └── useWebSocket.ts     # WebSocket connection hook
│   ├── store/                  # Redux store
│   │   ├── api/                # RTK Query APIs
│   │   ├── slices/             # Redux slices
│   │   └── store.ts            # Store configuration
│   ├── lib/                    # Utilities
│   └── types/                  # TypeScript types
├── public/                     # Static assets
├── tailwind.config.ts          # Tailwind configuration
└── next.config.ts              # Next.js configuration
```

## 🛠️ Development

### Prerequisites

- Node.js 20+
- npm or yarn

### Setup

```bash
# Install dependencies
npm install

# Start development server
npm run dev

# Build for production
npm run build

# Start production server
npm start
```

### Environment Variables

Create a `.env.local` file:

```env
NEXT_PUBLIC_API_URL=http://localhost:8080
NEXT_PUBLIC_WS_URL=ws://localhost:8080
```

## 🎨 Design System

### Colors (from brand)

| Color | Hex | Usage |
|-------|-----|-------|
| Background | `#07182a` | Main background |
| Primary | `#11e3de` | Buttons, accents |
| Accent | `#21ffe0` | Highlights |
| Secondary | `#00b4c8` | Secondary elements |
| Muted | `#b8fce5` | Muted text |

### Components

Built on top of [shadcn/ui](https://ui.shadcn.com/) patterns:

- `Button` - Various button styles
- `Input` - Form inputs
- `Avatar` - User avatars
- `Dialog` - Modal dialogs
- `ScrollArea` - Scrollable containers
- `Toast` - Notifications

## 📱 Features

- [x] Responsive design (mobile-first)
- [x] Real-time messaging via WebSocket
- [x] JWT authentication
- [x] Conversation list with search
- [x] Direct & group conversations
- [x] Message history
- [x] Smooth animations
- [x] Dark theme
- [ ] Typing indicators
- [ ] Read receipts
- [ ] File uploads
- [ ] Push notifications

