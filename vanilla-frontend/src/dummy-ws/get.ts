import { Chat, Message, WSGetResponse } from "./types"
import {
  chatMessages,
  chats,
  getLastMessageId,
  memberships,
  messages,
} from "./db"

// response data type for /chat/all
type ChatAllResponseData = number[]

// response data type for /chat/1/messages
type ChatMessagesResponseData = number[]

// response data type for /chat/1
type ChatResponseData = Chat | null

// response data type for /message/1
type MessageResponseData = Message | null

// response data type for /user/1/chat
type UserChatResponseData = number | null

const getHandler = ({ url, userId }: { url: string; userId: number }) => {
  // /chat/all, return array of ids of all users chats
  if (url.match(/\/chat\/all/)) {
    const userChats = Array.from(chats.entries())
      .filter(([chatId]) => memberships.get(chatId)?.includes(userId))
      .map(([chatId]) => ({
        id: chatId,
        lastMessageId: getLastMessageId(chatId) as number,
      }))
      .sort((a, b) => b.lastMessageId - a.lastMessageId)
      .map((chat) => chat.id)
    // Get chats where user is a member, sort by last message id

    const response: WSGetResponse<ChatAllResponseData> = {
      type: "get",
      item: {
        url,
        data: userChats,
      },
    }
    return response
  }

  // /chat/1/messages, return array of ids of messages
  if (url.match(/\/chat\/\d+\/messages/)) {
    const chatId = parseInt(url.split("/")[2])
    const data = chatMessages.get(chatId)

    if (!data) {
      throw new Error("Chat not found")
    }

    const response: WSGetResponse<ChatMessagesResponseData> = {
      type: "get",
      item: {
        url,
        data,
      },
    }
    return response
  }

  // /chat/1, return chat info
  if (url.match(/\/chat\/\d+/)) {
    const chatId = parseInt(url.split("/")[2])
    const chat = chats.get(chatId) || null

    if (!chat) {
      const response: WSGetResponse<ChatResponseData> = {
        type: "get",
        item: {
          url,
          data: null,
        },
      }
      return response
    }

    const data: ChatResponseData = {
      id: chatId, // TODO: maybe remove, can be inferred from url
      userId: memberships.get(chatId)?.find((id) => id !== userId) || userId,
      lastMessageId: getLastMessageId(chatId) as number,
    }

    const response: WSGetResponse<ChatResponseData> = {
      type: "get",
      item: {
        url,
        data,
      },
    }
    return response
  }

  // /message/1, return message info
  if (url.match(/\/message\/\d+/)) {
    const messageId = parseInt(url.split("/")[2])
    const message = messages.get(messageId) || null

    const data: MessageResponseData = message
      ? { ...message, id: messageId }
      : null

    const response: WSGetResponse<MessageResponseData> = {
      type: "get",
      item: {
        url,
        data,
      },
    }
    return response
  }

  // /user/1/chat, return chatId of chat between user and userId
  if (url.match(/\/user\/\d+\/chat/)) {
    const associatedId = parseInt(url.split("/")[2])
    const chatEntry = [...memberships.entries()].find(
      ([, ids]: [number, number[]]) => {
        return ids.includes(userId) && ids.includes(associatedId)
      }
    )

    const chatId = chatEntry ? chatEntry[0] : null

    const response: WSGetResponse<UserChatResponseData> = {
      type: "get",
      item: {
        url,
        data: chatId,
      },
    }
    return response
  }
  throw new Error("Invalid url")
}

export default getHandler
