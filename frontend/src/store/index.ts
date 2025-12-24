export { store } from "./store";
export type { RootState, AppDispatch } from "./store";
export { useAppDispatch, useAppSelector } from "./hooks";

// Auth
export { authApi, useLoginMutation, useRegisterMutation, useGetMeQuery } from "./api/authApi";
export {
  setCredentials,
  setUser,
  setTokens,
  logout,
  setLoading,
} from "./slices/authSlice";

// Conversations
export {
  conversationsApi,
  useGetConversationsQuery,
  useGetConversationQuery,
  useCreateConversationMutation,
  useGetMessagesQuery,
  useLazyGetMessagesQuery,
} from "./api/conversationsApi";
export {
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
} from "./slices/chatSlice";

