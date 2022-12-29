import { Post } from "../../custom"
import axios from "axios"

export const getPostsLocal = (): Promise<Post[]> =>
  axios(`http://${process.env.FORUM_BACKEND_LOCALHOST}/api/posts`).then((res) => res.data)

export const getPostLocal = (id: number) =>
  axios<Post | undefined>(`http://${process.env.FORUM_BACKEND_LOCALHOST}/api/posts/${id}`)
    .then((res) => res.data)
    .catch(() => undefined)

export const getUserPostsLocal = (user_id: number) =>
  axios<Post[]>(`http://${process.env.FORUM_BACKEND_LOCALHOST}/api/user/${user_id}/posts`)
    .then((res) => {
      return { posts: res.data }
    })
    .catch(() => {
      return { posts: [] as Post[] }
    })
