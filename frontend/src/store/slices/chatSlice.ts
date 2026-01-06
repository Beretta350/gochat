import { createSlice, PayloadAction } from "@reduxjs/toolkit";
import type { Conversation, Message } from "@/types";

interface ChatState {
  activeConversationId: string | null;
  conversations: Conversation[];
  messages: Record<string, Message[]>;
  typingUsers: Record<string, string[]>; // conversationId -> userIds
  isConnected: boolean;
  unreadCounts: Record<string, number>;
  onlineUsers: string[]; // List of online user IDs
}

const initialState: ChatState = {
  activeConversationId: null,
  conversations: [],
  messages: {},
  typingUsers: {},
  isConnected: false,
  unreadCounts: {},
  onlineUsers: [],
};

const chatSlice = createSlice({
  name: "chat",
  initialState,
  reducers: {
    setActiveConversation: (state, action: PayloadAction<string | null>) => {
      state.activeConversationId = action.payload;
      if (action.payload) {
        state.unreadCounts[action.payload] = 0;
      }
    },
    setConversations: (state, action: PayloadAction<Conversation[]>) => {
      state.conversations = action.payload;
    },
    addConversation: (state, action: PayloadAction<Conversation>) => {
      const exists = state.conversations.find(
        (c) => c.id === action.payload.id
      );
      if (!exists) {
        state.conversations.unshift(action.payload);
      }
    },
    updateConversation: (state, action: PayloadAction<Partial<Conversation> & { id: string }>) => {
      const index = state.conversations.findIndex(
        (c) => c.id === action.payload.id
      );
      if (index !== -1) {
        state.conversations[index] = {
          ...state.conversations[index],
          ...action.payload,
        };
      }
    },
    setMessages: (
      state,
      action: PayloadAction<{ conversationId: string; messages: Message[] }>
    ) => {
      state.messages[action.payload.conversationId] = action.payload.messages;
    },
    addMessage: (state, action: PayloadAction<Message>) => {
      const { conversation_id } = action.payload;
      if (!state.messages[conversation_id]) {
        state.messages[conversation_id] = [];
      }
      
      // Check if message already exists
      const exists = state.messages[conversation_id].find(
        (m) => m.id === action.payload.id
      );
      if (!exists) {
        state.messages[conversation_id].push(action.payload);
      }

      // Update conversation's last message
      const convIndex = state.conversations.findIndex(
        (c) => c.id === conversation_id
      );
      if (convIndex !== -1) {
        state.conversations[convIndex].last_message = action.payload;
        // Move conversation to top
        const [conv] = state.conversations.splice(convIndex, 1);
        state.conversations.unshift(conv);
      }

      // Increment unread count if not active conversation
      if (conversation_id !== state.activeConversationId) {
        state.unreadCounts[conversation_id] =
          (state.unreadCounts[conversation_id] || 0) + 1;
      }
    },
    prependMessages: (
      state,
      action: PayloadAction<{ conversationId: string; messages: Message[] }>
    ) => {
      const { conversationId, messages } = action.payload;
      if (!state.messages[conversationId]) {
        state.messages[conversationId] = [];
      }
      state.messages[conversationId] = [
        ...messages,
        ...state.messages[conversationId],
      ];
    },
    setTypingUser: (
      state,
      action: PayloadAction<{
        conversationId: string;
        userId: string;
        isTyping: boolean;
      }>
    ) => {
      const { conversationId, userId, isTyping } = action.payload;
      if (!state.typingUsers[conversationId]) {
        state.typingUsers[conversationId] = [];
      }
      if (isTyping) {
        if (!state.typingUsers[conversationId].includes(userId)) {
          state.typingUsers[conversationId].push(userId);
        }
      } else {
        state.typingUsers[conversationId] = state.typingUsers[
          conversationId
        ].filter((id) => id !== userId);
      }
    },
    setConnected: (state, action: PayloadAction<boolean>) => {
      state.isConnected = action.payload;
      // Clear online users when disconnected
      if (!action.payload) {
        state.onlineUsers = [];
      }
    },
    clearUnread: (state, action: PayloadAction<string>) => {
      state.unreadCounts[action.payload] = 0;
    },
    // Set initial list of online users (received on WebSocket connect)
    setOnlineUsers: (state, action: PayloadAction<string[]>) => {
      state.onlineUsers = action.payload;
    },
    // Update single user's online status
    setUserOnlineStatus: (
      state,
      action: PayloadAction<{ userId: string; isOnline: boolean }>
    ) => {
      const { userId, isOnline } = action.payload;
      if (isOnline) {
        if (!state.onlineUsers.includes(userId)) {
          state.onlineUsers.push(userId);
        }
      } else {
        state.onlineUsers = state.onlineUsers.filter((id) => id !== userId);
      }
    },
  },
});

export const {
  setActiveConversation,
  setConversations,
  addConversation,
  updateConversation,
  setMessages,
  addMessage,
  prependMessages,
  setTypingUser,
  setConnected,
  clearUnread,
  setOnlineUsers,
  setUserOnlineStatus,
} = chatSlice.actions;
export default chatSlice.reducer;

