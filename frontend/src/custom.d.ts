export type Post = {
  id: number
  title: string
  content: string
  likes: number
  dislikes: number
  date: string
  author: User
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
}
