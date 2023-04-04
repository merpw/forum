export type Post = {
  id: number
  title: string
  content: string
  likes_count: number
  dislikes_count: number
  comments_count: number
  date: string
  author: User
  comments: Comment[]
  categories: string
}

export type User = {
  name: string
  id: number
  email: string
  first_name: string
  last_name: string
  age: string
  gender: string
}

export type Comment = {
  id: number
  author: User
  content: string
  likes_count: number
  dislikes_count: number
  date: string
}
