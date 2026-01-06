// User types
export interface User {
  id: string;
  email: string;
  username: string;
  is_active: boolean;
  created_at: string;
}

export interface AuthTokens {
  access_token: string;
  refresh_token: string;
  expires_in: number;
}

export interface AuthResponse {
  user: User;
  tokens: AuthTokens;
}

export interface LoginRequest {
  email: string;
  password: string;
}

export interface RegisterRequest {
  email: string;
  username: string;
  password: string;
}

// Conversation types
export type ConversationType = "direct" | "group";

export interface Participant {
  id: string;
  email: string;
  username: string;
  is_active: boolean;
  role?: "admin" | "member";
  joined_at: string;
  created_at: string;
}

export interface Message {
  id: string;
  conversation_id: string;
  sender_id: string;
  sender_username?: string;
  content: string;
  type: "text" | "image" | "file" | "audio";
  sent_at: string | number;
}

export interface Conversation {
  id: string;
  type: ConversationType;
  name?: string;
  participants: Participant[];
  last_message?: Message;
  created_at: string;
  updated_at: string;
}

export interface CreateConversationRequest {
  participant_id?: string;
  participant_email?: string;
  participant_ids?: string[];
  participant_emails?: string[];
  name?: string;
}

export interface MessagesPage {
  messages: Message[];
  has_more: boolean;
  next_cursor?: string;
}

// WebSocket types
export interface WebSocketMessage {
  id: string;
  conversation_id: string;
  sender_id: string;
  sender_username?: string;
  content: string;
  type: string;
  sent_at: number;
}

export interface WebSocketError {
  error: boolean;
  message: string;
}

// Presence types
export interface PresenceEvent {
  type: "presence";
  user_id: string;
  username?: string;
  status: "online" | "offline";
}

export interface PresenceListEvent {
  type: "presence_list";
  online_users: string[];
}

export interface SendMessageRequest {
  conversation_id: string;
  content: string;
  type?: string;
}

