export { store } from "./store";
export type { RootState, AppDispatch } from "./store";
export { useAppDispatch, useAppSelector } from "./hooks";

// Auth
export {
  authApi,
  useLoginMutation,
  useRegisterMutation,
  useLogoutMutation,
  useGetMeQuery,
  useLazyGetMeQuery,
} from "./api/authApi";
export {
  setCredentials,
  setUser,
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
