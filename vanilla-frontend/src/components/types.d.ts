export type SignupForm = {
  name: string
  email: string
  password: string
  first_name: string
  last_name: string
  dob: string
  gender: string
}

export type LoginForm = {
  login: string
  password: string
  rememberMe: boolean
}

export type CreatePostBody = {
  Title: string
  Content: string
  Description: string
  Categories: string[]
}

export type SafePost = {
  HTML: string
  Content: string
}

// Change name to GET Module
export type User = {
  Id: number
  Name: string
}

export type Post = {
  Id: number
  Title: string
  Content: string
  Description: string
  Author: User
  Date: string
  CommentsCount: number
  Likes: number
  Dislikes: number
  Categories: string
}

export type Comment = {
  Id: number
  Content: string
  Author: User
  Date: string
  Likes: number
  Dislikes: number
}

export type Reaction = {
  Reaction: number
  Likes: number
  Dislikes: number
}

export type ActiveUser = {
  Name: string
  ID: number
  Online: boolean
  UnreadMSG: boolean
}

export type InactiveUser = {
  Name: string
  ID: number
  Online: boolean
}

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
  userId: number
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
