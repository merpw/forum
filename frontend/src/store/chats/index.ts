import { createSlice, PayloadAction } from "@reduxjs/toolkit"

import { Chat, Message, WSGetResponse, WSPostResponse } from "@/ws"

// TODO: split into separate files

// Map-like object without ES6 Map usage (as it is non-serializable)
type ObjectMap<K extends number | string | symbol, T> = Partial<Record<K, T>>

const initialState: {
  activeChatId: number | null
  unreadMessagesChatIds: number[]
  chatIds: number[] | undefined
  chats: ObjectMap<number, Chat | null>
  chatMessages: ObjectMap<number, number[]>
  messages: ObjectMap<number, Message | null>
  userChats: ObjectMap<number, number | null>
  usersOnline: number[] | undefined
} = {
  activeChatId: null,
  unreadMessagesChatIds: [],
  chatIds: undefined,
  chats: {},
  chatMessages: {},
  messages: {},
  userChats: {},
  usersOnline: undefined,
}

const chatSlice = createSlice({
  name: "chat",
  initialState,
  reducers: {
    setActiveChatId(state, action: PayloadAction<number | null>) {
      state.activeChatId = action.payload

      state.unreadMessagesChatIds = state.unreadMessagesChatIds.filter(
        (id) => id !== action.payload
      )
    },

    reset: () => initialState,

    handleChatAll: {
      reducer: (state, action: PayloadAction<number[]>) => {
        state.chatIds = action.payload
      },
      prepare: (response: WSGetResponse<number[]>) => ({ payload: response.item.data }),
    },

    handleChat: {
      reducer: (state, action: PayloadAction<{ chatId: number; data: Chat | null }>) => {
        state.chats[action.payload.chatId] = action.payload.data
        if (action.payload.data) {
          state.userChats[action.payload.data.companionId] = action.payload.chatId
        }
      },
      prepare: (response: WSGetResponse<Chat | null>) => {
        const chatId = +response.item.url.split("/")[2]
        return { payload: { chatId, data: response.item.data } }
      },
    },

    handleChatMessages: {
      reducer: (state, action: PayloadAction<{ chatId: number; messages: number[] }>) => {
        state.chatMessages[action.payload.chatId] = action.payload.messages
      },
      prepare: (response: WSGetResponse<number[]>) => {
        const chatId = +response.item.url.split("/")[2]
        return { payload: { chatId, messages: response.item.data } }
      },
    },

    // TODO: maybe rename to handleChatMessageAdd
    handleChatMessage: {
      reducer: (state, action: PayloadAction<{ chatId: number; messageId: number }>) => {
        state.chatMessages[action.payload.chatId]?.push(action.payload.messageId)
        if (state.activeChatId !== action.payload.chatId) {
          state.unreadMessagesChatIds.push(action.payload.chatId)
        }
        const changedChat = state.chats[action.payload.chatId]
        if (!changedChat) return

        changedChat.lastMessageId = action.payload.messageId

        state.chatIds?.sort((a, b) => {
          const chatA = state.chats[a]
          const chatB = state.chats[b]
          if (!chatA) return 1
          if (!chatB) return -1
          return chatB.lastMessageId - chatA.lastMessageId
        })
      },
      prepare: (response: WSPostResponse<number>) => {
        const chatId = +response.item.url.split("/")[2]
        return { payload: { chatId, messageId: response.item.data } }
      },
    },

    handleMessage: {
      reducer: (state, action: PayloadAction<{ messageId: number; data: Message | null }>) => {
        state.messages[action.payload.messageId] = action.payload.data
      },
      prepare(response: WSGetResponse<Message>) {
        const messageId = +response.item.url.split("/")[2]
        return { payload: { messageId, data: response.item.data } }
      },
    },
    handleUserChat: {
      reducer: (state, action: PayloadAction<{ userId: number; chatId: number | null }>) => {
        state.userChats[action.payload.userId] = action.payload.chatId
        if (action.payload.chatId !== null && !state.chatIds?.includes(action.payload.chatId)) {
          state.chatIds?.unshift(action.payload.chatId)
        }
      },
      prepare: (response: WSGetResponse<number | null>) => {
        const userId = +response.item.url.split("/")[2]
        return { payload: { userId, chatId: response.item.data } }
      },
    },

    handleChatCreate: {
      reducer: (state, action: PayloadAction<{ chatId: number; userId: number }>) => {
        state.chatIds?.unshift(action.payload.chatId)
        state.userChats[action.payload.userId] = action.payload.chatId
      },
      prepare: (response: WSPostResponse<{ chatId: number; userId: number }>) => {
        return { payload: response.item.data }
      },
    },

    handleUsersOnline: {
      reducer: (state, action: PayloadAction<number[]>) => {
        state.usersOnline = action.payload
      },
      prepare: (response: WSGetResponse<number[]>) => {
        return { payload: response.item.data }
      },
    },
  },
})

export const chatActions = chatSlice.actions

export const chatHandlers = [
  {
    regex: /^\/chat\/all$/,
    handler: chatActions.handleChatAll,
  },
  {
    regex: /^\/chat\/\d+\/messages$/,
    handler: chatActions.handleChatMessages,
  },
  {
    regex: /^\/chat\/\d+$/,
    handler: chatActions.handleChat,
  },
  {
    regex: /^\/message\/\d+$/,
    handler: chatActions.handleMessage,
  },
  {
    regex: /^\/chat\/\d+\/message$/,
    handler: chatActions.handleChatMessage,
  },
  {
    regex: /^\/users\/\d+\/chat$/,
    handler: chatActions.handleUserChat,
  },
  {
    regex: /^\/chat\/create$/,
    handler: chatActions.handleChatCreate,
  },
  {
    regex: /^\/users\/online$/,
    handler: chatActions.handleUsersOnline,
  },
]

export default chatSlice.reducer
