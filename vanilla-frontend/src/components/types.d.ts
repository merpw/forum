declare module "LoginModule" {
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
}

declare module "POST" {
  export type CreatePostBody = {
    Title: string
    Content: string
    Description: string
    Categories: string[]
  }
}

// Change name to GET Module
declare module "ContentModule" {
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
}
