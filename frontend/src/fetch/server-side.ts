import { Post, User } from "../custom"
import { posts, users } from "./dummy"

// TODO: add data fetching
export const getPostsLocal = (): Promise<Post[]> => Promise.resolve(posts)
// fetch(`http://${process.env.FORUM_BACKEND_LOCALHOST}/api/artists/`)
//   .then((res) => res.json())
//   .then((data) => data as Post[])

export const getPostLocal = (id: number): Promise<Post | undefined> =>
  Promise.resolve(posts.find((post) => post.id == id))

export const getUsersLocal = (): Promise<User[]> => Promise.resolve(users)
export const getUserLocal = (id: number): Promise<User | undefined> =>
  Promise.resolve(users.find((user) => user.id == id))
