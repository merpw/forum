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
  group_id?: number
  privacy: Privacy
  selectedUsers?: number[]
}

export type Privacy = "Public" | "Private" | "SuperPrivate"

export type User = {
  username: string
  id: number
  email?: string
  first_name?: string
  last_name?: string
  dob?: string
  gender?: string
  avatar?: string
  bio?: string
  privacy?: boolean
  follow_status?: FollowStatus

  followers_count?: number
  following_count?: number
}

/** 0 - not following, 1 - following, 2 - requested */
export type FollowStatus = 0 | 1 | 2

export type Comment = {
  id: number
  author: User
  content: string
  likes_count: number
  dislikes_count: number
  date: string
}
