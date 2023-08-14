export type Message = {
  id: number
  chatId: number
  content: string
  // -1 for system messages
  authorId: number
  timestamp: string
}

export type Chat = {
  id: number
  lastMessageId: number
} & (
  | {
      userId: number
    }
  | {
      groupId: number
    }
)

export type TypingData = {
  userId: number
  isTyping: boolean
}

export type WSBase = {
  type: string
  item: object
}

export type WSHandshake = WSBase & {
  type: "handshake"
  item: {
    token: string
  }
}

// Get requests and responses

export type WSGetRequest = WSBase & {
  type: "get"
  item: {
    url: string
  }
}

export type WSGetResponse<T> = WSBase & {
  type: "get"
  item: {
    url: string
    data: T
  }
}

//

// Post requests and responses

export type WSPostRequest<T> = WSBase & {
  type: "post"
  item: {
    url: string
    data: T
  }
}

export type WSPostResponse<T> = WSBase & {
  type: "post"
  item: {
    url: string
    data: T
  }
}

export type WSErrorResponse = WSBase & {
  type: "error"
  item: {
    message: string
  }
}

export type WebSocketRequest<T> = WSHandshake | WSGetRequest | WSPostRequest<T>

export type WebSocketResponse<T> =
  | WSHandshake
  | WSGetResponse<T>
  | WSPostResponse<T>
  | WSErrorResponse
