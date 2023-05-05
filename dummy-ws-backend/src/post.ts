import { WSPostResponse } from "./types"
import { chatMessages, chats, memberships, messages } from "./db"
import getHandler from "./get"
import { connections } from "./server"

// request data type for /chat/create
type CreateChatRequestData = {
  userId: number
}

// request data type for /chat/1/message
type SendMessageRequestData = {
  content: string
}

// response data type for /chat/1/message
type SendMessageResponseData = {
  messageId: number
}

const postHandler = ({
  url,
  data,
  userId,
}: {
  url: string
  data: never
  userId: number
}) => {
  if (url.match(/\/chat\/create/)) {
    const { userId: chatUserId } = data as CreateChatRequestData

    const chatId = chats.size

    chats.set(chatId, {
      type: 0,
    })

    memberships.set(chatId, [userId, chatUserId])

    const firstMessageId = messages.size

    messages.set(firstMessageId, {
      chatId,
      authorId: -1,
      content: "Chat created",
      timestamp: new Date().toISOString(),
    })

    chatMessages.set(chatId, [firstMessageId])

    const chatUserConn = [...connections.entries()].find(
      ([, { userId: connUserId }]) => connUserId === chatUserId
    )?.[0]

    chatUserConn?.send(
      JSON.stringify(
        getHandler({
          url: `/users/${userId}/chat`,
          userId: chatUserId,
        })
      )
    )

    return getHandler({
      url: `/users/${chatUserId}/chat`,
      userId,
    })
  }

  if (url.match(/\/chat\/\d+\/message/)) {
    const chatId = parseInt(url.split("/")[2])

    const chatHistory = chatMessages.get(chatId)

    if (!chatHistory) {
      throw new Error("Chat not found")
    }

    const { content } = data as SendMessageRequestData
    const messageId = messages.size

    messages.set(messageId, {
      chatId,
      authorId: userId,
      content,
      timestamp: new Date().toISOString(),
    })

    chatHistory.push(messageId)

    const response: WSPostResponse<SendMessageResponseData> = {
      type: "post",
      item: {
        url,
        data: {
          messageId,
        },
      },
    }

    const chatMembers = memberships.get(chatId)

    ;[...connections.entries()].forEach(
      ([connection, { userId: connUserId }]) => {
        userId !== connUserId &&
          chatMembers?.includes(connUserId) &&
          connection.send(JSON.stringify(response))
      }
    )

    return response
  }

  throw new Error("Invalid url")
}

export default postHandler
