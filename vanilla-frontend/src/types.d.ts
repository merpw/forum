// POST request types

export type SignupForm = {
  username: string
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
  // rememberMe: boolean
}

export type CreatePostBody = {
  Title: string
  Content: string
  Description: string
  Categories: string[]
}

// GET response types

export type User = {
  id: number
  name: string
}

export type Post = {
  id: number
  title: string
  content: string
  description: string
  author: User
  date: string
  commentsCount: number
  likes: number
  dislikes: number
  categories: string
}

export type Comment = {
  id: number
  content: string
  author: User
  date: string
  likes: number
  dislikes: number
}

// WS Types
// Authored by Maksim

export type Message = {
  id: number
  chatId: number
  content: string
  // -1 for system messages
  authorId: number
  timestamp: string
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
