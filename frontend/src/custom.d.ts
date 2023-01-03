export type Post = {
  id: number
  title: string
  content: string
  likes_count: number
  dislikes_count: number | undefined // defined only for post's author
  comments_count: number
  date: string
  author: User
  comments: Comment[]
  category: string
}

export type User = {
  name: string
  id: number
  email: string
}

export type Comment = {
  id : number
  author: User
  content: string
  likes_count: number
  dislikes_count: number | undefined // defined only for comment's author
  date: string
}
