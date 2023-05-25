// Dummy in-memory database

export const chats = new Map<
  number,
  {
    // 0 = 1:1, 1 = group
    type: 0
  }
>()

export const messages = new Map<
  number,
  {
    chatId: number
    authorId: number
    content: string
    timestamp: string
  }
>()

// chatId -> userId[]
export const memberships = new Map<number, number[]>()

// chatId -> messageId[]
export const chatMessages = new Map<number, number[]>()

export const getLastMessageId = (chatId: number) => {
  const messages = chatMessages.get(chatId)
  return messages ? messages.slice(-1)[0] : null
}
