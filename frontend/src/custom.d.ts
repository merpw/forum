export type Post = {
  id: number
  title: string
  content: string
  likes: number
  dislikes: number
  date: string
  author: User
  user_reaction: number | undefined
  comments: Comment[]
}

export type User = {
  name: string
  id: number
}

export type Comment = {
  author: User
  text: string
  date: string
  likes: number
  user_reaction: number | undefined
}
