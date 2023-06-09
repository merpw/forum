export type Post = {
  id: number
  title: string
  content: string
  description: string
  likes_count: number
  dislikes_count: number
  comments_count: number
  date: string
  author: User
  comments: Comment[]
  categories: string
}

export type User = {
  username: string
  id: number
  email: string
  first_name?: string
  last_name?: string
  dob?: string
  gender?: string
}

export type Comment = {
  id: number
  author: User
  content: string
  likes_count: number
  dislikes_count: number
  date: string
}
