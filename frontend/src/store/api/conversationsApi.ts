import { createApi } from "@reduxjs/toolkit/query/react";
import { baseQuery } from "./baseApi";
import type {
  Conversation,
  CreateConversationRequest,
  MessagesPage,
  Message,
  Participant,
} from "@/types";

// API response types (different from our internal types)
interface ConversationResponse {
  conversation: {
    id: string;
    type: "direct" | "group";
    name?: string;
    created_at: string;
    updated_at: string;
  };
  participants: Participant[];
  last_message?: Message;
}

interface GetConversationsResponse {
  conversations: ConversationResponse[];
  count: number;
}

export const conversationsApi = createApi({
  reducerPath: "conversationsApi",
  baseQuery,
  tagTypes: ["Conversations", "Messages"],
  endpoints: (builder) => ({
    getConversations: builder.query<Conversation[], void>({
      query: () => "/conversations",
      transformResponse: (response: GetConversationsResponse): Conversation[] => {
        if (!response?.conversations) return [];
        return response.conversations.map((item) => ({
          id: item.conversation.id,
          type: item.conversation.type,
          name: item.conversation.name,
          participants: item.participants,
          last_message: item.last_message,
          created_at: item.conversation.created_at,
          updated_at: item.conversation.updated_at,
        }));
      },
      providesTags: ["Conversations"],
    }),
    getConversation: builder.query<Conversation, string>({
      query: (id) => `/conversations/${id}`,
      transformResponse: (response: ConversationResponse): Conversation => ({
        id: response.conversation.id,
        type: response.conversation.type,
        name: response.conversation.name,
        participants: response.participants,
        created_at: response.conversation.created_at,
        updated_at: response.conversation.updated_at,
      }),
      providesTags: (_result, _error, id) => [{ type: "Conversations", id }],
    }),
    createConversation: builder.mutation<Conversation, CreateConversationRequest>({
      query: (data) => ({
        url: "/conversations",
        method: "POST",
        body: data,
      }),
      transformResponse: (response: ConversationResponse): Conversation => ({
        id: response.conversation.id,
        type: response.conversation.type,
        name: response.conversation.name,
        participants: response.participants,
        created_at: response.conversation.created_at,
        updated_at: response.conversation.updated_at,
      }),
      invalidatesTags: ["Conversations"],
    }),
    getMessages: builder.query<
      MessagesPage,
      { conversationId: string; cursor?: string; limit?: number }
    >({
      query: ({ conversationId, cursor, limit = 50 }) => {
        const params = new URLSearchParams();
        params.set("limit", limit.toString());
        if (cursor) params.set("cursor", cursor);
        return `/conversations/${conversationId}/messages?${params.toString()}`;
      },
      providesTags: (_result, _error, { conversationId }) => [
        { type: "Messages", id: conversationId },
      ],
    }),
  }),
});

export const {
  useGetConversationsQuery,
  useGetConversationQuery,
  useCreateConversationMutation,
  useGetMessagesQuery,
  useLazyGetMessagesQuery,
} = conversationsApi;
